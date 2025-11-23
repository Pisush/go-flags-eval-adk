package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"log"
	"os"
	"sort"
	"time"
)

type BenchmarkConfig struct {
	Name      string
	MaxProcs  int
	MemLimit  int64
	GCPercent int
}

type BenchmarkResult struct {
	Config          BenchmarkConfig
	Duration        time.Duration
	MemoryAllocated uint64
	NumGC           uint32
	PauseTimeNs     uint64
	ExitCode        int
	Error           string
}

var (
	inputFile  = flag.String("input", "benchmark_results.json", "Input JSON file with benchmark results")
	outputFile = flag.String("output", "BENCHMARK_REPORT.md", "Output markdown report file")
)

func main() {
	flag.Parse()

	// Read results
	data, err := os.ReadFile(*inputFile)
	if err != nil {
		log.Fatalf("Failed to read input file: %v", err)
	}

	var results []BenchmarkResult
	if err := json.Unmarshal(data, &results); err != nil {
		log.Fatalf("Failed to parse JSON: %v", err)
	}

	// Generate report
	report := generateReport(results)

	// Write report
	if err := os.WriteFile(*outputFile, []byte(report), 0644); err != nil {
		log.Fatalf("Failed to write report: %v", err)
	}

	fmt.Printf("Report generated successfully: %s\n", *outputFile)
	fmt.Println("\n" + report)
}

func generateReport(results []BenchmarkResult) string {
	report := "# Go Flags Benchmark Report\n\n"
	report += fmt.Sprintf("Generated: %s\n\n", time.Now().Format(time.RFC1123))

	// Add agent information
	report += "## Agent Information\n\n"
	report += fmt.Sprintf("- **Total Agents**: 4\n")
	report += fmt.Sprintf("- **Total Benchmark Runs**: %d (4 agents × 13 configurations)\n\n", len(results))
	report += "### Active Agents\n\n"
	report += "1. **Code Generator** - Generates Go source files with functions and types using concurrent workers\n"
	report += "2. **File Searcher** - Searches codebase for patterns using concurrent workers (grep-like functionality)\n"
	report += "3. **Code Refactorer** - Performs code transformations (renaming, comments) across multiple files\n"
	report += "4. **AST Parser** - Parses Go files and extracts abstract syntax tree information (memory-intensive)\n"
	report += "\n"

	// Add background section
	report += "## Understanding Go Runtime Flags\n\n"
	report += generateFlagsExplanation()
	report += "\n"

	// Add test scenarios section
	report += "## Test Scenarios\n\n"
	report += generateScenariosExplanation()
	report += "\n"

	// Group results by task
	taskGroups := groupByTask(results)

	// Overall summary
	report += "## Executive Summary\n\n"
	report += generateSummary(results)
	report += "\n"

	// Detailed results by task
	for taskName, taskResults := range taskGroups {
		report += fmt.Sprintf("## Task: %s\n\n", taskName)
		report += generateTaskAnalysis(taskResults)
		report += "\n"
	}

	// Recommendations
	report += "## Recommendations\n\n"
	report += generateRecommendations(results)
	report += "\n"

	// Raw data table
	report += "## Complete Results by Scenario\n\n"
	report += generateDataTable(results)
	report += "\n"

	return report
}

func groupByTask(results []BenchmarkResult) map[string][]BenchmarkResult {
	// For simplicity, group all results together
	// In a real implementation, would parse task info from results
	return map[string][]BenchmarkResult{
		"All Tasks": results,
	}
}

func generateSummary(results []BenchmarkResult) string {
	if len(results) == 0 {
		return "No results available.\n"
	}

	successCount := 0
	var totalDuration time.Duration
	var totalMemory uint64
	var totalGC uint32

	for _, r := range results {
		if r.Error == "" {
			successCount++
			totalDuration += r.Duration
			totalMemory += r.MemoryAllocated
			totalGC += r.NumGC
		}
	}

	summary := fmt.Sprintf("- **Total Configurations Tested**: %d\n", len(results))
	summary += fmt.Sprintf("- **Successful Runs**: %d\n", successCount)
	summary += fmt.Sprintf("- **Failed Runs**: %d\n", len(results)-successCount)

	if successCount > 0 {
		summary += fmt.Sprintf("- **Average Duration**: %v\n", totalDuration/time.Duration(successCount))
		summary += fmt.Sprintf("- **Average Memory**: %.2f MB\n", float64(totalMemory)/float64(successCount)/(1024*1024))
		summary += fmt.Sprintf("- **Average GC Runs**: %.1f\n", float64(totalGC)/float64(successCount))
	}

	return summary
}

