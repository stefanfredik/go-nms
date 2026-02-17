package database

import (
	"context"
	"fmt"
	"time"

	influxdb2 "github.com/influxdata/influxdb-client-go/v2"
	"github.com/yourorg/nms-go/internal/common/config"
)

func NewInfluxConnection(cfg config.InfluxConfig) (influxdb2.Client, error) {
	client := influxdb2.NewClient(cfg.URL, cfg.Token)

	// Verify connection
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	ok, err := client.Ping(ctx)
	if err != nil {
		return nil, fmt.Errorf("failed to connect to influxdb: %w", err)
	}

	if !ok {
		return nil, fmt.Errorf("influxdb health check failed")
	}

	return client, nil
}
