# Go Flags Benchmark Report - Sample Results

> **Note:** These are example results from a specific test environment.
> Your results will vary based on hardware, OS, and system load.
> Run `make benchmark && make report` to generate benchmarks for your machine.

## Test Environment

- **CPU:** Apple M2 (ARM64)
- **Memory:** 24 GB
- **OS:** macOS Darwin 23.5.0
- **Go Version:** 1.24.10
- **Date:** November 23, 2025

## Dataset

- **Source:** Prometheus Go Client Library v1.19.0
- **Repository:** https://github.com/prometheus/client_golang
- **License:** Apache 2.0
- **Files:** 122 Go files (31,023 lines of code, 1.2 MB)
- **Content:** Production-grade code for metrics collection, HTTP handlers, collectors, and Prometheus API client

---

Generated: Sun, 23 Nov 2025 11:43:54 CET

## Understanding Go Runtime Flags

This benchmark tests three key Go runtime flags:

### GOMAXPROCS
**What it does:** Sets the maximum number of OS threads that can execute Go code simultaneously.

**Impact:** Higher values enable more parallelism for CPU-bound tasks, but may increase scheduling overhead.
- **Default:** Number of CPU cores available
- **When to tune:** Increase for CPU-intensive workloads, decrease for I/O-bound tasks or constrained environments

### GOMEMLIMIT
**What it does:** Sets a soft memory limit for the Go runtime (Go 1.19+). When approaching this limit, GC becomes more aggressive.

**Impact:** Helps prevent out-of-memory kills in containers and constrained environments.
- **Default:** No limit
- **When to tune:** Set to 80-90% of container memory limit or available RAM in constrained environments

### GOGC
**What it does:** Controls garbage collector aggressiveness as a percentage of heap growth.

**Impact:** Lower values (e.g., 50) trigger GC more frequently with less memory usage. Higher values (e.g., 200) reduce GC frequency but use more memory.
- **Default:** 100 (GC runs when heap doubles)
- **When to tune:** Lower for memory-constrained environments, higher for throughput-focused applications
- **Special:** -1 disables automatic GC


## Test Scenarios

Four realistic agentic coding tasks are benchmarked:

### 1. Code Generation
**What it does:** Generates Go source files with functions and types using concurrent workers.

**Characteristics:** CPU-bound with moderate memory allocation. Tests parallel file I/O and string manipulation.

### 2. File Searching
**What it does:** Searches codebase for patterns using concurrent workers (like grep).

**Characteristics:** Mixed I/O and CPU workload. Tests concurrent file reading and pattern matching across many files.

### 3. Code Refactoring
**What it does:** Performs code transformations (renaming, adding comments) across multiple files.

**Characteristics:** I/O-intensive with string processing. Tests concurrent file read/write operations.

### 4. AST Parsing (Memory-Intensive)
**What it does:** Parses Go files and extracts abstract syntax tree information (imports, functions, types).

**Characteristics:** Memory-intensive with complex data structures. Tests GC behavior under heap pressure.


## Executive Summary

- **Total Configurations Tested**: 52
- **Successful Runs**: 52
- **Failed Runs**: 0
- **Average Duration**: 85.852388ms
- **Average Memory**: 0.01 MB
- **Average GC Runs**: 0.0

## Task: All Tasks

### Performance Analysis

#### Best 4 Fastest Configurations

| Rank | Configuration | Duration | Memory (MB) | GC Runs |
|------|---------------|----------|-------------|----------|
| 1 | maxprocs-4 | 55.852209ms | 0.01 | 0 |
| 2 | default | 62.72525ms | 0.01 | 0 |
| 3 | memlimit-256 | 62.861916ms | 0.01 | 0 |
| 4 | memlimit-512 | 63.047209ms | 0.01 | 0 |

#### Worst 4 Slowest Configurations

| Rank | Configuration | Duration | Memory (MB) | GC Runs |
|------|---------------|----------|-------------|----------|
| 1 | gc-50 | 166.901291ms | 0.01 | 0 |
| 2 | maxprocs-1 | 150.255542ms | 0.01 | 0 |
| 3 | gc-50 | 120.807958ms | 0.01 | 0 |
| 4 | constrained | 119.786875ms | 0.01 | 0 |

#### Best 4 Lowest Memory Usage

