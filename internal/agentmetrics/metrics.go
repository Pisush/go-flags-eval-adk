package agentmetrics

import (
	"encoding/json"
	"os"
	"time"
)

// Metrics represents performance metrics collected by an agent
type Metrics struct {
	Duration        time.Duration `json:"duration"`
	MemoryAllocated uint64        `json:"memory_allocated"` // Total bytes allocated
	HeapAllocated   uint64        `json:"heap_allocated"`   // Current heap size
	NumGC           uint32        `json:"num_gc"`           // Number of GC cycles
	PauseTimeNs     uint64        `json:"pause_time_ns"`    // Total GC pause time
	Goroutines      int           `json:"goroutines"`       // Number of goroutines

	// Agent-specific metrics
	TasksCompleted  int            `json:"tasks_completed,omitempty"`
	FilesProcessed  int            `json:"files_processed,omitempty"`
	Custom          map[string]any `json:"custom,omitempty"`
}

// WriteToFile writes metrics to a JSON file
func (m *Metrics) WriteToFile(filename string) error {
	data, err := json.MarshalIndent(m, "", "  ")
	if err != nil {
		return err
	}

	return os.WriteFile(filename, data, 0644)
}

// ReadFromFile reads metrics from a JSON file
func ReadFromFile(filename string) (*Metrics, error) {
	data, err := os.ReadFile(filename)
	if err != nil {
		return nil, err
	}

	var metrics Metrics
	if err := json.Unmarshal(data, &metrics); err != nil {
		return nil, err
	}

	return &metrics, nil
}
