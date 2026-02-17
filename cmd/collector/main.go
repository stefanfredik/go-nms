package main

import (
	"log"
	"os"
	"os/signal"
	"syscall"

	"github.com/yourorg/nms-go/internal/collector"
	"github.com/yourorg/nms-go/internal/common/config"
	"github.com/yourorg/nms-go/internal/common/database"
	"github.com/yourorg/nms-go/internal/common/queue"
	"github.com/yourorg/nms-go/internal/device/repository"
	"github.com/yourorg/nms-go/internal/device/service"
)

func main() {
	log.Println("Starting Collector Service...")

	cfg, err := config.LoadConfig()
	if err != nil {
		log.Fatalf("Failed to load config: %v", err)
	}

	// Connect to Database
	db, err := database.NewPostgresConnection(cfg.Database)
	if err != nil {
		log.Fatalf("Failed to connect to DB: %v", err)
	}

	// Connect to NATS
	nc, err := queue.NewNATSConnection(cfg.NATS)
	if err != nil {
		log.Fatalf("Failed to connect to NATS: %v", err)
	}
	defer nc.Close()

	// Initialize Services
	deviceRepo := repository.NewDeviceRepository(db)
	deviceService := service.NewDeviceService(deviceRepo)

	// Start Scheduler
	scheduler := collector.NewScheduler(deviceService, nc)
	go scheduler.Start()

	// Wait for shutdown signal
	c := make(chan os.Signal, 1)
	signal.Notify(c, os.Interrupt, syscall.SIGTERM)
	<-c

	log.Println("Stopping Collector Service...")
	scheduler.Stop()
}