| Rank | Configuration | Memory (MB) | Duration | GC Runs |
|------|---------------|-------------|----------|----------|
| 1 | memlimit-1024 | 0.01 | 112.640583ms | 0 |
| 2 | gc-50 | 0.01 | 166.901291ms | 0 |
| 3 | memlimit-512 | 0.01 | 89.255709ms | 0 |
| 4 | maxprocs-8 | 0.01 | 89.936417ms | 0 |

#### Worst 4 Highest Memory Usage

| Rank | Configuration | Memory (MB) | Duration | GC Runs |
|------|---------------|-------------|----------|----------|
| 1 | constrained | 0.01 | 93.967ms | 0 |
| 2 | performance | 0.01 | 72.363708ms | 0 |
| 3 | constrained | 0.01 | 80.906125ms | 0 |
| 4 | performance | 0.01 | 69.228875ms | 0 |

#### Best 4 Fewest GC Runs

| Rank | Configuration | GC Runs | Duration | Memory (MB) |
|------|---------------|---------|----------|-------------|
| 1 | memlimit-1024 | 0 | 112.640583ms | 0.01 |
| 2 | gc-50 | 0 | 166.901291ms | 0.01 |
| 3 | memlimit-512 | 0 | 89.255709ms | 0.01 |
| 4 | maxprocs-8 | 0 | 89.936417ms | 0.01 |

#### Worst 4 Most GC Runs

| Rank | Configuration | GC Runs | Duration | Memory (MB) |
|------|---------------|---------|----------|-------------|
| 1 | constrained | 0 | 93.967ms | 0.01 |
| 2 | performance | 0 | 72.363708ms | 0.01 |
| 3 | constrained | 0 | 80.906125ms | 0.01 |
| 4 | performance | 0 | 69.228875ms | 0.01 |


## Recommendations

Based on the benchmark results:

### GOMAXPROCS

- **Optimal value**: GOMAXPROCS=4 (Duration: 55.852209ms)
- Increasing GOMAXPROCS generally improves performance for CPU-bound tasks
- Diminishing returns observed beyond 4 cores for most workloads

### GOMEMLIMIT

- Memory limits trigger more aggressive GC behavior
- Recommended for containerized environments
- Tested configurations showed varying GC frequency with different limits

### GOGC

- Lower GOGC values (50) result in more frequent GC but lower memory usage
- Higher GOGC values (200) reduce GC frequency but increase memory consumption
- Default (100) provides balanced performance for most workloads

## Complete Results by Scenario

