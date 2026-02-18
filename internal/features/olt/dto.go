// Package olt provides the service and HTTP handler for querying OLT device data.
// It exposes OLT metrics (system, PON ports, ONTs) via REST API for consumption
// by openaccess and nms-rekayasa.
//
// All endpoints accept a JSON request body containing the SNMP target details.
// This means go-nms does NOT need its own device database — openaccess is the
// single source of truth for device inventory.
package olt

import "time"

// SNMPTarget describes the SNMP connection parameters for an OLT device.
// This is included in the request body of all OLT API endpoints.
type SNMPTarget struct {
	// IP is the management IP address of the OLT (required).
	IP string `json:"ip" binding:"required"`

	// Community is the SNMP v2c community string (default: "public").
	Community string `json:"community"`

	// Version is the SNMP version string: "2c" or "3" (default: "2c").
	// Currently only "2c" is supported.
	Version string `json:"version"`

	// Port is the SNMP UDP port (default: 161).
	Port uint16 `json:"port"`
}

// GetSystemMetricsRequest is the request body for POST /api/v1/olt/system.
type GetSystemMetricsRequest struct {
	Target SNMPTarget `json:"target" binding:"required"`
}

// GetPONPortsRequest is the request body for POST /api/v1/olt/pon-ports.
type GetPONPortsRequest struct {
	Target SNMPTarget `json:"target" binding:"required"`
}

// GetONTsRequest is the request body for POST /api/v1/olt/onts.
type GetONTsRequest struct {
	Target SNMPTarget `json:"target" binding:"required"`

	// PONPort filters results to a specific PON port index.
	// Set to 0 (or omit) to return ONTs from all PON ports.
	PONPort int `json:"pon_port"`
}

// SystemMetricsResponse is the API response for OLT system metrics.
type SystemMetricsResponse struct {
	IPAddress          string    `json:"ip_address"`
	Timestamp          time.Time `json:"timestamp"`
	SysDescr           string    `json:"sys_descr"`
	SysName            string    `json:"sys_name"`
	UptimeSeconds      int64     `json:"uptime_seconds"`
	CPUUsagePercent    float64   `json:"cpu_usage_percent"`
	MemoryTotalKB      int64     `json:"memory_total_kb"`
	MemoryUsedKB       int64     `json:"memory_used_kb"`
	MemoryUsagePercent float64   `json:"memory_usage_percent"`
	TemperatureCelsius float64   `json:"temperature_celsius"`
}

// PONPortResponse is the API response for a single PON port.
type PONPortResponse struct {
	IPAddress   string    `json:"ip_address"`
	Timestamp   time.Time `json:"timestamp"`
	PortIndex   int       `json:"port_index"`
	AdminStatus string    `json:"admin_status"`
	OperStatus  string    `json:"oper_status"`
	TxPowerDBm  float64   `json:"tx_power_dbm"`
	RxPowerDBm  float64   `json:"rx_power_dbm"`
	ONTCount    int       `json:"ont_count"`
}

// ONTResponse is the API response for a single ONT.
type ONTResponse struct {
	IPAddress      string    `json:"ip_address"`
	Timestamp      time.Time `json:"timestamp"`
	PONPortIndex   int       `json:"pon_port_index"`
	ONTIndex       int       `json:"ont_index"`
	SerialNumber   string    `json:"serial_number"`
	OperStatus     string    `json:"oper_status"`
	RxPowerDBm     float64   `json:"rx_power_dbm"`
	TxPowerDBm     float64   `json:"tx_power_dbm"`
	DistanceMeters int       `json:"distance_meters"`
	Description    string    `json:"description"`
}

// PONPortListResponse wraps a list of PON port responses.
type PONPortListResponse struct {
	IPAddress string            `json:"ip_address"`
	Count     int               `json:"count"`
	PonPorts  []PONPortResponse `json:"pon_ports"`
}

// ONTListResponse wraps a list of ONT responses.
type ONTListResponse struct {
	IPAddress string        `json:"ip_address"`
	Total     int           `json:"total"`
	ONTs      []ONTResponse `json:"onts"`
}
