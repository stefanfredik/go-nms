package service

import (
	"context"
	"fmt"
	"net"
	"os/exec"
	"sync"

	"github.com/yourorg/nms-go/internal/device/model"
)

type DiscoveryService interface {
	ScanSubnet(ctx context.Context, cidr string) ([]*model.Device, error)
}

type discoveryService struct {
}

func NewDiscoveryService() DiscoveryService {
	return &discoveryService{}
}

func (s *discoveryService) ScanSubnet(ctx context.Context, cidr string) ([]*model.Device, error) {
	ip, ipnet, err := net.ParseCIDR(cidr)
	if err != nil {
		return nil, fmt.Errorf("invalid CIDR: %w", err)
	}

	var devices []*model.Device
	var mu sync.Mutex
	var wg sync.WaitGroup
	sem := make(chan struct{}, 50) // Limit concurrency

	for ip := ip.Mask(ipnet.Mask); ipnet.Contains(ip); inc(ip) {
		currentIP := ip.String()
		wg.Add(1)
		sem <- struct{}{}

		go func(targetIP string) {
			defer wg.Done()
			defer func() { <-sem }()

			if checkPing(targetIP) {
				mu.Lock()
				devices = append(devices, &model.Device{
					Name:       fmt.Sprintf("Discovered Device %s", targetIP),
					IPAddress:  targetIP,
					DeviceType: model.DeviceTypeSwitch, // Default guess
					Status:     model.DeviceStatusOnline,
				})
				mu.Unlock()
			}
		}(currentIP)
	}

	wg.Wait()
	return devices, nil
}

func inc(ip net.IP) {
	for j := len(ip) - 1; j >= 0; j-- {
		ip[j]++
		if ip[j] > 0 {
			break
		}
	}
}

func checkPing(ip string) bool {
	// Simple ping command wrapper
	// Note: This relies on system 'ping' command
	cmd := exec.Command("ping", "-c", "1", "-W", "1", ip)
	err := cmd.Run()
	return err == nil
}

// Simple fingerprinting logic (future enhancement)
func fingerprint(ip string) model.DeviceType {
	// Try SSH banner grabbing or SNMP OID check
	return model.DeviceTypeSwitch
}
