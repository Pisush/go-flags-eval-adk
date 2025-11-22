# Quick Start Guide

Get started with Go flags evaluation in 5 minutes.

## Prerequisites

- Go 1.19 or later (for GOMEMLIMIT support)
- Basic familiarity with Go runtime concepts

## Installation

```bash
# Clone or navigate to the repository
cd vpn

# Build the evaluation tool
make build
```

## Run Your First Evaluation

```bash
# Run with default settings (30 second mixed workload)
./eval

# Run a specific workload type
./eval -workload=cpu -duration=1m

# Run with custom Go flags
./eval -maxprocs=4 -gcpercent=150 -memlimit=512
```

## Understanding the Output

```
=== Go Flags Evaluation ===
GOMAXPROCS: 4 (available CPUs: 8)
GOMEMLIMIT: 512 MB
GOGC: 100%
Workload: mixed
Duration: 30s

Starting workload...

=== Results ===
Duration: 30.002s
Memory Allocated: 1024.50 MB
Number of GC runs: 15
Total GC pause time: 45.2ms
Average GC pause time: 3.01ms
Active goroutines: 9
GOMAXPROCS: 4
```

**Key Metrics:**
- **Memory Allocated**: Total memory allocated during test
- **GC runs**: Number of garbage collection cycles
- **GC pause time**: Time spent in stop-the-world GC pauses
- **Active goroutines**: Concurrent goroutines at end

## Quick Experiments

### Experiment 1: Impact of GOMAXPROCS on CPU Workload

```bash
# Compare single-threaded vs multi-threaded
./eval -workload=cpu -maxprocs=1 -duration=30s
./eval -workload=cpu -maxprocs=8 -duration=30s

# Expected: Higher GOMAXPROCS = better CPU workload performance
```

### Experiment 2: Impact of GOGC on Memory Workload

```bash
# Aggressive GC
./eval -workload=memory -gcpercent=50 -duration=30s

# Conservative GC
./eval -workload=memory -gcpercent=200 -duration=30s

# Expected: Lower GOGC = more GC runs, higher GOGC = more memory usage
```

### Experiment 3: Memory Limits

```bash
# No limit
./eval -workload=memory -duration=30s

# With 256MB limit
./eval -workload=memory -memlimit=256 -duration=30s

# Expected: Memory limit triggers more aggressive GC
```

## Running Benchmarks

```bash
# Run all benchmarks
make bench

# Run specific benchmark
go test -bench=BenchmarkCPUWorkload ./benchmarks/

# With memory statistics
go test -bench=. -benchmem ./benchmarks/
```

## Example Output from Benchmarks

```
BenchmarkCPUWorkload/GOMAXPROCS_1-8       100    12543210 ns/op
BenchmarkCPUWorkload/GOMAXPROCS_2-8       150     7234512 ns/op
BenchmarkCPUWorkload/GOMAXPROCS_4-8       200     4123456 ns/op
```

**Reading benchmark results:**
- First number (100, 150, 200): Number of iterations
- Second number: Nanoseconds per operation (lower is better)

## Automated Evaluations

Run comprehensive tests with pre-configured scripts:

```bash
# Run all test scenarios
cd examples
./run-all.sh

# Compare GOMAXPROCS values
./compare-maxprocs.sh

# Test memory constraints
./memory-constrained.sh
```

Results are saved in `examples/results/` directory.

## Docker Quick Start

```bash
# Build Docker image
make docker-build

# Run in containers with different configurations
docker-compose -f examples/docker-compose.yml up
```

## Common Use Cases

### For Agent Development Kit Applications

```bash
# Test your agent workload characteristics
./eval -workload=mixed -verbose -duration=2m

# Optimize for throughput
./eval -maxprocs=8 -gcpercent=200

# Optimize for memory constraints
./eval -memlimit=256 -gcpercent=50 -maxprocs=2
```

### Finding Optimal Settings

1. **Start with defaults**: Run without any flags
2. **Identify bottleneck**: Check if CPU or memory limited
3. **Tune accordingly**:
   - CPU-bound: Increase GOMAXPROCS
   - Memory-bound: Set GOMEMLIMIT and tune GOGC
4. **Benchmark**: Compare before/after results

## Next Steps

- Read [README.md](README.md) for comprehensive documentation
- Check [docs/FLAGS.md](docs/FLAGS.md) for detailed flag explanations
- Review [examples/](examples/) for more test scenarios
- Integrate findings into your Agent Development Kit projects

## Troubleshooting

### Problem: High GC pause times

```bash
# Try more aggressive GC
./eval -gcpercent=50 -verbose
```

### Problem: Memory usage too high

```bash
# Set memory limit and aggressive GC
./eval -memlimit=512 -gcpercent=50 -verbose
```

### Problem: Poor CPU utilization

```bash
# Increase parallelism
./eval -maxprocs=<num_cpus> -verbose
```

## Getting Help

- Check the [README.md](README.md) for detailed documentation
- Review [docs/FLAGS.md](docs/FLAGS.md) for flag details
- Run `make help` to see all available commands
