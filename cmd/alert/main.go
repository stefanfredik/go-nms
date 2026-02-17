package main

import (
	"log"
	"os"
	"os/signal"
	"syscall"

	"github.com/yourorg/nms-go/internal/alert"
	"github.com/yourorg/nms-go/internal/common/config"
	"github.com/yourorg/nms-go/internal/common/queue"
	"github.com/yourorg/nms-go/internal/notification"
)

func main() {
	log.Println("Starting Alert Service...")

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

	// Initialize Services
	notifier := notification.NewEmailService()
	engine := alert.NewEngine(nc, notifier)
	go engine.Start()

	// Wait for shutdown signal
	c := make(chan os.Signal, 1)
	signal.Notify(c, os.Interrupt, syscall.SIGTERM)
	<-c

	log.Println("Stopping Alert Service...")
	engine.Stop()
}
