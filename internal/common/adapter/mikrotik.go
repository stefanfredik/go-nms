package adapter

import (
	"fmt"
	"strings"

	"github.com/go-routeros/routeros"
)

type MikrotikAdapter struct{}

func NewMikrotikAdapter() *MikrotikAdapter {
	return &MikrotikAdapter{}
}

// FetchSystemResources connects to the Mikrotik device and retrieves system resource data.
// Returns a map of metrics and true if successful, or nil and false if failed.
func (m *MikrotikAdapter) FetchSystemResources(ip, username, password string) (map[string]interface{}, bool) {
	// Default API port is 8728
	address := fmt.Sprintf("%s:8728", ip)

	// Dial the device
	c, err := routeros.Dial(address, username, password)
	if err != nil {
		return nil, false
	}
	defer c.Close()

	// Execute /system/resource/print
	reply, err := c.Run("/system/resource/print")
	if err != nil {
		return nil, false
	}

	if len(reply.Re) == 0 {
		return nil, false
	}

	// Parse the first result
	res := reply.Re[0].Map

	metrics := make(map[string]interface{})

	if val, ok := res["uptime"]; ok {
		metrics["uptime_str"] = val
	}

	if val, ok := res["cpu-load"]; ok {
		metrics["cpu_load"] = parsePercentage(val)
	}

	if val, ok := res["free-memory"]; ok {
		metrics["free_memory"] = parseBytes(val)
	}

	if val, ok := res["total-memory"]; ok {
		metrics["total_memory"] = parseBytes(val)
	}

	return metrics, true
}

// RunCommand executes a command via Mikrotik API
func (m *MikrotikAdapter) RunCommand(ip, username, password, command string) (string, error) {
	address := fmt.Sprintf("%s:8728", ip)
	c, err := routeros.Dial(address, username, password)
	if err != nil {
		return "", fmt.Errorf("failed to dial mikrotik: %w", err)
	}
	defer c.Close()

	// Mikrotik API commands are structured differently than SSH CLI.
	// Simple command like "/system/resource/print" works, but complex ones might need parsing.
	// We'll attempt to run it directly.
	// Note: API command format is different from CLI. CLI: /system resource print. API: /system/resource/print.
	// We might need to normalize the command for API if user inputs CLI format.

	// Normalize: replace spaces with /? Or assume user knows API format?
	// For "coba ambil resource", the command "/system/resource/print" is expected for API.
	// But user prompt might be generic. Let's try to handle basic conversion or just run it.

	// Basic CLI to API conversion (very naive)
	apiCommand := command
	if !strings.Contains(command, "/") {
		// If it's just words separated by space, try adding /
		apiCommand = "/" + strings.ReplaceAll(strings.TrimSpace(command), " ", "/")
	} else if strings.Contains(command, " ") && !strings.Contains(command, "/ ") {
		// If it has spaces but is a path like "/system resource print", convert to "/system/resource/print"
		apiCommand = strings.ReplaceAll(command, " ", "/")
	}

	reply, err := c.Run(apiCommand)
	if err != nil {
		return "", fmt.Errorf("failed to run command: %w", err)
	}

	// Format output from Reply
	var output strings.Builder
	for _, re := range reply.Re {
		for k, v := range re.Map {
			output.WriteString(fmt.Sprintf("%s=%s ", k, v))
		}
		output.WriteString("\n")
	}

	return output.String(), nil
}

func parsePercentage(s string) float64 {
	var f float64
	fmt.Sscanf(strings.TrimSuffix(s, "%"), "%f", &f)
	return f
}

func parseBytes(s string) int64 {
	var i int64
	fmt.Sscanf(s, "%d", &i)
	return i
}

// RunCommandStructured executes a command and returns the raw result map
func (m *MikrotikAdapter) RunCommandStructured(ip, username, password, command string) ([]map[string]string, error) {
	address := fmt.Sprintf("%s:8728", ip)
	c, err := routeros.Dial(address, username, password)
	if err != nil {
		return nil, fmt.Errorf("failed to dial mikrotik: %w", err)
	}
	defer c.Close()

	// Normalize command (same logic as RunCommand)
	apiCommand := command
	if !strings.Contains(command, "/") {
		apiCommand = "/" + strings.ReplaceAll(strings.TrimSpace(command), " ", "/")
	} else if strings.Contains(command, " ") && !strings.Contains(command, "/ ") {
		apiCommand = strings.ReplaceAll(command, " ", "/")
	}

	reply, err := c.Run(apiCommand)
	if err != nil {
		return nil, fmt.Errorf("failed to run command: %w", err)
	}

	var results []map[string]string
	for _, re := range reply.Re {
		results = append(results, re.Map)
	}

	return results, nil
}
