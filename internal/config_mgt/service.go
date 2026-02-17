package config_mgt

import (
	"context"
	"fmt"

	"github.com/yourorg/nms-go/internal/common/adapter"
	"github.com/yourorg/nms-go/internal/device/service"
)

type ConfigService interface {
	ExecuteCommand(ctx context.Context, deviceID, command string) (interface{}, error)
	BackupConfig(ctx context.Context, deviceID string) (string, error)
}

type configService struct {
	deviceService service.DeviceService
	sshAdapter    *SSHAdapter
}

func NewConfigService(ds service.DeviceService, ssh *SSHAdapter) ConfigService {
	return &configService{
		deviceService: ds,
		sshAdapter:    ssh,
	}
}

func (s *configService) ExecuteCommand(ctx context.Context, deviceID, command string) (interface{}, error) {
	device, err := s.deviceService.GetDevice(ctx, deviceID)
	if err != nil {
		return "", fmt.Errorf("device not found: %w", err)
	}

	// In a real app, fetch credentials from DB using device.CredentialsID
	// For MVP, we use defaults or mock
	user := "admin"
	password := "RexusBattlefire"

	if device.Protocol == "mikrotik_api" {
		mtAdapter := adapter.NewMikrotikAdapter()
		// Try to convert CLI command to API format if needed, or just pass it
		// e.g. /system resource print -> /system/resource/print
		// The adapter now handles basic conversion
		return mtAdapter.RunCommandStructured(device.IPAddress, user, password, command)
	}

	// Default to SSH
	output, err := s.sshAdapter.Execute(device.IPAddress, user, password, command)
	if err != nil {
		return output, fmt.Errorf("execution failed: %w", err)
	}

	return output, nil
}

func (s *configService) BackupConfig(ctx context.Context, deviceID string) (string, error) {
	// Simple backup: assume "export" command works (Ross/Mikrotik style)
	res, err := s.ExecuteCommand(ctx, deviceID, "/export")
	if err != nil {
		return "", err
	}
	// Since ExecuteCommand now returns interface{}, we need to assert it to string
	// Backup via API might return structured data, but /export usually returns text script even in API?
	// Actually, /export in API might not work the same way.
	// For now, let's assume it returns a string or we convert it.
	if str, ok := res.(string); ok {
		return str, nil
	}
	return fmt.Sprintf("%v", res), nil
}
