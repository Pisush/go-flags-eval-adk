#!/bin/bash

# Generate test data for benchmarking

set -e

TESTDATA_DIR="./testdata"

echo "Generating test data..."

# Create directory structure
for i in {1..5}; do
    mkdir -p "$TESTDATA_DIR/sample_$i"
done

# Generate sample Go files
for dir in "$TESTDATA_DIR"/sample_*; do
    for j in {1..10}; do
        cat > "$dir/file_$j.go" <<EOF
package $(basename $dir)

import (
	"fmt"
	"strings"
	"time"
)

// Function$j processes data
func Function$j(input string) string {
	oldVar := strings.ToUpper(input)
	result := fmt.Sprintf("Processed: %s", oldVar)
	return result
}

// Helper$j provides utility functionality
func Helper$j(x, y int) int {
	oldVar := x + y
	for i := 0; i < 100; i++ {
		oldVar += i
	}
	return oldVar
}

type Data$j struct {
	Value    string
	Count    int
	Timestamp time.Time
}

func (d *Data$j) Process() string {
	return fmt.Sprintf("%s: %d", d.Value, d.Count)
}
EOF
    done
done

# Copy sample.go to each directory
for dir in "$TESTDATA_DIR"/sample_*; do
    cp "$TESTDATA_DIR/sample.go" "$dir/"
done

echo "Test data generated successfully!"
echo "Created 5 directories with 11 Go files each (55 files total)"
