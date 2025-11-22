# Project Summary: Go Flags Evaluation for ADK

## What Was Built

This repository evaluates how different Go runtime flags (GOMAXPROCS, GOMEMLIMIT, GOGC) affect the performance of applications built with Google's Agent Development Kit for Go.

## Key Components

### 1. Agent Programs (4 realistic agentic coding tasks)

**Common Use Cases:**
- **Code Generator** (`cmd/agents/code_generator`) - Generates Go source files concurrently
- **File Searcher** (`cmd/agents/file_searcher`) - Searches codebases using concurrent workers
- **Refactorer** (`cmd/agents/refactor`) - Performs code transformations across multiple files

**Edge Cases:**
- **AST Parser** (`cmd/agents/ast_parser`) - Memory-intensive Go AST parsing

### 2. ADK Tools (`tools/coding_tools.go`)

Custom ADK tools demonstrating the framework:
- `FileReadTool` - Read file contents
- `FileWriteTool` - Write files to disk
- `GrepTool` - Search for patterns
- `ListFilesTool` - List directory contents

### 3. Benchmark Infrastructure

- **Benchmark Runner** (`cmd/benchmark`) - Tests 13 flag configurations across all agents (52 total runs)
- **Report Generator** (`cmd/report`) - Creates markdown reports with analysis and recommendations

### 4. Test Data

- 55 Go files across 5 directories
- Realistic code patterns for testing
- Generated via `scripts/generate_testdata.sh`

## How It Works

### Run Complete Benchmark Suite

```bash
make run-all
```

This:
1. Generates test data (55 Go files)
2. Runs each agent with 13 different flag configurations
3. Measures duration, memory usage, and GC behavior
4. Generates a comprehensive markdown report

### Test Individual Agents

```bash
# Code generation
GOMAXPROCS=4 go run ./cmd/agents/code_generator -files=20

# File searching
GOMEMLIMIT=256MiB go run ./cmd/agents/file_searcher

# Refactoring
GOGC=50 go run ./cmd/agents/refactor -operation=add-comments

# AST parsing (memory-intensive)
GOMEMLIMIT=128MiB GOGC=50 go run ./cmd/agents/ast_parser
```

## Benchmark Configurations Tested

1. **Default** - Baseline with Go defaults
2. **GOMAXPROCS variations** - 1, 2, 4, 8 cores
3. **GOMEMLIMIT variations** - 256MB, 512MB, 1GB
4. **GOGC variations** - 50 (aggressive), 200 (conservative), off
5. **Combined scenarios**:
   - **Constrained**: GOMAXPROCS=2, GOMEMLIMIT=256MB, GOGC=50
   - **Performance**: GOMAXPROCS=8, GOMEMLIMIT=2GB, GOGC=200

## Output

### Benchmark Results (JSON)

```json
{
  "Config": {"Name": "maxprocs-8", "MaxProcs": 8, "MemLimit": 0, "GCPercent": 100},
  "Duration": "1.2s",
  "MemoryAllocated": 52428800,
  "NumGC": 8,
  "PauseTimeNs": 1200000
}
```

### Generated Report (Markdown)

- Executive summary with overall statistics
- Performance analysis by task
- Fastest/lowest memory/fewest GC rankings
- Recommendations for different use cases
- Complete data table

## Integration with ADK

The project demonstrates:
- Creating custom ADK tools with `functiontool`
- Tool definitions with JSON schema
- Concurrent agent operations
- Real-world coding scenarios

## Technologies

- **Go 1.24** with toolchain support
- **Google ADK** v0.2.0
- **Runtime tuning**: GOMAXPROCS, GOMEMLIMIT, GOGC
- **Concurrency**: Goroutines, channels, sync primitives
- **AST parsing**: Go parser and ast packages

## Use Cases

### For ADK Developers

- Understand performance characteristics of ADK agents
- Optimize flag settings for your specific workload
- Test under different resource constraints

### For Go Developers

- Learn how runtime flags affect different workload types
- See real-world examples of concurrent Go programs
- Benchmark methodology for performance tuning

### For DevOps/SRE

- Container resource planning (memory limits, CPU quotas)
- Performance optimization for cloud deployments
- GC tuning for different environments

## Quick Commands

```bash
# Build everything
make build

# Run all agents quickly
make all-agents

# Run specific task benchmark
make benchmark-task TASK=code-gen

# Generate test data
make testdata

# Clean up
make clean

# Show help
make help
```

## Files Created

- 6 Go programs (4 agents + benchmark + report)
- 1 ADK tools package
- 55 test data files
- Comprehensive documentation
- Build automation (Makefile)
- Scripts for data generation

## Next Steps

1. Run `make run-all` to generate your first benchmark report
2. Experiment with custom configurations
3. Add your own agent tasks
4. Integrate findings into your ADK projects
