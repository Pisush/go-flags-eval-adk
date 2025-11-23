package main

import (
	"context"
	"encoding/json"
	"flag"
	"fmt"
	"log"
	"os"
	"os/exec"
	"runtime"
	"runtime/debug"
	"time"

	"github.com/natalie/go-flags-eval/internal/agentmetrics"
)

// BenchmarkConfig defines a set of Go runtime flags to test
type BenchmarkConfig struct {
	Name      string
	MaxProcs  int
	MemLimit  int64 // in MB
	GCPercent int
}

// BenchmarkResult stores the results of a benchmark run
type BenchmarkResult struct {
	Config           BenchmarkConfig
	Duration         time.Duration
	MemoryAllocated  uint64
	NumGC            uint32
	PauseTimeNs      uint64
	ExitCode         int
	Error            string
}

// AgentTask represents a task for an agent to perform
type AgentTask struct {
	Name        string
	Command     string
	Args        []string
	Description string
}

var (
	outputFile = flag.String("output", "benchmark_results.json", "Output file for benchmark results")
	taskName   = flag.String("task", "all", "Specific task to run (all, code-gen, file-search, code-analysis)")
)

func main() {
	flag.Parse()

	ctx := context.Background()

	// Define benchmark configurations
	configs := []BenchmarkConfig{
		{"default", 0, 0, 100},
		{"maxprocs-1", 1, 0, 100},
		{"maxprocs-2", 2, 0, 100},
		{"maxprocs-4", 4, 0, 100},
		{"maxprocs-8", 8, 0, 100},
		{"memlimit-256", 0, 256, 100},
		{"memlimit-512", 0, 512, 100},
		{"memlimit-1024", 0, 1024, 100},
		{"gc-50", 0, 0, 50},
		{"gc-200", 0, 0, 200},
		{"gc-off", 0, 0, -1},
		{"constrained", 2, 256, 50},
		{"performance", 8, 2048, 200},
	}

	// Define agent tasks
	tasks := []AgentTask{
		{
			Name:        "code-gen",
			Command:     "go",
			Args:        []string{"run", "./cmd/agents/code_generator", "-files=10", "-lines=100"},
			Description: "Generate 10 Go files with 100 lines each",
		},
		{
			Name:        "file-search",
			Command:     "go",
			Args:        []string{"run", "./cmd/agents/file_searcher", "-pattern=func", "-dir=./testdata"},
			Description: "Search for 'func' pattern in test data",
		},
		{
			Name:        "refactor",
			Command:     "go",
			Args:        []string{"run", "./cmd/agents/refactor", "-target=./testdata", "-operation=rename"},
			Description: "Rename variables across multiple files",
		},
		{
			Name:        "ast-parser",
			Command:     "go",
			Args:        []string{"run", "./cmd/agents/ast_parser", "-target=./testdata"},
			Description: "Parse Go files and extract AST information (memory-intensive)",
		},
	}

	// Filter tasks if specific task requested
	if *taskName != "all" {
		filtered := []AgentTask{}
		for _, task := range tasks {
			if task.Name == *taskName {
				filtered = append(filtered, task)
			}
		}
		if len(filtered) == 0 {
			log.Fatalf("Unknown task: %s", *taskName)
		}
		tasks = filtered
	}

	// Run benchmarks
	results := []BenchmarkResult{}
	for _, task := range tasks {
		fmt.Printf("\n=== Running Task: %s ===\n", task.Name)
		fmt.Printf("Description: %s\n\n", task.Description)

		for _, cfg := range configs {
			fmt.Printf("Testing configuration: %s... ", cfg.Name)
			result := runBenchmark(ctx, task, cfg)
			results = append(results, result)

			if result.Error != "" {
				fmt.Printf("ERROR: %s\n", result.Error)
			} else {
				fmt.Printf("Duration: %v, Memory: %.2f MB, GC runs: %d\n",
					result.Duration,
					float64(result.MemoryAllocated)/(1024*1024),
					result.NumGC)
			}
		}
	}

	// Save results to JSON
	if err := saveResults(*outputFile, results); err != nil {
		log.Fatalf("Failed to save results: %v", err)
	}

	fmt.Printf("\n\nResults saved to: %s\n", *outputFile)
	printSummary(results)
}

