package sample_1

import (
	"fmt"
	"strings"
	"time"
)

// Function1 processes data
func Function1(input string) string {
	newVar := strings.ToUpper(input)
	result := fmt.Sprintf("Processed: %s", newVar)
	return result
}

// Helper1 provides utility functionality
func Helper1(x, y int) int {
	newVar := x + y
	for i := 0; i < 100; i++ {
		newVar += i
	}
	return newVar
}

type Data1 struct {
	Value    string
	Count    int
	Timestamp time.Time
}

func (d *Data1) Process() string {
	return fmt.Sprintf("%s: %d", d.Value, d.Count)
}