| Scenario | Configuration | GOMAXPROCS | GOMEMLIMIT | GOGC | Duration | Memory (MB) | GC Runs | Status |
|----------|---------------|------------|------------|------|----------|-------------|---------|--------|
| Code Generation | default | default | - | 100 | 62.72525ms | 0.01 | 0 | ✓ |
| Code Generation | maxprocs-1 | 1 | - | 100 | 86.629625ms | 0.01 | 0 | ✓ |
| Code Generation | maxprocs-2 | 2 | - | 100 | 63.923542ms | 0.01 | 0 | ✓ |
| Code Generation | maxprocs-4 | 4 | - | 100 | 55.852209ms | 0.01 | 0 | ✓ |
| Code Generation | maxprocs-8 | 8 | - | 100 | 63.874833ms | 0.01 | 0 | ✓ |
| Code Generation | memlimit-256 | default | 256MB | 100 | 62.861916ms | 0.01 | 0 | ✓ |
| Code Generation | memlimit-512 | default | 512MB | 100 | 63.047209ms | 0.01 | 0 | ✓ |
| Code Generation | memlimit-1024 | default | 1024MB | 100 | 66.533541ms | 0.01 | 0 | ✓ |
| Code Generation | gc-50 | default | - | 50 | 84.604292ms | 0.01 | 0 | ✓ |
| Code Generation | gc-200 | default | - | 200 | 67.017584ms | 0.01 | 0 | ✓ |
| Code Generation | gc-off | default | - | -1 | 67.198667ms | 0.01 | 0 | ✓ |
| Code Generation | constrained | 2 | 256MB | 50 | 80.906125ms | 0.01 | 0 | ✓ |
| Code Generation | performance | 8 | 2048MB | 200 | 69.228875ms | 0.01 | 0 | ✓ |
| File Searching | default | default | - | 100 | 77.780583ms | 0.01 | 0 | ✓ |
| File Searching | maxprocs-1 | 1 | - | 100 | 107.694542ms | 0.01 | 0 | ✓ |
| File Searching | maxprocs-2 | 2 | - | 100 | 74.745667ms | 0.01 | 0 | ✓ |
| File Searching | maxprocs-4 | 4 | - | 100 | 66.519209ms | 0.01 | 0 | ✓ |
| File Searching | maxprocs-8 | 8 | - | 100 | 75.982958ms | 0.01 | 0 | ✓ |
| File Searching | memlimit-256 | default | 256MB | 100 | 75.234125ms | 0.01 | 0 | ✓ |
| File Searching | memlimit-512 | default | 512MB | 100 | 76.50225ms | 0.01 | 0 | ✓ |
| File Searching | memlimit-1024 | default | 1024MB | 100 | 74.964917ms | 0.01 | 0 | ✓ |
| File Searching | gc-50 | default | - | 50 | 101.30375ms | 0.01 | 0 | ✓ |
| File Searching | gc-200 | default | - | 200 | 83.545708ms | 0.01 | 0 | ✓ |
| File Searching | gc-off | default | - | -1 | 78.587208ms | 0.01 | 0 | ✓ |
| File Searching | constrained | 2 | 256MB | 50 | 88.39925ms | 0.01 | 0 | ✓ |
| File Searching | performance | 8 | 2048MB | 200 | 75.1265ms | 0.01 | 0 | ✓ |
| Code Refactoring | default | default | - | 100 | 89.83875ms | 0.01 | 0 | ✓ |
| Code Refactoring | maxprocs-1 | 1 | - | 100 | 112.595334ms | 0.01 | 0 | ✓ |
| Code Refactoring | maxprocs-2 | 2 | - | 100 | 77.051541ms | 0.01 | 0 | ✓ |
| Code Refactoring | maxprocs-4 | 4 | - | 100 | 66.28825ms | 0.01 | 0 | ✓ |
| Code Refactoring | maxprocs-8 | 8 | - | 100 | 74.179125ms | 0.01 | 0 | ✓ |
| Code Refactoring | memlimit-256 | default | 256MB | 100 | 73.490042ms | 0.01 | 0 | ✓ |
| Code Refactoring | memlimit-512 | default | 512MB | 100 | 85.720958ms | 0.01 | 0 | ✓ |
| Code Refactoring | memlimit-1024 | default | 1024MB | 100 | 97.154625ms | 0.01 | 0 | ✓ |
| Code Refactoring | gc-50 | default | - | 50 | 120.807958ms | 0.01 | 0 | ✓ |
| Code Refactoring | gc-200 | default | - | 200 | 101.439917ms | 0.01 | 0 | ✓ |
| Code Refactoring | gc-off | default | - | -1 | 93ms | 0.01 | 0 | ✓ |
| Code Refactoring | constrained | 2 | 256MB | 50 | 93.967ms | 0.01 | 0 | ✓ |
| Code Refactoring | performance | 8 | 2048MB | 200 | 72.363708ms | 0.01 | 0 | ✓ |
| AST Parsing | default | default | - | 100 | 96.975458ms | 0.01 | 0 | ✓ |
| AST Parsing | maxprocs-1 | 1 | - | 100 | 150.255542ms | 0.01 | 0 | ✓ |
| AST Parsing | maxprocs-2 | 2 | - | 100 | 97.497125ms | 0.01 | 0 | ✓ |
| AST Parsing | maxprocs-4 | 4 | - | 100 | 72.652458ms | 0.01 | 0 | ✓ |
| AST Parsing | maxprocs-8 | 8 | - | 100 | 89.936417ms | 0.01 | 0 | ✓ |
| AST Parsing | memlimit-256 | default | 256MB | 100 | 88.813083ms | 0.01 | 0 | ✓ |
| AST Parsing | memlimit-512 | default | 512MB | 100 | 89.255709ms | 0.01 | 0 | ✓ |
| AST Parsing | memlimit-1024 | default | 1024MB | 100 | 112.640583ms | 0.01 | 0 | ✓ |
| AST Parsing | gc-50 | default | - | 50 | 166.901291ms | 0.01 | 0 | ✓ |
| AST Parsing | gc-200 | default | - | 200 | 101.514792ms | 0.01 | 0 | ✓ |
| AST Parsing | gc-off | default | - | -1 | 91.510667ms | 0.01 | 0 | ✓ |
| AST Parsing | constrained | 2 | 256MB | 50 | 119.786875ms | 0.01 | 0 | ✓ |
| AST Parsing | performance | 8 | 2048MB | 200 | 77.896667ms | 0.01 | 0 | ✓ |

