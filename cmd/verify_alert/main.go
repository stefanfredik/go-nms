package main

import (
	"encoding/json"
	"log"
	"time"

	"github.com/yourorg/nms-go/internal/common/config"
	commonModel "github.com/yourorg/nms-go/internal/common/model"
	"github.com/yourorg/nms-go/internal/common/queue"
)

func main() {
	// 1. Load Config (for NATS)
	cfg := config.Config{
		NATS: config.NATSConfig{
			URL: "nats://localhost:4222",
		},
	}

	// 2. Connect to NATS
	nc, err := queue.NewNATSConnection(cfg.NATS)
	if err != nil {
		log.Fatalf("Failed to connect to NATS: %v", err)
	}
	defer nc.Close()

	// 3. Publish High Latency Metric
	highLatency := commonModel.Metric{
		DeviceID:   "test-device-1",
		DeviceName: "Core Switch",
		IPAddress:  "192.168.1.1",
		Timestamp:  time.Now(),
		Values: map[string]interface{}{
			"rtt_ms":  150.0,
			"success": true,
		},
	}
	
	payload, _ := json.Marshal(highLatency)
	nc.Publish("nms.metrics", payload)
	log.Println("Sent High Latency Metric (>100ms)")

	// 4. Publish Device Down Metric
	deviceDown := commonModel.Metric{
		DeviceID:   "test-device-2",
		DeviceName: "Access Point",
		IPAddress:  "192.168.1.20",
		Timestamp:  time.Now(),
		Values: map[string]interface{}{
			"rtt_ms":  0.0,
			"success": false,
		},
	}
	
	payload, _ = json.Marshal(deviceDown)
	nc.Publish("nms.metrics", payload)
	log.Println("Sent Device Down Metric (success=false)")

	time.Sleep(2 * time.Second)
	log.Println("Verification messages sent. Check Alert Engine logs.")
}
