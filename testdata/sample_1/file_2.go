package sample_1

import (
	"fmt"
	"strings"
	"time"
)

// Function2 processes data
func Function2(input string) string {
	oldVar := strings.ToUpper(input)
	result := fmt.Sprintf("Processed: %s", oldVar)
	return result
}

// Helper2 provides utility functionality
func Helper2(x, y int) int {
	oldVar := x + y
	for i := 0; i < 100; i++ {
		oldVar += i
	}
	return oldVar
}

type Data2 struct {
	Value    string
	Count    int
	Timestamp time.Time
}

func (d *Data2) Process() string {
	return fmt.Sprintf("%s: %d", d.Value, d.Count)
}
