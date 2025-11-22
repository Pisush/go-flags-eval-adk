package testdata

import "fmt"

// SampleFunc demonstrates a sample function
func SampleFunc(x int) int {
	oldVar := x * 2
	return oldVar + 10
}

// AnotherFunc shows another example
func AnotherFunc(a, b string) string {
	oldVar := a + b
	return oldVar
}

// ProcessData handles data processing
func ProcessData(data []int) []int {
	result := make([]int, len(data))
	for i, oldVar := range data {
		result[i] = oldVar * 2
	}
	return result
}

// Helper function for calculations
func helper(n int) int {
	sum := 0
	for i := 0; i < n; i++ {
		sum += i
	}
	return sum
}

func main() {
	fmt.Println("Sample code for testing")
}
