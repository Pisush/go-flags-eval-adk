#!/bin/bash

# Test behavior under various memory constraints

set -e

EVAL_BIN="../cmd/eval/main.go"
DURATION="30s"
WORKLOAD="memory"
OUTPUT_DIR="./results/memory-limits"

mkdir -p "$OUTPUT_DIR"

echo "=== Memory Limit Evaluation ==="
echo "Workload: $WORKLOAD"
echo "Duration: $DURATION"
echo ""

# Test with different memory limits
echo "Test 1: No memory limit"
go run "$EVAL_BIN" -duration="$DURATION" -workload="$WORKLOAD" \
    | tee "$OUTPUT_DIR/no-limit.txt"
echo ""

echo "Test 2: 128MB limit with aggressive GC"
go run "$EVAL_BIN" -duration="$DURATION" -workload="$WORKLOAD" -memlimit=128 -gcpercent=50 \
    | tee "$OUTPUT_DIR/128mb-gc50.txt"
echo ""

echo "Test 3: 256MB limit with default GC"
go run "$EVAL_BIN" -duration="$DURATION" -workload="$WORKLOAD" -memlimit=256 -gcpercent=100 \
    | tee "$OUTPUT_DIR/256mb-gc100.txt"
echo ""

echo "Test 4: 512MB limit with conservative GC"
go run "$EVAL_BIN" -duration="$DURATION" -workload="$WORKLOAD" -memlimit=512 -gcpercent=200 \
    | tee "$OUTPUT_DIR/512mb-gc200.txt"
echo ""

echo "Test 5: 1GB limit with default GC"
go run "$EVAL_BIN" -duration="$DURATION" -workload="$WORKLOAD" -memlimit=1024 -gcpercent=100 \
    | tee "$OUTPUT_DIR/1gb-gc100.txt"
echo ""

echo "=== Evaluation Complete ==="
echo "Results saved to $OUTPUT_DIR"
echo ""
echo "GC Statistics Summary:"
grep "Number of GC runs:" "$OUTPUT_DIR"/*.txt
echo ""
grep "Total GC pause time:" "$OUTPUT_DIR"/*.txt
