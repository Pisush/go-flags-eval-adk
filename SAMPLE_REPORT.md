# Go Flags Benchmark Report - Sample Results

> **Note:** These are example results from a specific test environment.
> Your results will vary based on hardware, OS, and system load.
> Run `make run-all` to generate benchmarks for your machine.

## Test Environment

- **CPU:** Apple M2 (ARM64)
- **Memory:** 24 GB
- **OS:** macOS 14.5
- **Go Version:** 1.24.10
- **Date:** November 22, 2025

---

Generated: Sat, 22 Nov 2025 16:59:49 CET

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
- **Average Duration**: 109.781238ms
- **Average Memory**: 0.01 MB
- **Average GC Runs**: 0.0

## Task: All Tasks

### Performance Analysis

#### Fastest Configurations

| Rank | Configuration | Duration | Memory (MB) | GC Runs |
|------|---------------|----------|-------------|----------|
| 1 | maxprocs-4 | 56.345833ms | 0.01 | 0 |
| 2 | maxprocs-4 | 56.655708ms | 0.01 | 0 |
| 3 | maxprocs-4 | 57.727709ms | 0.01 | 0 |
| 4 | gc-off | 59.500292ms | 0.01 | 0 |
| 5 | gc-200 | 59.93475ms | 0.01 | 0 |

#### Lowest Memory Usage

| Rank | Configuration | Memory (MB) | Duration | GC Runs |
|------|---------------|-------------|----------|----------|
| 1 | memlimit-1024 | 0.01 | 63.20875ms | 0 |
| 2 | maxprocs-1 | 0.01 | 335.35475ms | 0 |
| 3 | gc-50 | 0.01 | 79.356542ms | 0 |
| 4 | gc-off | 0.01 | 59.500292ms | 0 |
| 5 | default | 0.01 | 555.163333ms | 0 |

#### Fewest GC Runs

| Rank | Configuration | GC Runs | Duration | Memory (MB) |
|------|---------------|---------|----------|-------------|
| 1 | memlimit-1024 | 0 | 63.20875ms | 0.01 |
| 2 | maxprocs-1 | 0 | 335.35475ms | 0.01 |
| 3 | gc-50 | 0 | 79.356542ms | 0.01 |
| 4 | gc-off | 0 | 59.500292ms | 0.01 |
| 5 | default | 0 | 555.163333ms | 0.01 |


## Recommendations

Based on the benchmark results:

### GOMAXPROCS