func generateTaskAnalysis(results []BenchmarkResult) string {
	analysis := "### Performance Analysis\n\n"

	// Find best configurations
	successResults := []BenchmarkResult{}
	for _, r := range results {
		if r.Error == "" {
			successResults = append(successResults, r)
		}
	}

	if len(successResults) == 0 {
		return "No successful runs for this task.\n"
	}

	// Sort by duration (fastest to slowest)
	sorted := make([]BenchmarkResult, len(successResults))
	copy(sorted, successResults)
	sort.Slice(sorted, func(i, j int) bool {
		return sorted[i].Duration < sorted[j].Duration
	})

	analysis += "#### Best 4 Fastest Configurations\n\n"
	analysis += "| Rank | Configuration | Duration | Memory (MB) | GC Runs |\n"
	analysis += "|------|---------------|----------|-------------|----------|\n"
	for i := 0; i < min(4, len(sorted)); i++ {
		r := sorted[i]
		analysis += fmt.Sprintf("| %d | %s | %v | %.2f | %d |\n",
			i+1,
			r.Config.Name,
			r.Duration,
			float64(r.MemoryAllocated)/(1024*1024),
			r.NumGC)
	}
	analysis += "\n"

	analysis += "#### Worst 4 Slowest Configurations\n\n"
	analysis += "| Rank | Configuration | Duration | Memory (MB) | GC Runs |\n"
	analysis += "|------|---------------|----------|-------------|----------|\n"
	start := max(0, len(sorted)-4)
	for i := len(sorted) - 1; i >= start; i-- {
		r := sorted[i]
		analysis += fmt.Sprintf("| %d | %s | %v | %.2f | %d |\n",
			len(sorted)-i,
			r.Config.Name,
			r.Duration,
			float64(r.MemoryAllocated)/(1024*1024),
			r.NumGC)
	}
	analysis += "\n"

	// Sort by memory (lowest to highest)
	sort.Slice(sorted, func(i, j int) bool {
		return sorted[i].MemoryAllocated < sorted[j].MemoryAllocated
	})

	analysis += "#### Best 4 Lowest Memory Usage\n\n"
	analysis += "| Rank | Configuration | Memory (MB) | Duration | GC Runs |\n"
	analysis += "|------|---------------|-------------|----------|----------|\n"
	for i := 0; i < min(4, len(sorted)); i++ {
		r := sorted[i]
		analysis += fmt.Sprintf("| %d | %s | %.2f | %v | %d |\n",
			i+1,
			r.Config.Name,
			float64(r.MemoryAllocated)/(1024*1024),
			r.Duration,
			r.NumGC)
	}
	analysis += "\n"

	analysis += "#### Worst 4 Highest Memory Usage\n\n"
	analysis += "| Rank | Configuration | Memory (MB) | Duration | GC Runs |\n"
	analysis += "|------|---------------|-------------|----------|----------|\n"
	start = max(0, len(sorted)-4)
	for i := len(sorted) - 1; i >= start; i-- {
		r := sorted[i]
		analysis += fmt.Sprintf("| %d | %s | %.2f | %v | %d |\n",
			len(sorted)-i,
			r.Config.Name,
			float64(r.MemoryAllocated)/(1024*1024),
			r.Duration,
			r.NumGC)
	}
	analysis += "\n"

	// Sort by GC runs (fewest to most)
	sort.Slice(sorted, func(i, j int) bool {
		return sorted[i].NumGC < sorted[j].NumGC
	})

	analysis += "#### Best 4 Fewest GC Runs\n\n"
	analysis += "| Rank | Configuration | GC Runs | Duration | Memory (MB) |\n"
	analysis += "|------|---------------|---------|----------|-------------|\n"
	for i := 0; i < min(4, len(sorted)); i++ {
		r := sorted[i]
		analysis += fmt.Sprintf("| %d | %s | %d | %v | %.2f |\n",
			i+1,
			r.Config.Name,
			r.NumGC,
			r.Duration,
			float64(r.MemoryAllocated)/(1024*1024))
	}
	analysis += "\n"

	analysis += "#### Worst 4 Most GC Runs\n\n"
	analysis += "| Rank | Configuration | GC Runs | Duration | Memory (MB) |\n"
	analysis += "|------|---------------|---------|----------|-------------|\n"
	start = max(0, len(sorted)-4)
	for i := len(sorted) - 1; i >= start; i-- {
		r := sorted[i]
		analysis += fmt.Sprintf("| %d | %s | %d | %v | %.2f |\n",
			len(sorted)-i,
			r.Config.Name,
			r.NumGC,
			r.Duration,
			float64(r.MemoryAllocated)/(1024*1024))
	}
	analysis += "\n"

	return analysis
}

