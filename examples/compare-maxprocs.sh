#!/bin/bash

# Compare performance with different GOMAXPROCS values

set -e

EVAL_BIN="../cmd/eval/main.go"
DURATION="30s"
WORKLOAD="cpu"
OUTPUT_DIR="./results/maxprocs-comparison"

mkdir -p "$OUTPUT_DIR"

echo "=== GOMAXPROCS Comparison ==="
echo "Workload: $WORKLOAD"
echo "Duration: $DURATION"
echo ""

# Test with different GOMAXPROCS values
for procs in 1 2 4 8 16; do
    echo "Testing GOMAXPROCS=$procs..."
    go run "$EVAL_BIN" -duration="$DURATION" -workload="$WORKLOAD" -maxprocs="$procs" \
        | tee "$OUTPUT_DIR/maxprocs-$procs.txt"
    echo ""
done

echo "=== Comparison Complete ==="
echo "Results saved to $OUTPUT_DIR"
echo ""
echo "Summary:"
grep "Duration:" "$OUTPUT_DIR"/*.txt
echo ""
grep "Number of GC runs:" "$OUTPUT_DIR"/*.txt
