package monitoring

import (
	"context"
	"fmt"
	"log"
	"time"

	"github.com/yourorg/nms-go/internal/device/model"
	"github.com/yourorg/nms-go/internal/worker/protocols/mikrotik"
)

// PollDevice connects to a device, gathers metrics, and writes them
func PollDevice(ctx context.Context, target DeviceTarget, writer MetricWriter) error {
	// Construct temporary device model
	device := &model.Device{
		ID:        target.IP, // using IP as ID for simplicity in ad-hoc polling
		IPAddress: target.IP,
		Protocol:  model.ProtocolMikrotikAPI,
		Credentials: &model.DeviceCredentials{
			Username:          target.Username,
			PasswordEncrypted: target.Password,
		},
	}

	client := mikrotik.NewMikrotikClient(10 * time.Second)
	if err := client.Connect(ctx, device); err != nil {
		return fmt.Errorf("failed to connect to %s: %w", target.IP, err)
	}
	defer client.Disconnect()

	// 1. Get System Metrics
	sysMetrics, err := client.GetSystemMetrics(ctx)
	if err != nil {
		log.Printf("Error collecting system metrics for %s: %v", target.IP, err)
	} else {
		writer.WriteSystemMetrics(sysMetrics)
	}

	// 2. Get Interface Metrics
	ifMetrics, err := client.GetInterfaceMetrics(ctx)
	if err != nil {
		log.Printf("Error collecting interface metrics for %s: %v", target.IP, err)
	} else {
		writer.WriteInterfaceMetrics(ifMetrics)
	}

	return nil
}
