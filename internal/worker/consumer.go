package worker

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"time"

	influxdb2 "github.com/influxdata/influxdb-client-go/v2"
	"github.com/nats-io/nats.go"
	"github.com/yourorg/nms-go/internal/common/adapter"
	"github.com/yourorg/nms-go/internal/common/config"
	commonModel "github.com/yourorg/nms-go/internal/common/model"
)

type Worker struct {
	natsConn     *nats.Conn
	influxClient influxdb2.Client
	influxConfig config.InfluxConfig
	stopChan     chan struct{}
}

func NewWorker(nc *nats.Conn, ic influxdb2.Client, iConfig config.InfluxConfig) *Worker {
	return &Worker{
		natsConn:     nc,
		influxClient: ic,
		influxConfig: iConfig,
		stopChan:     make(chan struct{}),
	}
}

func (w *Worker) Start() {
	log.Println("Worker started, subscribing to nms.poll.tasks")

	sub, err := w.natsConn.Subscribe("nms.poll.tasks", func(msg *nats.Msg) {
		var task commonModel.PollTask
		if err := json.Unmarshal(msg.Data, &task); err != nil {
			log.Printf("Error unmarshalling task: %v", err)
			return
		}

		fmt.Printf("Initial worker received task: %v\n", task)
		go w.processTask(task)
	})

	if err != nil {
		log.Fatalf("Error communicating with NATS: %v", err)
	}
	defer sub.Unsubscribe()

	<-w.stopChan
}

func (w *Worker) Stop() {
	close(w.stopChan)
}

func (w *Worker) processTask(task commonModel.PollTask) {
	// Adapter selection logic
	var rtt time.Duration
	var success bool
	var metrics map[string]interface{}

	// Measure total poll duration
	pollStart := time.Now()

	if task.Protocol == "mikrotik_api" {
		// TODO: Fetch credentials from somewhere secure.
		// For MVP, hardcoded or passed in task (security risk)
		// Assuming "admin" / "admin" for test
		mtAdapter := adapter.NewMikrotikAdapter()
		m, ok := mtAdapter.FetchSystemResources(task.IPAddress, "admin", "admin")
		success = ok
		metrics = m

		// Also do a ping for RTT
		pingAdapter := &PingAdapter{}
		rtt, _ = pingAdapter.Ping(task.IPAddress)

	} else {
		// Default to Ping
		pingAdapter := &PingAdapter{}
		rtt, success = pingAdapter.Ping(task.IPAddress)
	}

	duration := time.Since(pollStart)

	// Write metrics to Influx
	writeAPI := w.influxClient.WriteAPIBlocking(w.influxConfig.Org, w.influxConfig.Bucket)

	rttMs := float64(rtt.Microseconds()) / 1000.0
	p := influxdb2.NewPoint(
		"device_poll",
		map[string]string{
			"device_id":   task.DeviceID,
			"ip_address":  task.IPAddress,
			"device_type": task.DeviceType,
		},
		map[string]interface{}{
			"rtt_ms":           rttMs,
			"success":          success,
			"poll_duration_ms": float64(duration.Milliseconds()),
		},
		time.Now(),
	)

	if err := writeAPI.WritePoint(context.Background(), p); err != nil {
		log.Printf("Error writing metrics to Influx: %v", err)
	}

	// Prepare Values map
	values := map[string]interface{}{
		"rtt_ms":  rttMs,
		"success": success,
	}

	// Add other collected metrics (e.g. from Mikrotik)
	for k, v := range metrics {
		values[k] = v
	}

	// Publish metric to Alert Engine
	metric := commonModel.Metric{
		DeviceID:  task.DeviceID,
		IPAddress: task.IPAddress,
		Timestamp: time.Now(),
		Values:    values,
	}

	payload, _ := json.Marshal(metric)
	if err := w.natsConn.Publish("nms.metrics", payload); err != nil {
		log.Printf("Error publishing metrics to NATS: %v", err)
	}
}
