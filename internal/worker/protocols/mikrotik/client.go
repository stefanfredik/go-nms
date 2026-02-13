package mikrotik

import (
	"context"
	"fmt"
	"strconv"
	"time"

	"github.com/go-routeros/routeros"
	"github.com/yourorg/nms-go/internal/device/model"
)

// SystemMetrics represents system-level metrics from a device
type SystemMetrics struct {
	DeviceID    string
	Timestamp   time.Time
	CPUUsage    float64
	MemoryTotal uint64
	MemoryFree  uint64
	MemoryUsage float64
	DiskTotal   uint64
	DiskFree    uint64
	DiskUsage   float64
	Uptime      int64 // seconds
	Temperature float64
	Voltage     float64
}

// InterfaceMetrics represents interface-level metrics
type InterfaceMetrics struct {
	DeviceID      string
	InterfaceName string
	Timestamp     time.Time
	Status        string // up, down, disabled
	BytesIn       uint64
	BytesOut      uint64
	PacketsIn     uint64
	PacketsOut    uint64
	ErrorsIn      uint64
	ErrorsOut     uint64
	DropsIn       uint64
	DropsOut      uint64
	Speed         string // 100Mbps, 1Gbps, etc
}

// MikrotikClient implements protocol.DeviceProtocol for Mikrotik devices
type MikrotikClient struct {
	client  *routeros.Client
	device  *model.Device
	timeout time.Duration
}

// NewMikrotikClient creates a new Mikrotik API client
func NewMikrotikClient(timeout time.Duration) *MikrotikClient {
	return &MikrotikClient{
		timeout: timeout,
	}
}

// Connect establishes connection to Mikrotik device
func (m *MikrotikClient) Connect(ctx context.Context, device *model.Device) error {
	m.device = device
	
	// Validate credentials
	if device.Credentials == nil {
		return fmt.Errorf("device credentials not loaded")
	}
	
	// Connect with timeout
	dialCtx, cancel := context.WithTimeout(ctx, m.timeout)
	defer cancel()
	
	address := fmt.Sprintf("%s:8728", device.IPAddress) // Default Mikrotik API port
	
	client, err := routeros.DialTimeout(address, device.Credentials.Username, 
		device.Credentials.PasswordEncrypted, m.timeout) // TODO: Decrypt password
	if err != nil {
		return fmt.Errorf("failed to connect: %w", err)
	}
	
	m.client = client
	return nil
}

// Disconnect closes the connection
func (m *MikrotikClient) Disconnect() error {
	if m.client != nil {
		m.client.Close()
	}
	return nil
}

// ExecuteCommand executes a RouterOS command
func (m *MikrotikClient) ExecuteCommand(ctx context.Context, command string) (string, error) {
	if m.client == nil {
		return "", fmt.Errorf("not connected")
	}
	
	reply, err := m.client.Run(command)
	if err != nil {
		return "", fmt.Errorf("command execution failed: %w", err)
	}
	
	// Convert reply to string representation
	result := ""
	for _, sentence := range reply.Re {
		for key, value := range sentence.Map {
			result += fmt.Sprintf("%s: %s\n", key, value)
		}
	}
	
	return result, nil
}

// GetSystemMetrics retrieves system-level metrics
func (m *MikrotikClient) GetSystemMetrics(ctx context.Context) (*SystemMetrics, error) {
	if m.client == nil {
		return nil, fmt.Errorf("not connected")
	}
	
	metrics := &SystemMetrics{
		DeviceID:  m.device.ID,
		Timestamp: time.Now(),
	}
	
	// Get system resources
	reply, err := m.client.Run("/system/resource/print")
	if err != nil {
		return nil, fmt.Errorf("failed to get system resources: %w", err)
	}
	
	if len(reply.Re) == 0 {
		return nil, fmt.Errorf("no system resource data returned")
	}
	
	res := reply.Re[0].Map
	
	// Parse CPU usage
	if cpuLoad, ok := res["cpu-load"]; ok {
		if cpu, err := strconv.ParseFloat(cpuLoad, 64); err == nil {
			metrics.CPUUsage = cpu
		}
	}
	
	// Parse memory
	if totalMem, ok := res["total-memory"]; ok {
		if mem, err := strconv.ParseUint(totalMem, 10, 64); err == nil {
			metrics.MemoryTotal = mem
		}
	}
	
	if freeMem, ok := res["free-memory"]; ok {
		if mem, err := strconv.ParseUint(freeMem, 10, 64); err == nil {
			metrics.MemoryFree = mem
		}
	}
	
	// Calculate memory usage percentage
	if metrics.MemoryTotal > 0 {
		used := metrics.MemoryTotal - metrics.MemoryFree
		metrics.MemoryUsage = float64(used) / float64(metrics.MemoryTotal) * 100
	}
	
	// Parse disk
	if totalHdd, ok := res["total-hdd-space"]; ok {
		if disk, err := strconv.ParseUint(totalHdd, 10, 64); err == nil {
			metrics.DiskTotal = disk
		}
	}
	
	if freeHdd, ok := res["free-hdd-space"]; ok {
		if disk, err := strconv.ParseUint(freeHdd, 10, 64); err == nil {
			metrics.DiskFree = disk
		}
	}
	
	// Calculate disk usage percentage
	if metrics.DiskTotal > 0 {
		used := metrics.DiskTotal - metrics.DiskFree
		metrics.DiskUsage = float64(used) / float64(metrics.DiskTotal) * 100
	}
	
	// Parse uptime
	if uptime, ok := res["uptime"]; ok {
		// Uptime format: 1w2d3h4m5s - need to parse this
		duration := parseRouterOSUptime(uptime)
		metrics.Uptime = int64(duration.Seconds())
	}
	
	// Get health info
	healthReply, err := m.client.Run("/system/health/print")
	if err == nil && len(healthReply.Re) > 0 {
		health := healthReply.Re[0].Map
		
		if temp, ok := health["temperature"]; ok {
			if t, err := strconv.ParseFloat(temp, 64); err == nil {
				metrics.Temperature = t
			}
		}
		
		if voltage, ok := health["voltage"]; ok {
			if v, err := strconv.ParseFloat(voltage, 64); err == nil {
				metrics.Voltage = v
			}
		}
	}
	
	return metrics, nil
}

