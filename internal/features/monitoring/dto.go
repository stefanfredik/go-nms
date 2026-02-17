package monitoring

import "github.com/yourorg/nms-go/internal/features/execution"

// SyncRequest represents the payload from OpenAccess to sync inventory
type SyncRequest struct {
	Targets []execution.Target `json:"targets" binding:"required"`
}

// DeviceTarget is the internal representation of a monitoring target
type DeviceTarget struct {
	IP       string
	Driver   string
	Username string
	Password string
	Port     int
}
