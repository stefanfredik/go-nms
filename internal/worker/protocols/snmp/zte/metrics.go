package zte

import "time"

// OLTSystemMetrics holds system-level metrics collected from a ZTE C320 OLT.
type OLTSystemMetrics struct {
	// DeviceID is the identifier of the OLT device in go-nms.
	DeviceID string `json:"device_id"`

	// Timestamp is when the metrics were collected.
	Timestamp time.Time `json:"timestamp"`

	// SysDescr is the textual description of the OLT.
	SysDescr string `json:"sys_descr"`

	// SysName is the administratively-assigned name of the OLT.
	SysName string `json:"sys_name"`

	// UptimeSeconds is the time in seconds since the OLT was last restarted.
	UptimeSeconds int64 `json:"uptime_seconds"`

	// CPUUsagePercent is the current CPU utilization (0-100).
	CPUUsagePercent float64 `json:"cpu_usage_percent"`

	// MemoryTotalKB is the total installed memory in kilobytes.
	MemoryTotalKB int64 `json:"memory_total_kb"`

	// MemoryUsedKB is the currently used memory in kilobytes.
	MemoryUsedKB int64 `json:"memory_used_kb"`

	// MemoryUsagePercent is the calculated memory utilization (0-100).
	MemoryUsagePercent float64 `json:"memory_usage_percent"`

	// TemperatureCelsius is the chassis temperature in degrees Celsius.
	TemperatureCelsius float64 `json:"temperature_celsius"`
}

// PONPortMetrics holds metrics for a single PON port on a ZTE C320 OLT.
type PONPortMetrics struct {
	// DeviceID is the identifier of the parent OLT device.
	DeviceID string `json:"device_id"`

	// Timestamp is when the metrics were collected.
	Timestamp time.Time `json:"timestamp"`

	// PortIndex is the SNMP index of the PON port.
	PortIndex int `json:"port_index"`

	// AdminStatus is the administrative status of the port.
	AdminStatus PONPortStatus `json:"admin_status"`

	// OperStatus is the operational (actual) status of the port.
	OperStatus PONPortStatus `json:"oper_status"`

	// TxPowerDBm is the transmit optical power in dBm.
	TxPowerDBm float64 `json:"tx_power_dbm"`

	// RxPowerDBm is the receive optical power in dBm.
	RxPowerDBm float64 `json:"rx_power_dbm"`

	// ONTCount is the number of registered ONTs on this PON port.
	ONTCount int `json:"ont_count"`
}

// ONTMetrics holds metrics for a single ONT registered on a ZTE C320 OLT.
type ONTMetrics struct {
	// DeviceID is the identifier of the parent OLT device.
	DeviceID string `json:"device_id"`

	// Timestamp is when the metrics were collected.
	Timestamp time.Time `json:"timestamp"`

	// PONPortIndex is the index of the PON port this ONT is connected to.
	PONPortIndex int `json:"pon_port_index"`

	// ONTIndex is the SNMP index of the ONT within its PON port.
	ONTIndex int `json:"ont_index"`

	// SerialNumber is the factory serial number of the ONT.
	SerialNumber string `json:"serial_number"`

	// OperStatus is the current operational status of the ONT.
	OperStatus ONTStatus `json:"oper_status"`

	// RxPowerDBm is the receive optical power measured at the OLT side in dBm.
	// A value below -27 dBm typically indicates a signal problem.
	RxPowerDBm float64 `json:"rx_power_dbm"`

	// TxPowerDBm is the transmit optical power of the ONT in dBm.
	TxPowerDBm float64 `json:"tx_power_dbm"`

	// DistanceMeters is the physical distance from the OLT to the ONT in meters.
	DistanceMeters int `json:"distance_meters"`

	// Description is the user-configured description of the ONT.
	Description string `json:"description"`
}