func runBenchmark(ctx context.Context, task AgentTask, cfg BenchmarkConfig) BenchmarkResult {
	result := BenchmarkResult{
		Config: cfg,
	}

	// Create temporary file for metrics
	metricsFile, err := os.CreateTemp("", "agent-metrics-*.json")
	if err != nil {
		result.Error = fmt.Sprintf("Failed to create metrics file: %v", err)
		return result
	}
	metricsPath := metricsFile.Name()
	metricsFile.Close()
	defer os.Remove(metricsPath)

	// Prepare environment
	env := os.Environ()
	if cfg.MaxProcs > 0 {
		env = append(env, fmt.Sprintf("GOMAXPROCS=%d", cfg.MaxProcs))
	}
	if cfg.MemLimit > 0 {
		env = append(env, fmt.Sprintf("GOMEMLIMIT=%dMiB", cfg.MemLimit))
	}
	if cfg.GCPercent != 100 {
		env = append(env, fmt.Sprintf("GOGC=%d", cfg.GCPercent))
	}

	// Add metrics output flag to args
	args := append(task.Args, fmt.Sprintf("-metrics-output=%s", metricsPath))

	// Run command
	cmd := exec.CommandContext(ctx, task.Command, args...)
	cmd.Env = env
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr

	startTime := time.Now()
	err = cmd.Run()
	result.Duration = time.Since(startTime)

	if err != nil {
		result.Error = err.Error()
		if exitErr, ok := err.(*exec.ExitError); ok {
			result.ExitCode = exitErr.ExitCode()
		}
	}

	// Read metrics from agent
	metrics, err := agentmetrics.ReadFromFile(metricsPath)
	if err != nil {
		// Fall back to duration from benchmark harness
		log.Printf("Warning: Could not read agent metrics: %v", err)
	} else {
		// Use metrics from the actual agent process
		result.MemoryAllocated = metrics.MemoryAllocated
		result.NumGC = metrics.NumGC
		result.PauseTimeNs = metrics.PauseTimeNs
	}

	return result
}

func saveResults(filename string, results []BenchmarkResult) error {
	data, err := json.MarshalIndent(results, "", "  ")
	if err != nil {
		return fmt.Errorf("failed to marshal results: %w", err)
	}

	return os.WriteFile(filename, data, 0644)
}

func printSummary(results []BenchmarkResult) {
	fmt.Println("\n=== Summary ===")

	// Find best configurations
	var fastestDuration time.Duration = time.Hour
	var lowestMemory uint64 = ^uint64(0)
	var fewestGC uint32 = ^uint32(0)

	var fastestConfig, lowestMemConfig, fewestGCConfig string

	for _, r := range results {
		if r.Error == "" {
			if r.Duration < fastestDuration {
				fastestDuration = r.Duration
				fastestConfig = r.Config.Name
			}
			if r.MemoryAllocated < lowestMemory {
				lowestMemory = r.MemoryAllocated
				lowestMemConfig = r.Config.Name
			}
			if r.NumGC < fewestGC {
				fewestGC = r.NumGC
				fewestGCConfig = r.Config.Name
			}
		}
	}

	fmt.Printf("Fastest execution: %s (%v)\n", fastestConfig, fastestDuration)
	fmt.Printf("Lowest memory: %s (%.2f MB)\n", lowestMemConfig, float64(lowestMemory)/(1024*1024))
	fmt.Printf("Fewest GC runs: %s (%d runs)\n", fewestGCConfig, fewestGC)
}

func init() {
	// Set reasonable defaults for the benchmark runner itself
	runtime.GOMAXPROCS(runtime.NumCPU())
	debug.SetGCPercent(100)
}
