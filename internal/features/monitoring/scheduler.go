package monitoring

import (
	"context"
	"log"
	"sync"
	"time"
)

type Scheduler struct {
	store  *TargetStore
	writer MetricWriter
	ticker *time.Ticker
	quit   chan struct{}
	wg     sync.WaitGroup
}

func NewScheduler(store *TargetStore, writer MetricWriter) *Scheduler {
	return &Scheduler{
		store:  store,
		writer: writer,
		quit:   make(chan struct{}),
	}
}

func (s *Scheduler) Start(interval time.Duration) {
	s.ticker = time.NewTicker(interval)
	go func() {
		for {
			select {
			case <-s.ticker.C:
				s.runCollection()
			case <-s.quit:
				s.ticker.Stop()
				return
			}
		}
	}()
	log.Printf("Monitoring Scheduler started with interval %v", interval)
}

func (s *Scheduler) Stop() {
	close(s.quit)
	s.wg.Wait()
	s.writer.Close()
	log.Println("Monitoring Scheduler stopped")
}

func (s *Scheduler) runCollection() {
	targets := s.store.GetAll()
	log.Printf("Starting collection for %d devices", len(targets))

	for _, target := range targets {
		s.wg.Add(1)
		go func(t DeviceTarget) {
			defer s.wg.Done()

			// Context with timeout for every poll
			ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
			defer cancel()

			if err := PollDevice(ctx, t, s.writer); err != nil {
				log.Printf("Failed to poll %s: %v", t.IP, err)
			}
		}(target)
	}
}
