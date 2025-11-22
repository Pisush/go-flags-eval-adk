# Dockerfile for running Go flags evaluation in containers

FROM golang:1.22-alpine AS builder

WORKDIR /app

# Copy go mod files
COPY go.mod ./
COPY go.sum* ./

# Download dependencies
RUN go mod download

# Copy source code
COPY . .

# Build the evaluation tool
RUN go build -o /eval ./cmd/eval

FROM alpine:latest

WORKDIR /root/

# Copy the binary from builder
COPY --from=builder /eval .

# Set default environment variables
ENV WORKLOAD=mixed
ENV DURATION=30s
ENV VERBOSE=false

# Run the evaluation
CMD sh -c "./eval -workload=${WORKLOAD} -duration=${DURATION} -verbose=${VERBOSE}"
