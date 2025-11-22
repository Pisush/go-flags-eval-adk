# Go Flags Evaluation for Agent Development Kit

A comprehensive benchmarking repository to evaluate the impact of Go runtime flags (GOMAXPROCS, GOMEMLIMIT, GOGC) on applications built with the Agent Development Kit for Go.

## Overview

This project benchmarks real agentic coding tasks with different Go runtime configurations to help you optimize your ADK-based applications for:
- **Performance**: Execution speed and throughput
- **Memory efficiency**: RAM usage and GC behavior
- **Resource constraints**: Container and cloud deployments

## Agent Tasks

The repository includes realistic agentic coding scenarios:

### Common Use Cases
1. **Code Generator** - Generates Go source files with functions and types
2. **File Searcher** - Searches codebases for patterns using concurrent workers
3. **Refactorer** - Performs code transformations across multiple files

### Edge Cases
4. **AST Parser** - Memory-intensive parsing of Go abstract syntax trees

Each task is designed to stress different aspects of the Go runtime.

## Quick Start

```bash
# Generate test data
./scripts/generate_testdata.sh

# Run all benchmarks
go run ./cmd/benchmark

# Generate report
go run ./cmd/report -input=benchmark_results.json -output=BENCHMARK_REPORT.md
```

## Go Runtime Flags

### GOMAXPROCS

Controls maximum OS threads for Go code execution.

```bash
GOMAXPROCS=4 go run ./cmd/agents/code_generator
```

**Impact:**
- CPU-bound tasks: Higher values improve parallelism
- I/O-bound tasks: Lower values reduce overhead
- Default: `runtime.NumCPU()`

### GOMEMLIMIT

Sets soft memory limit for the runtime (Go 1.19+).

```bash
GOMEMLIMIT=512MiB go run ./cmd/agents/ast_parser
```

**Impact:**
- Triggers aggressive GC near limit
- Essential for containers with memory limits
- Prevents OOM kills in constrained environments

### GOGC

Controls garbage collector aggressiveness.

```bash
GOGC=200 go run ./cmd/agents/file_searcher
```

**Impact:**
- Lower values (50): More frequent GC, less memory
- Higher values (200): Less frequent GC, more memory
- -1: Disables automatic GC
- Default: 100

## Project Structure

```
.
├── cmd/
│   ├── agents/              # Agent programs for benchmarking
│   │   ├── code_generator/  # Generates Go code files
│   │   ├── file_searcher/   # Searches files concurrently
│   │   ├── refactor/        # Refactors code across files
│   │   └── ast_parser/      # Parses Go AST (memory-intensive)
│   ├── benchmark/           # Benchmark runner
│   └── report/              # Report generator
├── tools/                   # Reusable ADK tools
├── testdata/                # Test files for benchmarking
├── scripts/                 # Helper scripts
└── examples/                # Example configurations

```

## Running Benchmarks

### All Tasks with All Configurations

```bash
go run ./cmd/benchmark
```

This tests 13 configurations across 4 agent tasks (52 total runs):
- Default settings
- GOMAXPROCS variations (1, 2, 4, 8)
- GOMEMLIMIT variations (256MB, 512MB, 1GB)
- GOGC variations (50, 200, off)
- Combined scenarios (constrained, performance)

### Specific Task

```bash
# Run only code generation benchmarks
go run ./cmd/benchmark -task=code-gen

# Run only AST parser (memory-intensive)
go run ./cmd/benchmark -task=ast-parser
```

### Custom Configurations

Modify `cmd/benchmark/main.go` to add your own configurations:

```go
configs := []BenchmarkConfig{
    {"my-config", 4, 1024, 150},
    // MaxProcs, MemLimit(MB), GCPercent
}
```

## Analyzing Results

### Generate Report

```bash
go run ./cmd/report -input=results/benchmark_results.json -output=BENCHMARK_REPORT.md
```

The report includes:
- **Executive Summary**: Overall statistics
- **Task Analysis**: Best configurations per task
- **Recommendations**: Flag tuning guidance based on results
- **Complete Data**: Full results table

### View Sample Results

