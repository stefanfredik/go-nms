package model

import "time"

// PollTask represents a task to poll a specific device
type PollTask struct {
	DeviceID   string      `json:"device_id"`
	IPAddress  string      `json:"ip_address"`
	DeviceType string      `json:"device_type"`
	Protocol   string      `json:"protocol"`
	Timestamp  time.Time   `json:"timestamp"`
}
