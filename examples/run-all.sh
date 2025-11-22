#!/bin/bash

# Script to run evaluation with various flag combinations

set -e

EVAL_BIN="../cmd/eval/main.go"
DURATION="30s"
OUTPUT_DIR="./results"

mkdir -p "$OUTPUT_DIR"

echo "=== Running Go Flags Evaluation Suite ==="
echo "Duration per test: $DURATION"
echo "Output directory: $OUTPUT_DIR"
echo ""

# Test 1: Default settings
echo "Test 1: Default settings"
go run "$EVAL_BIN" -duration="$DURATION" -workload=mixed | tee "$OUTPUT_DIR/01-default.txt"
echo ""

# Test 2: Low GOMAXPROCS
echo "Test 2: GOMAXPROCS=1"
go run "$EVAL_BIN" -duration="$DURATION" -workload=cpu -maxprocs=1 | tee "$OUTPUT_DIR/02-maxprocs-1.txt"
echo ""

# Test 3: High GOMAXPROCS
echo "Test 3: GOMAXPROCS=8"
go run "$EVAL_BIN" -duration="$DURATION" -workload=cpu -maxprocs=8 | tee "$OUTPUT_DIR/03-maxprocs-8.txt"
echo ""

# Test 4: Memory limit 256MB
echo "Test 4: Memory limit 256MB"
go run "$EVAL_BIN" -duration="$DURATION" -workload=memory -memlimit=256 | tee "$OUTPUT_DIR/04-memlimit-256.txt"
echo ""

# Test 5: Memory limit 1GB
echo "Test 5: Memory limit 1GB"
go run "$EVAL_BIN" -duration="$DURATION" -workload=memory -memlimit=1024 | tee "$OUTPUT_DIR/05-memlimit-1024.txt"
echo ""

# Test 6: Aggressive GC (GOGC=50)
echo "Test 6: Aggressive GC (GOGC=50)"
go run "$EVAL_BIN" -duration="$DURATION" -workload=memory -gcpercent=50 | tee "$OUTPUT_DIR/06-gc-50.txt"
echo ""

# Test 7: Conservative GC (GOGC=200)
echo "Test 7: Conservative GC (GOGC=200)"
go run "$EVAL_BIN" -duration="$DURATION" -workload=memory -gcpercent=200 | tee "$OUTPUT_DIR/07-gc-200.txt"
echo ""

# Test 8: GC disabled (GOGC=-1)
echo "Test 8: GC disabled (GOGC=-1)"
go run "$EVAL_BIN" -duration="$DURATION" -workload=memory -gcpercent=-1 | tee "$OUTPUT_DIR/08-gc-off.txt"
echo ""

# Test 9: Memory constrained configuration
echo "Test 9: Memory constrained (256MB limit, GOGC=50, MAXPROCS=2)"
go run "$EVAL_BIN" -duration="$DURATION" -workload=mixed -memlimit=256 -gcpercent=50 -maxprocs=2 | tee "$OUTPUT_DIR/09-constrained.txt"
echo ""

# Test 10: Performance optimized configuration
echo "Test 10: Performance optimized (GOGC=200, MAXPROCS=8)"
go run "$EVAL_BIN" -duration="$DURATION" -workload=mixed -gcpercent=200 -maxprocs=8 | tee "$OUTPUT_DIR/10-optimized.txt"
echo ""

echo "=== Evaluation Complete ==="
echo "Results saved to $OUTPUT_DIR"
