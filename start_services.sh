#!/bin/bash
echo "Starting go-nms services..."

# Start Alert Service
go run cmd/alert/main.go > alert.log 2>&1 &
echo "Alert Service started (PID: $!)"

# Start Worker Service
go run cmd/worker/main.go > worker.log 2>&1 &
echo "Worker Service started (PID: $!)"

# Start Collector Service
go run cmd/collector/main.go > collector.log 2>&1 &
echo "Collector Service started (PID: $!)"

echo "All services are running in the background."
