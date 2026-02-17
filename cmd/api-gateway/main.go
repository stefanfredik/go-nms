package main

import (
	"fmt"
	"log"
	"time"

	apigateway "github.com/yourorg/nms-go/internal/api-gateway"
	"github.com/yourorg/nms-go/internal/common/config"
	"github.com/yourorg/nms-go/internal/common/database"
	"github.com/yourorg/nms-go/internal/device/model"
	"github.com/yourorg/nms-go/internal/features/monitoring"
	// "github.com/yourorg/nms-go/internal/common/database"
)

func main() {
	cfg, err := config.LoadConfig()
	if err != nil {
		log.Fatalf("Failed to load config: %v", err)
	}

	// db, err := database.NewPostgresConnection(cfg.Database) ...

	// Connect to Database
	db, err := database.NewPostgresConnection(cfg.Database)
	if err != nil {
		log.Fatalf("Failed to connect to database: %v", err)
	}

	// Auto Migrate
	if err := database.Migrate(db, &model.Device{}, &model.DeviceCredentials{}, &model.DeviceGroup{}); err != nil {
		log.Printf("Failed to run migrations: %v", err)
	}

	// Initialize Monitoring Components
	targetStore := monitoring.NewTargetStore()

	influxWriter := monitoring.NewInfluxDBWriter(
		cfg.Influx.URL,
		cfg.Influx.Token,
		cfg.Influx.Org,
		cfg.Influx.Bucket,
	)

	scheduler := monitoring.NewScheduler(targetStore, influxWriter)
	scheduler.Start(60 * time.Second) // Poll every 60s
	defer scheduler.Stop()

	monitoringHandler := monitoring.NewHandler(targetStore)

	r := apigateway.NewRouter(cfg, db, monitoringHandler)

	addr := fmt.Sprintf(":%d", cfg.Server.Port)
	log.Printf("Starting API Gateway on %s", addr)
	if err := r.Run(addr); err != nil {
		log.Fatalf("Failed to start server: %v", err)
	}
}
