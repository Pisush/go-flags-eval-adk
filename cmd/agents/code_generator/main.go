package main

import (
	"flag"
	"fmt"
	"log"
	"math/rand"
	"os"
	"path/filepath"
	"runtime"
	"runtime/debug"
	"time"
)

var (
	numFiles = flag.Int("files", 10, "Number of files to generate")
	numLines = flag.Int("lines", 100, "Number of lines per file")
	outputDir = flag.String("output", "./generated", "Output directory")
)

func main() {
	flag.Parse()

	start := time.Now()

	// Create output directory
	if err := os.MkdirAll(*outputDir, 0755); err != nil {
		log.Fatalf("Failed to create output directory: %v", err)
	}

	// Report configuration
	fmt.Printf("Code Generator Agent\n")
	fmt.Printf("====================\n")
	fmt.Printf("GOMAXPROCS: %d\n", runtime.GOMAXPROCS(-1))
	fmt.Printf("GOGC: %d\n", debug.SetGCPercent(-1))
	debug.SetGCPercent(debug.SetGCPercent(-1)) // Restore
	fmt.Printf("Files to generate: %d\n", *numFiles)
	fmt.Printf("Lines per file: %d\n", *numLines)
	fmt.Printf("\n")

	// Generate files concurrently
	type result struct {
		filename string
		err      error
	}

	results := make(chan result, *numFiles)
	sem := make(chan struct{}, runtime.GOMAXPROCS(-1))

	for i := 0; i < *numFiles; i++ {
		go func(fileNum int) {
			sem <- struct{}{}
			defer func() { <-sem }()

			filename := filepath.Join(*outputDir, fmt.Sprintf("generated_%d.go", fileNum))
			err := generateGoFile(filename, *numLines)
			results <- result{filename, err}
		}(i)
	}

	// Collect results
	successCount := 0
	for i := 0; i < *numFiles; i++ {
		res := <-results
		if res.err != nil {
			log.Printf("Error generating %s: %v", res.filename, res.err)
		} else {
			successCount++
		}
	}

	elapsed := time.Since(start)

	// Print statistics
	var ms runtime.MemStats
	runtime.ReadMemStats(&ms)

	fmt.Printf("\nResults:\n")
	fmt.Printf("========\n")
	fmt.Printf("Files generated: %d/%d\n", successCount, *numFiles)
	fmt.Printf("Duration: %v\n", elapsed)
	fmt.Printf("Memory allocated: %.2f MB\n", float64(ms.TotalAlloc)/(1024*1024))
	fmt.Printf("GC runs: %d\n", ms.NumGC)
	fmt.Printf("Goroutines: %d\n", runtime.NumGoroutine())
}

func generateGoFile(filename string, lines int) error {
	f, err := os.Create(filename)
	if err != nil {
		return err
	}
	defer f.Close()

	packageName := filepath.Base(filepath.Dir(filename))
	if packageName == "." || packageName == "/" {
		packageName = "main"
	}

	// Write package declaration
	fmt.Fprintf(f, "package %s\n\n", packageName)
	fmt.Fprintf(f, "import (\n\t\"fmt\"\n\t\"math/rand\"\n)\n\n")

	// Generate functions
	numFuncs := lines / 10
	if numFuncs < 1 {
		numFuncs = 1
	}

	for i := 0; i < numFuncs; i++ {
		funcName := fmt.Sprintf("Function%d", i)
		fmt.Fprintf(f, "// %s performs operation %d\n", funcName, i)
		fmt.Fprintf(f, "func %s(x, y int) int {\n", funcName)

		// Generate function body
		linesInFunc := lines / numFuncs
		for j := 0; j < linesInFunc; j++ {
			operation := rand.Intn(4)
			switch operation {
			case 0:
				fmt.Fprintf(f, "\tx = x + y\n")
			case 1:
				fmt.Fprintf(f, "\ty = y * 2\n")
			case 2:
				fmt.Fprintf(f, "\tx = x - y\n")
			case 3:
				fmt.Fprintf(f, "\tfmt.Printf(\"Calculating: %%d\\n\", x+y)\n")
			}
		}

		fmt.Fprintf(f, "\treturn x + y\n")
		fmt.Fprintf(f, "}\n\n")
	}

	return nil
}
