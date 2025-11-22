.PHONY: build test clean help testdata benchmark report all-agents run-all

# Build all tools
build:
	@echo "Building all tools..."
	@go build -o bin/benchmark ./cmd/benchmark
	@go build -o bin/report ./cmd/report
	@go build -o bin/code-generator ./cmd/agents/code_generator
	@go build -o bin/file-searcher ./cmd/agents/file_searcher
	@go build -o bin/refactor ./cmd/agents/refactor
	@go build -o bin/ast-parser ./cmd/agents/ast_parser
	@echo "Build complete! Binaries in ./bin/"

# Run tests
test:
	@echo "Running tests..."
	@go test -v ./...

# Generate test data
testdata:
	@echo "Generating test data..."
	@./scripts/generate_testdata.sh

# Run complete benchmark suite
benchmark:
	@echo "Running benchmark suite..."
	@mkdir -p results
	@go run ./cmd/benchmark -output=results/benchmark_results.json

# Run specific task benchmark
benchmark-task:
	@echo "Running benchmark for task: $(TASK)"
	@mkdir -p results
	@go run ./cmd/benchmark -task=$(TASK) -output=results/benchmark_$(TASK).json

# Generate report from benchmark results
report:
	@echo "Generating report..."
	@go run ./cmd/report -input=results/benchmark_results.json -output=BENCHMARK_REPORT.md
	@echo "Report generated: BENCHMARK_REPORT.md"

# Run all: testdata, benchmark, report
run-all: testdata
	@echo "Running complete benchmark suite..."
	@make benchmark
	@make report
	@echo ""
	@echo "=== Benchmark Complete ==="
	@echo "Results: results/benchmark_results.json"
	@echo "Report: BENCHMARK_REPORT.md"

# Run individual agents
run-code-gen:
	@go run ./cmd/agents/code_generator -files=10 -lines=100

run-file-search:
	@go run ./cmd/agents/file_searcher -pattern=func -dir=./testdata

run-refactor:
	@go run ./cmd/agents/refactor -target=./testdata -operation=rename

run-ast-parser:
	@go run ./cmd/agents/ast_parser -target=./testdata

# Quick test of all agents
all-agents:
	@echo "Testing all agents..."
	@echo "\n=== Code Generator ==="
	@go run ./cmd/agents/code_generator -files=5 -lines=50 -output=./generated
	@echo "\n=== File Searcher ==="
	@go run ./cmd/agents/file_searcher -pattern=func -dir=./testdata
	@echo "\n=== Refactorer ==="
	@go run ./cmd/agents/refactor -target=./testdata -operation=format
	@echo "\n=== AST Parser ==="
	@go run ./cmd/agents/ast_parser -target=./testdata

# Clean build artifacts
clean:
	@echo "Cleaning build artifacts..."
	@rm -rf bin/
	@rm -rf results/
	@rm -rf generated/
	@rm -f BENCHMARK_REPORT.md
	@rm -f benchmark_results.json
	@go clean

# Show help
help:
	@echo "Go Flags Evaluation for ADK - Makefile Commands"
	@echo ""
	@echo "Usage: make [target]"
	@echo ""
	@echo "Main Targets:"
	@echo "  run-all          - Generate testdata, run benchmarks, generate report"
	@echo "  benchmark        - Run complete benchmark suite"
	@echo "  report           - Generate markdown report from results"
	@echo "  testdata         - Generate test files for benchmarking"
	@echo ""
	@echo "Build Targets:"
	@echo "  build            - Build all binaries"
	@echo "  test             - Run all tests"
	@echo ""
	@echo "Individual Agents:"
	@echo "  run-code-gen     - Run code generator"
	@echo "  run-file-search  - Run file searcher"
	@echo "  run-refactor     - Run refactorer"
	@echo "  run-ast-parser   - Run AST parser"
	@echo "  all-agents       - Quick test of all agents"
	@echo ""
	@echo "Task-Specific Benchmarks:"
	@echo "  benchmark-task TASK=code-gen      - Benchmark code generation"
	@echo "  benchmark-task TASK=file-search   - Benchmark file searching"
	@echo "  benchmark-task TASK=refactor      - Benchmark refactoring"
	@echo "  benchmark-task TASK=ast-parser    - Benchmark AST parsing"
	@echo ""
	@echo "Utility:"
	@echo "  clean            - Remove build artifacts and results"
	@echo "  help             - Show this help message"
