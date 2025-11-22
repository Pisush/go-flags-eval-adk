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
	target    = flag.String("target", "./testdata", "Target directory for refactoring")
	operation = flag.String("operation", "rename", "Operation: rename, add-comments, format")
	oldName   = flag.String("old", "oldVar", "Old variable name (for rename)")
	newName   = flag.String("new", "newVar", "New variable name (for rename)")
)

func main() {
	flag.Parse()

	start := time.Now()

	// Report configuration
	fmt.Printf("Refactor Agent\n")
	fmt.Printf("==============\n")
	fmt.Printf("GOMAXPROCS: %d\n", runtime.GOMAXPROCS(-1))
	fmt.Printf("GOGC: %d\n", debug.SetGCPercent(-1))
	debug.SetGCPercent(debug.SetGCPercent(-1)) // Restore
	fmt.Printf("Target: %s\n", *target)
	fmt.Printf("Operation: %s\n", *operation)
	if *operation == "rename" {
		fmt.Printf("Rename: %s -> %s\n", *oldName, *newName)
	}
	fmt.Printf("\n")

	// Find all Go files
	var files []string
	err := filepath.Walk(*target, func(path string, info os.FileInfo, err error) error {
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

	fmt.Printf("Found %d Go files to refactor\n\n", len(files))

	// Refactor files concurrently
	var wg sync.WaitGroup
	results := make(chan refactorResult, len(files))
	sem := make(chan struct{}, runtime.GOMAXPROCS(-1))

	for _, file := range files {
		wg.Add(1)
		go func(filename string) {
			defer wg.Done()
			sem <- struct{}{}
			defer func() { <-sem }()

			result := refactorFile(filename, *operation, *oldName, *newName)
			results <- result
		}(file)
	}

	go func() {
		wg.Wait()
		close(results)
	}()

	// Collect results
	totalChanges := 0
	filesModified := 0

	for result := range results {
		if result.err != nil {
			log.Printf("Error processing %s: %v", result.filename, result.err)
		} else if result.changes > 0 {
			filesModified++
			totalChanges += result.changes
		}
	}

	elapsed := time.Since(start)

	// Print statistics
	var ms runtime.MemStats
	runtime.ReadMemStats(&ms)

	fmt.Printf("\nResults:\n")
	fmt.Printf("========\n")
	fmt.Printf("Files processed: %d\n", len(files))
	fmt.Printf("Files modified: %d\n", filesModified)
	fmt.Printf("Total changes: %d\n", totalChanges)
	fmt.Printf("Duration: %v\n", elapsed)
	fmt.Printf("Memory allocated: %.2f MB\n", float64(ms.TotalAlloc)/(1024*1024))
	fmt.Printf("GC runs: %d\n", ms.NumGC)
}

type refactorResult struct {
	filename string
	changes  int
	err      error
}

func refactorFile(filename, operation, oldName, newName string) refactorResult {
	result := refactorResult{filename: filename}

	// Read file
	f, err := os.Open(filename)
	if err != nil {
		result.err = err
		return result
	}
	defer f.Close()

	var lines []string
	scanner := bufio.NewScanner(f)
	for scanner.Scan() {
		lines = append(lines, scanner.Text())
	}

	if err := scanner.Err(); err != nil {
		result.err = err
		return result
	}

	// Apply refactoring
	switch operation {
	case "rename":
		result.changes = renameVariable(lines, oldName, newName)
	case "add-comments":
		result.changes = addComments(lines)
	case "format":
		result.changes = formatCode(lines)
	default:
		result.err = fmt.Errorf("unknown operation: %s", operation)
		return result
	}

	// Write back if changes were made
	if result.changes > 0 {
		output := strings.Join(lines, "\n") + "\n"
		if err := os.WriteFile(filename, []byte(output), 0644); err != nil {
			result.err = err
		}
	}

	return result
}

func renameVariable(lines []string, oldName, newName string) int {
	changes := 0
	for i, line := range lines {
		if strings.Contains(line, oldName) {
			lines[i] = strings.ReplaceAll(line, oldName, newName)
			changes++
		}
	}
	return changes
}

func addComments(lines []string) int {
	changes := 0
	for i, line := range lines {
		trimmed := strings.TrimSpace(line)
		if strings.HasPrefix(trimmed, "func ") && i > 0 {
			prevLine := strings.TrimSpace(lines[i-1])
			if !strings.HasPrefix(prevLine, "//") {
				// Add comment before function
				indent := strings.Repeat("\t", countLeadingTabs(line))
				lines[i] = fmt.Sprintf("%s// %s\n%s", indent, extractFuncName(trimmed), line)
				changes++
			}
		}
	}
	return changes
}

func formatCode(lines []string) int {
	changes := 0
	for i, line := range lines {
		// Remove trailing whitespace
		trimmed := strings.TrimRight(line, " \t")
		if trimmed != line {
			lines[i] = trimmed
			changes++
		}
	}
	return changes
}

func countLeadingTabs(s string) int {
	count := 0
	for _, r := range s {
		if r == '\t' {
			count++
		} else {
			break
		}
	}
	return count
}

func extractFuncName(line string) string {
	parts := strings.Fields(line)
	if len(parts) >= 2 {
		name := parts[1]
		if idx := strings.Index(name, "("); idx > 0 {
			return name[:idx] + " function"
		}
	}
	return "Function"
}