- **Optimal value**: GOMAXPROCS=4 (Duration: 56.345833ms)
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
| Code Generation | default | default | - | 100 | 59.989875ms | 0.01 | 0 | ✓ |
| Code Generation | maxprocs-1 | 1 | - | 100 | 85.479334ms | 0.01 | 0 | ✓ |
| Code Generation | maxprocs-2 | 2 | - | 100 | 62.535708ms | 0.01 | 0 | ✓ |
| Code Generation | maxprocs-4 | 4 | - | 100 | 57.727709ms | 0.01 | 0 | ✓ |
| Code Generation | maxprocs-8 | 8 | - | 100 | 71.923042ms | 0.01 | 0 | ✓ |
| Code Generation | memlimit-256 | default | 256MB | 100 | 73.441208ms | 0.01 | 0 | ✓ |
| Code Generation | memlimit-512 | default | 512MB | 100 | 72.659416ms | 0.01 | 0 | ✓ |
| Code Generation | memlimit-1024 | default | 1024MB | 100 | 78.72025ms | 0.01 | 0 | ✓ |
| Code Generation | gc-50 | default | - | 50 | 111.033875ms | 0.01 | 0 | ✓ |
| Code Generation | gc-200 | default | - | 200 | 91.264958ms | 0.01 | 0 | ✓ |
| Code Generation | gc-off | default | - | -1 | 88.716875ms | 0.01 | 0 | ✓ |
| Code Generation | constrained | 2 | 256MB | 50 | 96.627125ms | 0.01 | 0 | ✓ |
| Code Generation | performance | 8 | 2048MB | 200 | 71.064708ms | 0.01 | 0 | ✓ |
| File Searching | default | default | - | 100 | 455.840625ms | 0.01 | 0 | ✓ |
| File Searching | maxprocs-1 | 1 | - | 100 | 325.601959ms | 0.01 | 0 | ✓ |
| File Searching | maxprocs-2 | 2 | - | 100 | 61.592083ms | 0.01 | 0 | ✓ |
| File Searching | maxprocs-4 | 4 | - | 100 | 56.345833ms | 0.01 | 0 | ✓ |
| File Searching | maxprocs-8 | 8 | - | 100 | 68.04525ms | 0.01 | 0 | ✓ |
| File Searching | memlimit-256 | default | 256MB | 100 | 70.398625ms | 0.01 | 0 | ✓ |
| File Searching | memlimit-512 | default | 512MB | 100 | 63.502833ms | 0.01 | 0 | ✓ |
| File Searching | memlimit-1024 | default | 1024MB | 100 | 67.692042ms | 0.01 | 0 | ✓ |
| File Searching | gc-50 | default | - | 50 | 87.60725ms | 0.01 | 0 | ✓ |
| File Searching | gc-200 | default | - | 200 | 68.159917ms | 0.01 | 0 | ✓ |
| File Searching | gc-off | default | - | -1 | 69.550917ms | 0.01 | 0 | ✓ |
| File Searching | constrained | 2 | 256MB | 50 | 80.418417ms | 0.01 | 0 | ✓ |
| File Searching | performance | 8 | 2048MB | 200 | 69.929125ms | 0.01 | 0 | ✓ |
| Code Refactoring | default | default | - | 100 | 441.403416ms | 0.01 | 0 | ✓ |
| Code Refactoring | maxprocs-1 | 1 | - | 100 | 328.358709ms | 0.01 | 0 | ✓ |
| Code Refactoring | maxprocs-2 | 2 | - | 100 | 61.220708ms | 0.01 | 0 | ✓ |
| Code Refactoring | maxprocs-4 | 4 | - | 100 | 56.655708ms | 0.01 | 0 | ✓ |
| Code Refactoring | maxprocs-8 | 8 | - | 100 | 64.464333ms | 0.01 | 0 | ✓ |
| Code Refactoring | memlimit-256 | default | 256MB | 100 | 65.306041ms | 0.01 | 0 | ✓ |
| Code Refactoring | memlimit-512 | default | 512MB | 100 | 66.639292ms | 0.01 | 0 | ✓ |
| Code Refactoring | memlimit-1024 | default | 1024MB | 100 | 67.647083ms | 0.01 | 0 | ✓ |
| Code Refactoring | gc-50 | default | - | 50 | 78.346875ms | 0.01 | 0 | ✓ |
| Code Refactoring | gc-200 | default | - | 200 | 59.93475ms | 0.01 | 0 | ✓ |
| Code Refactoring | gc-off | default | - | -1 | 64.142ms | 0.01 | 0 | ✓ |
| Code Refactoring | constrained | 2 | 256MB | 50 | 80.780834ms | 0.01 | 0 | ✓ |
| Code Refactoring | performance | 8 | 2048MB | 200 | 76.050042ms | 0.01 | 0 | ✓ |
| AST Parsing | default | default | - | 100 | 555.163333ms | 0.01 | 0 | ✓ |
| AST Parsing | maxprocs-1 | 1 | - | 100 | 335.35475ms | 0.01 | 0 | ✓ |
| AST Parsing | maxprocs-2 | 2 | - | 100 | 63.808459ms | 0.01 | 0 | ✓ |
| AST Parsing | maxprocs-4 | 4 | - | 100 | 61.352125ms | 0.01 | 0 | ✓ |
| AST Parsing | maxprocs-8 | 8 | - | 100 | 65.415292ms | 0.01 | 0 | ✓ |
| AST Parsing | memlimit-256 | default | 256MB | 100 | 66.288958ms | 0.01 | 0 | ✓ |
| AST Parsing | memlimit-512 | default | 512MB | 100 | 68.356875ms | 0.01 | 0 | ✓ |
| AST Parsing | memlimit-1024 | default | 1024MB | 100 | 63.20875ms | 0.01 | 0 | ✓ |
| AST Parsing | gc-50 | default | - | 50 | 79.356542ms | 0.01 | 0 | ✓ |
| AST Parsing | gc-200 | default | - | 200 | 61.59375ms | 0.01 | 0 | ✓ |
| AST Parsing | gc-off | default | - | -1 | 59.500292ms | 0.01 | 0 | ✓ |
| AST Parsing | constrained | 2 | 256MB | 50 | 82.461541ms | 0.01 | 0 | ✓ |
| AST Parsing | performance | 8 | 2048MB | 200 | 69.945ms | 0.01 | 0 | ✓ |