See **[SAMPLE_REPORT.md](SAMPLE_REPORT.md)** for example benchmark results from an Apple M2 with 24GB RAM. Your results will vary based on your hardware.

**Key findings from sample run:**
- **Fastest**: GOMAXPROCS=4 (56ms)
- **52 configurations tested** across 4 agent tasks
- **All runs successful** with minimal memory usage
- **Optimal sweet spot**: 4 cores for these workloads

## Individual Agent Usage

Run agents independently to test specific scenarios:

### Code Generator

```bash
go run ./cmd/agents/code_generator \
  -files=20 \
  -lines=200 \
  -output=./generated

# With custom flags
GOMAXPROCS=4 GOGC=200 go run ./cmd/agents/code_generator -files=50
```

### File Searcher

```bash
go run ./cmd/agents/file_searcher \
  -pattern="func" \
  -dir=./testdata \
  -workers=8

# Memory constrained
GOMEMLIMIT=256MiB GOGC=50 go run ./cmd/agents/file_searcher
```

### Refactorer

```bash
go run ./cmd/agents/refactor \
  -target=./testdata \
  -operation=rename \
  -old=oldVar \
  -new=newVar

# Or add comments
go run ./cmd/agents/refactor -operation=add-comments
```

### AST Parser (Memory-Intensive)

```bash
go run ./cmd/agents/ast_parser \
  -target=./testdata \
  -imports=true \
  -funcs=true \
  -types=true

# Test memory limits
GOMEMLIMIT=128MiB go run ./cmd/agents/ast_parser
```

## Integration with ADK

This repository uses Google's Agent Development Kit for Go:

```go
import "google.golang.org/adk"
```

The `tools/` package demonstrates creating custom ADK tools for file operations:
- `FileReadTool` - Read files
- `FileWriteTool` - Write files
- `GrepTool` - Search patterns
- `ListFilesTool` - List directories

## Recommendations by Use Case

### High-Throughput API (CPU-Bound)

```bash
GOMAXPROCS=8 GOGC=200 GOMEMLIMIT=2GiB
```

Maximize parallelism, reduce GC frequency.

### Memory-Constrained Container

```bash
GOMAXPROCS=2 GOGC=50 GOMEMLIMIT=450MiB  # for 512MB container
```

Aggressive GC, respect memory limits.

### Batch Processing

```bash
GOMAXPROCS=16 GOGC=400
```

Maximum performance, plenty of memory.

### Low-Latency Service

```bash
GOMAXPROCS=4 GOGC=75 GOMEMLIMIT=1GiB
```

Frequent GC for consistent latency.

## Best Practices

1. **Establish Baseline**: Run benchmarks with default settings first
2. **Measure Real Workloads**: Use tasks similar to your production code
3. **Test Under Load**: Simulate production traffic patterns
4. **Monitor in Production**: Validate benchmark findings with real metrics
5. **Container Awareness**: Always set GOMEMLIMIT in containerized environments

## Performance Metrics

Each benchmark reports:
- **Duration**: Total execution time
- **Memory Allocated**: Bytes allocated during execution
- **GC Runs**: Number of garbage collection cycles
- **GC Pause Time**: Total time spent in GC pauses
- **Exit Code**: Success/failure status

## Contributing

Contributions welcome! Areas for improvement:
- Additional agent tasks (API clients, database operations)
- Cloud-specific optimizations (GKE, Cloud Run)
- Integration with profiling tools (pprof)
- Statistical analysis of results

## References

- [Agent Development Kit for Go](https://github.com/google/adk-go)
- [Go Runtime Package](https://pkg.go.dev/runtime)
- [Go GC Guide](https://tip.golang.org/doc/gc-guide)
- [GOMEMLIMIT Blog Post](https://go.dev/blog/go119-memlimit)

## License

MIT

---

**Sources:**
- [Go - Agent Development Kit](https://google.github.io/adk-docs/get-started/go/)
- [GitHub - google/adk-go](https://github.com/google/adk-go)
- [Announcing the Agent Development Kit for Go](https://developers.googleblog.com/en/announcing-the-agent-development-kit-for-go-build-powerful-ai-agents-with-your-favorite-languages/)
