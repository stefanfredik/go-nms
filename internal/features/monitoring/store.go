package monitoring

import (
	"sync"
)

// TargetStore manages the in-memory list of devices to poll
type TargetStore struct {
	mu      sync.RWMutex
	targets map[string]DeviceTarget
}

func NewTargetStore() *TargetStore {
	return &TargetStore{
		targets: make(map[string]DeviceTarget),
	}
}

// ReplaceAll replaces the entire store with new targets (Full Sync)
func (s *TargetStore) ReplaceAll(newTargets []DeviceTarget) {
	s.mu.Lock()
	defer s.mu.Unlock()

	s.targets = make(map[string]DeviceTarget)
	for _, t := range newTargets {
		s.targets[t.IP] = t
	}
}

// GetAll returns a copy of all targets
func (s *TargetStore) GetAll() []DeviceTarget {
	s.mu.RLock()
	defer s.mu.RUnlock()

	result := make([]DeviceTarget, 0, len(s.targets))
	for _, t := range s.targets {
		result = append(result, t)
	}
	return result
}
