package main

import (
	"flag"
	"fmt"
	"go/ast"
	"go/parser"
	"go/token"
	"log"
	"os"
	"path/filepath"
	"runtime"
	"runtime/debug"
	"strings"
	"sync"
	"time"

	"github.com/natalie/go-flags-eval/internal/agentmetrics"
)

var (
	target        = flag.String("target", "./testdata", "Target directory to parse")
	findImports   = flag.Bool("imports", true, "Find all imports")
	findFuncs     = flag.Bool("funcs", true, "Find all functions")
	findTypes     = flag.Bool("types", true, "Find all type definitions")
	metricsOutput = flag.String("metrics-output", "", "File to write performance metrics (JSON)")
)

type ParsedFile struct {
	Path    string
	Imports []string
	Funcs   []string
	Types   []string
}

func main() {
	flag.Parse()

	start := time.Now()

	// Report configuration
	fmt.Printf("AST Parser Agent (Memory-Intensive)\n")
	fmt.Printf("====================================\n")
	fmt.Printf("GOMAXPROCS: %d\n", runtime.GOMAXPROCS(-1))
	gcVal := debug.SetGCPercent(-1)
	debug.SetGCPercent(gcVal)
	fmt.Printf("GOGC: %d\n", gcVal)
	fmt.Printf("Target: %s\n", *target)
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

	fmt.Printf("Found %d Go files to parse\n\n", len(files))

	// Parse files concurrently
	var wg sync.WaitGroup
	results := make(chan ParsedFile, len(files))
	sem := make(chan struct{}, runtime.GOMAXPROCS(-1))

	for _, file := range files {
		wg.Add(1)
		go func(filename string) {
			defer wg.Done()
			sem <- struct{}{}
			defer func() { <-sem }()

			parsed := parseFile(filename)
			if parsed != nil {
				results <- *parsed
			}
		}(file)
	}

	go func() {
		wg.Wait()
		close(results)
	}()

	// Collect results
	allParsed := []ParsedFile{}
	totalImports := 0
	totalFuncs := 0
	totalTypes := 0

	for parsed := range results {
		allParsed = append(allParsed, parsed)
		totalImports += len(parsed.Imports)
		totalFuncs += len(parsed.Funcs)
		totalTypes += len(parsed.Types)
	}

	elapsed := time.Since(start)

	// Collect statistics
	var ms runtime.MemStats
	runtime.ReadMemStats(&ms)

	// Print statistics
	fmt.Printf("\nResults:\n")
	fmt.Printf("========\n")
	fmt.Printf("Files parsed: %d\n", len(allParsed))
	fmt.Printf("Total imports: %d\n", totalImports)
	fmt.Printf("Total functions: %d\n", totalFuncs)
	fmt.Printf("Total types: %d\n", totalTypes)
	fmt.Printf("Duration: %v\n", elapsed)
	fmt.Printf("Memory allocated: %.2f MB\n", float64(ms.TotalAlloc)/(1024*1024))
	fmt.Printf("Heap allocated: %.2f MB\n", float64(ms.HeapAlloc)/(1024*1024))
	fmt.Printf("GC runs: %d\n", ms.NumGC)
	fmt.Printf("Goroutines: %d\n", runtime.NumGoroutine())

	// Write metrics to file if requested
	if *metricsOutput != "" {
		metrics := &agentmetrics.Metrics{
			Duration:        elapsed,
			MemoryAllocated: ms.TotalAlloc,
			HeapAllocated:   ms.HeapAlloc,
			NumGC:           ms.NumGC,
			PauseTimeNs:     ms.PauseTotalNs,
			Goroutines:      runtime.NumGoroutine(),
			FilesProcessed:  len(allParsed),
			Custom: map[string]any{
				"total_imports":   totalImports,
				"total_functions": totalFuncs,
				"total_types":     totalTypes,
			},
		}

		if err := metrics.WriteToFile(*metricsOutput); err != nil {
			log.Printf("Failed to write metrics: %v", err)
		}
	}
}

func parseFile(filename string) *ParsedFile {
	fset := token.NewFileSet()
	node, err := parser.ParseFile(fset, filename, nil, parser.ParseComments)
	if err != nil {
		log.Printf("Error parsing %s: %v", filename, err)
		return nil
	}

	result := &ParsedFile{
		Path: filename,
	}

	// Extract imports
	if *findImports {
		for _, imp := range node.Imports {
			if imp.Path != nil {
				result.Imports = append(result.Imports, imp.Path.Value)
			}
		}
	}

	// Extract functions and types
	ast.Inspect(node, func(n ast.Node) bool {
		switch x := n.(type) {
		case *ast.FuncDecl:
			if *findFuncs && x.Name != nil {
				result.Funcs = append(result.Funcs, x.Name.Name)
			}
		case *ast.TypeSpec:
			if *findTypes && x.Name != nil {
				result.Types = append(result.Types, x.Name.Name)
			}
		}
		return true
	})

	return result
}
