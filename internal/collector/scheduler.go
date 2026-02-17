package collector

import (
	"context"
	"encoding/json"
	"log"
	"time"

	"github.com/nats-io/nats.go"
	commonModel "github.com/yourorg/nms-go/internal/common/model"
	"github.com/yourorg/nms-go/internal/device/service"
)

type Scheduler struct {
	deviceService service.DeviceService
	natsConn      *nats.Conn
	stopChan      chan struct{}
}

func NewScheduler(ds service.DeviceService, nc *nats.Conn) *Scheduler {
	return &Scheduler{
		deviceService: ds,
		natsConn:      nc,
		stopChan:      make(chan struct{}),
	}
}

func (s *Scheduler) Start() {
	ticker := time.NewTicker(10 * time.Second) // Check every 10 seconds (simplification)
	defer ticker.Stop()

	log.Println("Collector Scheduler started")

	for {
		select {
		case <-ticker.C:
			s.schedulePolls()
		case <-s.stopChan:
			log.Println("Collector Scheduler stopped")
			return
		}
	}
}

func (s *Scheduler) Stop() {
	close(s.stopChan)
}

func (s *Scheduler) schedulePolls() {
	ctx := context.Background()
	// In a real app, we would query DB for devices "due" for polling.
	// For now, we fetch all enabled devices and dispatch tasks.
	// Optimization: Use pagination or specific DB query for 'next_poll_at'
	
	devices, _, err := s.deviceService.ListDevices(ctx, 1, 1000)
	if err != nil {
		log.Printf("Error fetching devices: %v", err)
		return
	}

	for _, d := range devices {
		if !d.Enabled {
			continue
		}

		task := commonModel.PollTask{
			DeviceID:   d.ID,
			IPAddress:  d.IPAddress,
			DeviceType: string(d.DeviceType),
			Protocol:   string(d.Protocol),
			Timestamp:  time.Now(),
		}

		payload, _ := json.Marshal(task)
		err := s.natsConn.Publish("nms.poll.tasks", payload)
		if err != nil {
			log.Printf("Error publishing task for device %s: %v", d.Name, err)
		} else {
			// log.Printf("Scheduled poll for %s (%s)", d.Name, d.IPAddress)
		}
	}
}
