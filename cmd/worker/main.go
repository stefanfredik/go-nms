package main

import (
	"log"
	"os"
	"os/signal"
	"syscall"

	"github.com/yourorg/nms-go/internal/common/config"
	"github.com/yourorg/nms-go/internal/common/database"
	"github.com/yourorg/nms-go/internal/common/queue"
	"github.com/yourorg/nms-go/internal/worker"
)

func main() {
	log.Println("Starting Worker Service...")

	cfg, err := config.LoadConfig()
	if err != nil {
		log.Fatalf("Failed to load config: %v", err)
	}

	// Connect to NATS
	nc, err := queue.NewNATSConnection(cfg.NATS)
	if err != nil {
		log.Fatalf("Failed to connect to NATS: %v", err)
	}
	defer nc.Close()

	// Connect to InfluxDB
	influxClient, err := database.NewInfluxConnection(cfg.Influx)
	if err != nil {
		log.Fatalf("Failed to connect to InfluxDB: %v", err)
	}
	defer influxClient.Close()

	// Start Worker
	w := worker.NewWorker(nc, influxClient, cfg.Influx)
	go w.Start()

	// Wait for shutdown signal
	c := make(chan os.Signal, 1)
	signal.Notify(c, os.Interrupt, syscall.SIGTERM)
	<-c

	log.Println("Stopping Worker Service...")
	w.Stop()
}
