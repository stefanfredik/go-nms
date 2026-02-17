package main

import (
	"context"
	"log"
	"time"

	"github.com/yourorg/nms-go/internal/common/config"
	"github.com/yourorg/nms-go/internal/common/database"
	"github.com/yourorg/nms-go/internal/common/queue"
)

func main() {
	// 1. Load Config
	// Manually set config for verification if env vars aren't set
	cfg := config.Config{
		Database: config.DatabaseConfig{
			Host:     "localhost",
			Port:     5432,
			User:     "nms",
			Password: "nms_password",
			DBName:   "nms_db",
			SSLMode:  "disable",
		},
		Redis: config.RedisConfig{
			Addr:     "localhost:6379",
			Password: "",
			DB:       0,
		},
		NATS: config.NATSConfig{
			URL: "nats://localhost:4222",
		},
		Influx: config.InfluxConfig{
			URL:    "http://localhost:8086",
			Token:  "my-token", // NOTE: Needs actual token after Influx setup
			Org:    "nms_org",
			Bucket: "nms_bucket",
		},
	}

	log.Println("Verifying infrastructure connections...")

	// 2. Verify Postgres
	db, err := database.NewPostgresConnection(cfg.Database)
	if err != nil {
		log.Printf("❌ Postgres: Failed (%v)", err)
	} else {
		sqlDB, _ := db.DB()
		if err := sqlDB.Ping(); err != nil {
			log.Printf("❌ Postgres: Ping failed (%v)", err)
		} else {
			log.Println("✅ Postgres: Connected")
		}
	}

	// 3. Verify Redis
	rdb, err := database.NewRedisConnection(cfg.Redis)
	if err != nil {
		log.Printf("❌ Redis: Failed (%v)", err)
	} else {
		log.Println("✅ Redis: Connected")
		rdb.Close()
	}

	// 4. Verify NATS
	nc, err := queue.NewNATSConnection(cfg.NATS)
	if err != nil {
		log.Printf("❌ NATS: Failed (%v)", err)
	} else {
		log.Println("✅ NATS: Connected")
		nc.Close()
	}

	// 5. Verify InfluxDB
	// Note: Influx might fail auth if token is wrong, but we check connectivity
	influxClient, err := database.NewInfluxConnection(cfg.Influx)
	if err != nil {
		log.Printf("❌ InfluxDB: Failed (%v)", err)
	} else {
		log.Println("✅ InfluxDB: Connected")
		influxClient.Close()
	}
}
