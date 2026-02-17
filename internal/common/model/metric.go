package model

import "time"

// Metric represents a data point collected from a device
type Metric struct {
	DeviceID   string                 `json:"device_id"`
	DeviceName string                 `json:"device_name"` // Optional, for easier alerting
	IPAddress  string                 `json:"ip_address"`
	Timestamp  time.Time              `json:"timestamp"`
	Values     map[string]interface{} `json:"values"` // e.g. "cpu": 80.5, "rtt": 20.0
	Tags       map[string]string      `json:"tags,omitempty"`
}