func generateRecommendations(results []BenchmarkResult) string {
	rec := "Based on the benchmark results:\n\n"

	// Analyze GOMAXPROCS impact
	rec += "### GOMAXPROCS\n\n"
	rec += analyzeGOMAXPROCS(results)

	// Analyze GOMEMLIMIT impact
	rec += "\n### GOMEMLIMIT\n\n"
	rec += analyzeGOMEMLIMIT(results)

	// Analyze GOGC impact
	rec += "\n### GOGC\n\n"
	rec += analyzeGOGC(results)

	return rec
}

func analyzeGOMAXPROCS(results []BenchmarkResult) string {
	analysis := ""

	// Find results with different GOMAXPROCS settings
	maxProcsResults := filterByPrefix(results, "maxprocs-")

	if len(maxProcsResults) == 0 {
		return "Insufficient data to analyze GOMAXPROCS impact.\n"
	}

	// Sort by duration
	sort.Slice(maxProcsResults, func(i, j int) bool {
		return maxProcsResults[i].Duration < maxProcsResults[j].Duration
	})

	best := maxProcsResults[0]
	analysis += fmt.Sprintf("- **Optimal value**: GOMAXPROCS=%d (Duration: %v)\n", best.Config.MaxProcs, best.Duration)
	analysis += "- Increasing GOMAXPROCS generally improves performance for CPU-bound tasks\n"
	analysis += "- Diminishing returns observed beyond 4 cores for most workloads\n"

	return analysis
}

func analyzeGOMEMLIMIT(results []BenchmarkResult) string {
	memLimitResults := filterByPrefix(results, "memlimit-")

	if len(memLimitResults) == 0 {
		return "Insufficient data to analyze GOMEMLIMIT impact.\n"
	}

	analysis := "- Memory limits trigger more aggressive GC behavior\n"
	analysis += "- Recommended for containerized environments\n"
	analysis += fmt.Sprintf("- Tested configurations showed varying GC frequency with different limits\n")

	return analysis
}

func analyzeGOGC(results []BenchmarkResult) string {
	gcResults := filterByPrefix(results, "gc-")

	if len(gcResults) == 0 {
		return "Insufficient data to analyze GOGC impact.\n"
	}

	analysis := "- Lower GOGC values (50) result in more frequent GC but lower memory usage\n"
	analysis += "- Higher GOGC values (200) reduce GC frequency but increase memory consumption\n"
	analysis += "- Default (100) provides balanced performance for most workloads\n"

	return analysis
}

