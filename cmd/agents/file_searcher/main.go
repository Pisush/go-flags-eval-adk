package main

import (
	"bufio"
	"flag"
	"fmt"
	"log"
	"os"
	"path/filepath"
	"runtime"
	"runtime/debug"
	"strings"
	"sync"
	"time"
)

var (
	pattern = flag.String("pattern", "func", "Pattern to search for")
	dir     = flag.String("dir", "./testdata", "Directory to search in")
	workers = flag.Int("workers", 0, "Number of worker goroutines (0 = GOMAXPROCS)")
)

type Match struct {
	File    string
	Line    int
	Content string
}

func main() {
	flag.Parse()

	start := time.Now()

	if *workers == 0 {
		*workers = runtime.GOMAXPROCS(-1)
	}

	// Report configuration
	fmt.Printf("File Searcher Agent\n")
	fmt.Printf("===================\n")
	fmt.Printf("GOMAXPROCS: %d\n", runtime.GOMAXPROCS(-1))
	fmt.Printf("GOGC: %d\n", debug.SetGCPercent(-1))
	debug.SetGCPercent(debug.SetGCPercent(-1)) // Restore
	fmt.Printf("Pattern: %s\n", *pattern)
	fmt.Printf("Directory: %s\n", *dir)
	fmt.Printf("Workers: %d\n", *workers)
	fmt.Printf("\n")

	// Find all Go files
	var files []string
	err := filepath.Walk(*dir, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}
		if !info.IsDir() && strings.HasSuffix(path, ".go") {
			files = append(files, path)
		}
		return nil
	})

	if err != nil {
		log.Fatalf("Failed to walk directory: %v", err)
	}

	fmt.Printf("Found %d Go files to search\n\n", len(files))

	// Search files concurrently
	fileChan := make(chan string, len(files))
	matchChan := make(chan Match, 100)

	var wg sync.WaitGroup

	// Start workers
	for i := 0; i < *workers; i++ {
		wg.Add(1)
		go worker(&wg, fileChan, matchChan, *pattern)
	}

	// Send files to workers
	go func() {
		for _, file := range files {
			fileChan <- file
		}
		close(fileChan)
	}()

	// Collect matches
	go func() {
		wg.Wait()
		close(matchChan)
	}()

	matches := []Match{}
	for match := range matchChan {
		matches = append(matches, match)
	}

	elapsed := time.Since(start)

	// Print results
	var ms runtime.MemStats
	runtime.ReadMemStats(&ms)

	fmt.Printf("\nResults:\n")
	fmt.Printf("========\n")
	fmt.Printf("Files searched: %d\n", len(files))
	fmt.Printf("Matches found: %d\n", len(matches))
	fmt.Printf("Duration: %v\n", elapsed)
	fmt.Printf("Memory allocated: %.2f MB\n", float64(ms.TotalAlloc)/(1024*1024))
	fmt.Printf("GC runs: %d\n", ms.NumGC)

	// Print first 10 matches
	if len(matches) > 0 {
		fmt.Printf("\nFirst 10 matches:\n")
		for i, match := range matches {
			if i >= 10 {
				break
			}
			fmt.Printf("%s:%d: %s\n", match.File, match.Line, strings.TrimSpace(match.Content))
		}
	}
}

func worker(wg *sync.WaitGroup, files <-chan string, matches chan<- Match, pattern string) {
	defer wg.Done()

	for file := range files {
		searchFile(file, pattern, matches)
	}
}

func searchFile(filename, pattern string, matches chan<- Match) {
	f, err := os.Open(filename)
	if err != nil {
		return
	}
	defer f.Close()

	scanner := bufio.NewScanner(f)
	lineNum := 0

	for scanner.Scan() {
		lineNum++
		line := scanner.Text()
		if strings.Contains(line, pattern) {
			matches <- Match{
				File:    filename,
				Line:    lineNum,
				Content: line,
			}
		}
	}
}
