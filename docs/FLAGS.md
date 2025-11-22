# Go Runtime Flags Detailed Guide

This document provides detailed information about Go runtime flags and their impact on application performance.

## Table of Contents

1. [GOMAXPROCS](#gomaxprocs)
2. [GOMEMLIMIT](#gomemlimit)
3. [GOGC](#gogc)
4. [GODEBUG](#godebug)
5. [GOTRACEBACK](#gotraceback)
6. [Best Practices](#best-practices)

## GOMAXPROCS

### What It Does

`GOMAXPROCS` sets the maximum number of operating system threads that can execute user-level Go code simultaneously. This is one of the most important tuning parameters for Go applications.

### How to Set

```bash
# Environment variable
export GOMAXPROCS=4
go run main.go

# Or inline
GOMAXPROCS=4 go run main.go

# Runtime API
runtime.GOMAXPROCS(4)
```

### Default Value

By default, GOMAXPROCS is set to the number of logical CPUs available (`runtime.NumCPU()`).

### When to Adjust

**Increase GOMAXPROCS when:**
- You have CPU-intensive workloads
- Your application can benefit from parallelism
- You have more CPUs than the default
- Profiling shows CPU underutilization

**Decrease GOMAXPROCS when:**
- Running in resource-constrained environments
- Your workload is mostly I/O-bound
- You want to reduce scheduling overhead
- Running multiple Go processes on the same machine

### Example Impact

```
GOMAXPROCS=1:  Serial execution, no parallelism
GOMAXPROCS=2:  Can utilize 2 CPUs in parallel
GOMAXPROCS=8:  Can utilize 8 CPUs in parallel
GOMAXPROCS=NumCPU(): Uses all available CPUs
```

### Performance Considerations

- **CPU-bound tasks**: Higher values generally better
- **I/O-bound tasks**: Lower values may reduce overhead
- **Goroutine scheduling**: More threads = more scheduling work
- **Cache coherency**: Too many threads can increase cache misses

## GOMEMLIMIT

### What It Does

`GOMEMLIMIT` (Go 1.19+) sets a soft memory limit for the Go runtime. When memory usage approaches this limit, the garbage collector becomes more aggressive.

### How to Set

```bash
# Environment variable (supports units: B, KiB, MiB, GiB, TiB)
export GOMEMLIMIT=512MiB
go run main.go

# Runtime API (in bytes)
debug.SetMemoryLimit(512 * 1024 * 1024)
```

### Default Value

By default, there is no memory limit (effectively unlimited).

### When to Use

**Set GOMEMLIMIT when:**
- Running in containers with memory limits
- Running in Kubernetes with memory requests/limits
- You want to prevent OOM kills
- Operating in memory-constrained environments

**Recommended Values:**
- **Containers**: Set to 80-90% of container memory limit
- **Kubernetes**: Set to match memory request/limit
- **Shared hosts**: Set based on fair share of available memory

### Example Configurations

```bash
# Container with 512MB limit - set Go limit to 450MB
GOMEMLIMIT=450MiB

# Kubernetes pod with 2GB limit - set to 1800MB
GOMEMLIMIT=1800MiB

# Development machine - no limit needed
# (don't set GOMEMLIMIT)
```

### Important Notes

- This is a **soft limit**, not a hard cap
- Runtime may exceed limit temporarily
- Helps prevent OOM but doesn't guarantee it
- Works in conjunction with GOGC

## GOGC

### What It Does

`GOGC` controls the garbage collector's target heap growth before triggering a GC cycle.

### How to Set

```bash
# Environment variable (percentage)
export GOGC=200
go run main.go

# Runtime API
debug.SetGCPercent(200)
```

### Default Value

Default is `100`, meaning GC runs when heap doubles in size.

### Understanding the Values

- **GOGC=100** (default): GC when heap size = 2x live data
- **GOGC=200**: GC when heap size = 3x live data (less frequent GC)
- **GOGC=50**: GC when heap size = 1.5x live data (more frequent GC)
- **GOGC=off** or **GOGC=-1**: Disable automatic GC

### When to Adjust

**Increase GOGC (e.g., 200-400) when:**
- You have plenty of memory
- GC overhead is too high
- You want to maximize throughput
- Memory allocation spikes are acceptable

**Decrease GOGC (e.g., 20-75) when:**
- Memory is constrained
- You need predictable memory usage
- Working with GOMEMLIMIT
- Reducing GC pause time is important

**Disable GOGC (=-1) when:**
- Extremely memory-rich environment
- You'll manually trigger GC
- Building specialized applications

### Trade-offs

| GOGC Value | GC Frequency | Memory Usage | CPU Usage | Latency |
|------------|--------------|--------------|-----------|---------|
| 50         | High         | Low          | High      | More consistent |
| 100        | Medium       | Medium       | Medium    | Balanced |
| 200        | Low          | High         | Low       | Can spike |
| -1         | Manual only  | Very High    | Lowest    | Unpredictable |

## GODEBUG

### What It Does

`GODEBUG` enables various runtime debugging features.

### Common Options

```bash
# GC trace - prints GC information
GODEBUG=gctrace=1 go run main.go

# Memory allocation trace
GODEBUG=allocfreetrace=1 go run main.go

# Scheduler trace
GODEBUG=schedtrace=1000 go run main.go

# Multiple options
GODEBUG=gctrace=1,schedtrace=1000 go run main.go
```

### Useful GODEBUG Settings for Evaluation

#### gctrace=1

Prints GC statistics to stderr:

```
gc 1 @0.001s 2%: 0.009+0.23+0.003 ms clock, 0.037+0.10/0.22/0.068+0.013 ms cpu, 4->4->0 MB, 5 MB goal, 4 P
```

Breakdown:
- `gc 1`: GC cycle number
- `@0.001s`: Time since program start
- `2%`: Percentage of time in GC
- Memory sizes: before->after->live data

#### schedtrace=X

Prints scheduler information every X milliseconds:

```
SCHED 0ms: gomaxprocs=4 idleprocs=0 threads=5 spinningthreads=0 idlethreads=0 runqueue=0 [0 0 0 0]
```

#### allocfreetrace=1

Prints every memory allocation and free (very verbose).

## GOTRACEBACK

### What It Does

Controls the amount of detail in stack traces when a panic occurs.

### Values

```bash
GOTRACEBACK=none    # Omit goroutine stack traces
GOTRACEBACK=single  # Single goroutine (default)
GOTRACEBACK=all     # All goroutines
GOTRACEBACK=system  # All goroutines + runtime frames
GOTRACEBACK=crash   # All goroutines + runtime frames + OS core dump
```

### When to Use

- **Development**: `all` or `system` for maximum debugging info
- **Production**: `single` (default) to reduce log size
- **Critical debugging**: `crash` to get core dumps

## Best Practices

### 1. Start with Defaults

Always start with Go's defaults and only tune when you have evidence of a problem.

### 2. Measure Before Tuning

Use this evaluation tool to establish baselines:

```bash
# Baseline with defaults
go run cmd/eval/main.go -duration=1m -workload=mixed -verbose

# Test with changes
go run cmd/eval/main.go -duration=1m -workload=mixed -verbose -maxprocs=4 -gcpercent=200
```

### 3. Container-Specific Settings

For containerized applications:

```bash
# Set GOMEMLIMIT to 90% of container limit
# Set GOMAXPROCS to match CPU quota
GOMAXPROCS=2 GOMEMLIMIT=450MiB go run main.go
```

### 4. Production Recommendations

```bash
# CPU-optimized service
GOMAXPROCS=<num_cores> GOGC=200

# Memory-constrained service
GOMEMLIMIT=<80%_of_limit>MiB GOGC=50

# Balanced service (start here)
GOMAXPROCS=<num_cores> GOMEMLIMIT=<80%_of_limit>MiB GOGC=100
```

### 5. Monitoring in Production

Always monitor these metrics:
- GC pause times
- GC frequency
- Memory usage
- CPU utilization
- Request latency (P50, P95, P99)

### 6. Avoid Premature Optimization

Only tune when you have:
- Performance requirements not being met
- Measurements showing specific bottlenecks
- Understanding of your workload characteristics

## Common Scenarios

### Scenario 1: High-Throughput API

```bash
# Maximize parallelism, reduce GC frequency
GOMAXPROCS=8 GOGC=200 GOMEMLIMIT=2GiB
```

### Scenario 2: Memory-Limited Container

```bash
# Aggressive GC, respect memory limit
GOMAXPROCS=2 GOGC=50 GOMEMLIMIT=450MiB
```

### Scenario 3: Batch Processing

```bash
# Max performance, plenty of memory
GOMAXPROCS=16 GOGC=400
```

### Scenario 4: Low-Latency Service

```bash
# Frequent GC for consistent latency
GOMAXPROCS=4 GOGC=75 GOMEMLIMIT=1GiB
```

## References

- [Go Runtime Package](https://pkg.go.dev/runtime)
- [Go GC Guide](https://tip.golang.org/doc/gc-guide)
- [GOMEMLIMIT Blog Post](https://go.dev/blog/go119-memlimit)
- [Runtime Environment Variables](https://pkg.go.dev/runtime#hdr-Environment_Variables)