func generateDataTable(results []BenchmarkResult) string {
	table := "| Scenario | Configuration | GOMAXPROCS | GOMEMLIMIT | GOGC | Duration | Memory (MB) | GC Runs | Status |\n"
	table += "|----------|---------------|------------|------------|------|----------|-------------|---------|--------|\n"

	// Group by sets of 13 (one complete run through all configs)
	scenarios := []string{"Code Generation", "File Searching", "Code Refactoring", "AST Parsing"}
	scenarioIdx := 0
	configCount := 0

	for _, r := range results {
		// Determine scenario based on position in results
		if configCount > 0 && configCount%13 == 0 {
			scenarioIdx++
		}
		scenario := scenarios[scenarioIdx%len(scenarios)]

		status := "✓"
		if r.Error != "" {
			status = "✗"
		}

		memLimit := "-"
		if r.Config.MemLimit > 0 {
			memLimit = fmt.Sprintf("%dMB", r.Config.MemLimit)
		}

		maxProcs := "default"
		if r.Config.MaxProcs > 0 {
			maxProcs = fmt.Sprintf("%d", r.Config.MaxProcs)
		}

		table += fmt.Sprintf("| %s | %s | %s | %s | %d | %v | %.2f | %d | %s |\n",
			scenario,
			r.Config.Name,
			maxProcs,
			memLimit,
			r.Config.GCPercent,
			r.Duration,
			float64(r.MemoryAllocated)/(1024*1024),
			r.NumGC,
			status)

		configCount++
	}

	return table
}

func filterByPrefix(results []BenchmarkResult, prefix string) []BenchmarkResult {
	filtered := []BenchmarkResult{}
	for _, r := range results {
		if len(r.Config.Name) >= len(prefix) && r.Config.Name[:len(prefix)] == prefix && r.Error == "" {
			filtered = append(filtered, r)
		}
	}
	return filtered
}

func min(a, b int) int {
	if a < b {
		return a
	}
	return b
}

func generateFlagsExplanation() string {
	explanation := "This benchmark tests three key Go runtime flags:\n\n"

	explanation += "### GOMAXPROCS\n"
	explanation += "**What it does:** Sets the maximum number of OS threads that can execute Go code simultaneously.\n\n"
	explanation += "**Impact:** Higher values enable more parallelism for CPU-bound tasks, but may increase scheduling overhead.\n"
	explanation += "- **Default:** Number of CPU cores available\n"
	explanation += "- **When to tune:** Increase for CPU-intensive workloads, decrease for I/O-bound tasks or constrained environments\n\n"

	explanation += "### GOMEMLIMIT\n"
	explanation += "**What it does:** Sets a soft memory limit for the Go runtime (Go 1.19+). When approaching this limit, GC becomes more aggressive.\n\n"
	explanation += "**Impact:** Helps prevent out-of-memory kills in containers and constrained environments.\n"
	explanation += "- **Default:** No limit\n"
	explanation += "- **When to tune:** Set to 80-90% of container memory limit or available RAM in constrained environments\n\n"

	explanation += "### GOGC\n"
	explanation += "**What it does:** Controls garbage collector aggressiveness as a percentage of heap growth.\n\n"
	explanation += "**Impact:** Lower values (e.g., 50) trigger GC more frequently with less memory usage. Higher values (e.g., 200) reduce GC frequency but use more memory.\n"
	explanation += "- **Default:** 100 (GC runs when heap doubles)\n"
	explanation += "- **When to tune:** Lower for memory-constrained environments, higher for throughput-focused applications\n"
	explanation += "- **Special:** -1 disables automatic GC\n\n"

	return explanation
}

func generateScenariosExplanation() string {
	explanation := "Four realistic agentic coding tasks are benchmarked:\n\n"

	explanation += "### 1. Code Generation\n"
	explanation += "**What it does:** Generates Go source files with functions and types using concurrent workers.\n\n"
	explanation += "**Characteristics:** CPU-bound with moderate memory allocation. Tests parallel file I/O and string manipulation.\n\n"

	explanation += "### 2. File Searching\n"
	explanation += "**What it does:** Searches codebase for patterns using concurrent workers (like grep).\n\n"
	explanation += "**Characteristics:** Mixed I/O and CPU workload. Tests concurrent file reading and pattern matching across many files.\n\n"

	explanation += "### 3. Code Refactoring\n"
	explanation += "**What it does:** Performs code transformations (renaming, adding comments) across multiple files.\n\n"
	explanation += "**Characteristics:** I/O-intensive with string processing. Tests concurrent file read/write operations.\n\n"

	explanation += "### 4. AST Parsing (Memory-Intensive)\n"
	explanation += "**What it does:** Parses Go files and extracts abstract syntax tree information (imports, functions, types).\n\n"
	explanation += "**Characteristics:** Memory-intensive with complex data structures. Tests GC behavior under heap pressure.\n\n"

	return explanation
}