// GetInterfaceMetrics retrieves metrics for all interfaces
func (m *MikrotikClient) GetInterfaceMetrics(ctx context.Context) ([]*InterfaceMetrics, error) {
	if m.client == nil {
		return nil, fmt.Errorf("not connected")
	}
	
	// Get interface statistics
	reply, err := m.client.Run("/interface/print", "=stats")
	if err != nil {
		return nil, fmt.Errorf("failed to get interface stats: %w", err)
	}
	
	var metrics []*InterfaceMetrics
	timestamp := time.Now()
	
	for _, iface := range reply.Re {
		m := &InterfaceMetrics{
			DeviceID:      m.device.ID,
			InterfaceName: iface.Map["name"],
			Timestamp:     timestamp,
			Status:        iface.Map["running"],
		}
		
		// Parse counters
		if val, ok := iface.Map["rx-byte"]; ok {
			if n, err := strconv.ParseUint(val, 10, 64); err == nil {
				m.BytesIn = n
			}
		}
		
		if val, ok := iface.Map["tx-byte"]; ok {
			if n, err := strconv.ParseUint(val, 10, 64); err == nil {
				m.BytesOut = n
			}
		}
		
		if val, ok := iface.Map["rx-packet"]; ok {
			if n, err := strconv.ParseUint(val, 10, 64); err == nil {
				m.PacketsIn = n
			}
		}
		
		if val, ok := iface.Map["tx-packet"]; ok {
			if n, err := strconv.ParseUint(val, 10, 64); err == nil {
				m.PacketsOut = n
			}
		}
		
		if val, ok := iface.Map["rx-error"]; ok {
			if n, err := strconv.ParseUint(val, 10, 64); err == nil {
				m.ErrorsIn = n
			}
		}
		
		if val, ok := iface.Map["tx-error"]; ok {
			if n, err := strconv.ParseUint(val, 10, 64); err == nil {
				m.ErrorsOut = n
			}
		}
		
		if val, ok := iface.Map["rx-drop"]; ok {
			if n, err := strconv.ParseUint(val, 10, 64); err == nil {
				m.DropsIn = n
			}
		}
		
		if val, ok := iface.Map["tx-drop"]; ok {
			if n, err := strconv.ParseUint(val, 10, 64); err == nil {
				m.DropsOut = n
			}
		}
		
		metrics = append(metrics, m)
	}
	
	return metrics, nil
}

// GetWirelessMetrics retrieves wireless-specific metrics
func (m *MikrotikClient) GetWirelessMetrics(ctx context.Context) ([]map[string]interface{}, error) {
	if m.client == nil {
		return nil, fmt.Errorf("not connected")
	}
	
	// Get wireless interfaces
	reply, err := m.client.Run("/interface/wireless/print")
	if err != nil {
		return nil, fmt.Errorf("failed to get wireless interfaces: %w", err)
	}
	
	var metrics []map[string]interface{}
	
	for _, iface := range reply.Re {
		metric := make(map[string]interface{})
		metric["device_id"] = m.device.ID
		metric["interface"] = iface.Map["name"]
		metric["ssid"] = iface.Map["ssid"]
		metric["frequency"] = iface.Map["frequency"]
		metric["band"] = iface.Map["band"]
		
		// Get registration table for client count
		regReply, err := m.client.Run("/interface/wireless/registration-table/print",
			fmt.Sprintf("?interface=%s", iface.Map["name"]))
		if err == nil {
			metric["connected_clients"] = len(regReply.Re)
		}
		
		metrics = append(metrics, metric)
	}
	
	return metrics, nil
}

// parseRouterOSUptime parses RouterOS uptime format (e.g., "1w2d3h4m5s")
func parseRouterOSUptime(uptime string) time.Duration {
	// This is a simplified parser - production code should be more robust
	var duration time.Duration
	
	// Parse weeks, days, hours, minutes, seconds
	// Implementation depends on actual format received from device
	// This is a placeholder
	
	return duration
}

// ValidateConnection performs a quick connection test
func (m *MikrotikClient) ValidateConnection(ctx context.Context, device *model.Device) error {
	if err := m.Connect(ctx, device); err != nil {
		return err
	}
	defer m.Disconnect()
	
	// Try to get system identity as a validation
	reply, err := m.client.Run("/system/identity/print")
	if err != nil {
		return fmt.Errorf("connection validation failed: %w", err)
	}
	
	if len(reply.Re) == 0 {
		return fmt.Errorf("no response from device")
	}
	
	return nil
}
