package execution

import (
	"context"
	"fmt"
	"time"

	"github.com/yourorg/nms-go/internal/device/model"
	"github.com/yourorg/nms-go/internal/worker/protocols/mikrotik"
)

type ExecutionService interface {
	ExecuteCommand(ctx context.Context, req ExecuteCommandRequest) (*ExecuteCommandResponse, error)
	GetStats(ctx context.Context, req GetStatsRequest) (*GetStatsResponse, error)
}

type executionService struct {
}

func NewExecutionService() ExecutionService {
	return &executionService{}
}

func (s *executionService) ExecuteCommand(ctx context.Context, req ExecuteCommandRequest) (*ExecuteCommandResponse, error) {
	// 1. Create temporary device model
	device := &model.Device{
		ID:        "adhoc",
		IPAddress: req.Target.IP,
		Protocol:  model.ProtocolMikrotikAPI,
		Credentials: &model.DeviceCredentials{
			Username:          req.Target.Auth.Username,
			PasswordEncrypted: req.Target.Auth.Password, // Passing plain password as expected by current client implementation
		},
	}

	// 2. Select driver
	if req.Target.Driver != "mikrotik" {
		return nil, fmt.Errorf("unsupported driver: %s", req.Target.Driver)
	}

	// 3. Initiate Client
	client := mikrotik.NewMikrotikClient(10 * time.Second)

	// 4. Connect
	if err := client.Connect(ctx, device); err != nil {
		return &ExecuteCommandResponse{
			Status: "error",
			Error:  fmt.Sprintf("failed to connect: %v", err),
		}, nil
	}
	defer client.Disconnect()

	// 5. Execute
	output, err := client.ExecuteCommand(ctx, req.Command)
	if err != nil {
		return &ExecuteCommandResponse{
			Status: "error",
			Error:  fmt.Sprintf("execution failed: %v", err),
		}, nil
	}

	return &ExecuteCommandResponse{
		Status: "success",
		Output: output,
	}, nil
}

func (s *executionService) GetStats(ctx context.Context, req GetStatsRequest) (*GetStatsResponse, error) {
	// 1. Create temporary device model
	device := &model.Device{
		ID:        "adhoc",
		IPAddress: req.Target.IP,
		Protocol:  model.ProtocolMikrotikAPI,
		Credentials: &model.DeviceCredentials{
			Username:          req.Target.Auth.Username,
			PasswordEncrypted: req.Target.Auth.Password,
		},
	}

	// 2. Select driver
	if req.Target.Driver != "mikrotik" {
		return nil, fmt.Errorf("unsupported driver: %s", req.Target.Driver)
	}

	// 3. Initiate Client
	client := mikrotik.NewMikrotikClient(10 * time.Second)

	// 4. Connect
	if err := client.Connect(ctx, device); err != nil {
		return &GetStatsResponse{
			Status: "error",
			Error:  fmt.Sprintf("failed to connect: %v", err),
		}, nil
	}
	defer client.Disconnect()

	// 5. Get Metrics
	systemMetrics, err := client.GetSystemMetrics(ctx)
	if err != nil {
		return &GetStatsResponse{
			Status: "error",
			Error:  fmt.Sprintf("failed to get system metrics: %v", err),
		}, nil
	}

	interfaceMetrics, err := client.GetInterfaceMetrics(ctx)
	if err != nil {
		return &GetStatsResponse{
			Status: "error",
			Error:  fmt.Sprintf("failed to get interface metrics: %v", err),
		}, nil
	}

	// 6. Construct Response
	data := map[string]interface{}{
		"system":     systemMetrics,
		"interfaces": interfaceMetrics,
	}

	return &GetStatsResponse{
		Status: "success",
		Data:   data,
	}, nil
}
