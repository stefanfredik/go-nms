package main

import (
	"log"

	"github.com/yourorg/nms-go/internal/common/config"
	"github.com/yourorg/nms-go/internal/common/database"
	"github.com/yourorg/nms-go/internal/device/model"
)

func main() {
	// Load configuration
	cfg, err := config.LoadConfig()
	if err != nil {
		log.Fatalf("Failed to load config: %v", err)
	}

	// Connect to database
	db, err := database.NewPostgresConnection(cfg.Database)
	if err != nil {
		log.Fatalf("Failed to connect to database: %v", err)
	}

	// Run migrations
	log.Println("Starting migration...")
	err = db.AutoMigrate(
		&model.Device{},
		&model.DeviceCredentials{},
		&model.DeviceGroup{},
	)
	if err != nil {
		log.Fatalf("Migration failed: %v", err)
	}

	log.Println("Migration completed successfully!")
}
