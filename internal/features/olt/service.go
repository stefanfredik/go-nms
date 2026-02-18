package olt

import (
	"context"
	"fmt"
	"time"

	devicemodel "github.com/yourorg/nms-go/internal/device/model"
	"github.com/yourorg/nms-go/internal/worker/protocols/snmp/zte"
)

// OLTService defines the interface for querying OLT device data via SNMP.
// It is consumed by the HTTP handler and can be mocked in tests.
//
// All methods accept an SNMPTarget instead of a device ID, so go-nms does not
// need its own device registry — openaccess is the single source of truth.
type OLTService interface {
	// GetSystemMetrics returns system-level metrics for the OLT at the given target.
	GetSystemMetrics(ctx context.Context, target SNMPTarget) (*SystemMetricsResponse, error)

	// GetPONPorts returns metrics for all PON ports on the OLT at the given target.
	GetPONPorts(ctx context.Context, target SNMPTarget) (*PONPortListResponse, error)

	// GetONTs returns metrics for all ONTs on the OLT at the given target.
	// If ponPortIndex > 0, only ONTs on that specific PON port are returned.
	GetONTs(ctx context.Context, target SNMPTarget, ponPortIndex int) (*ONTListResponse, error)
}

type oltService struct {
	timeout time.Duration
}

// NewOLTService creates a new OLTService.
// No device repository is needed — connection details come from the request body.
func NewOLTService() OLTService {
	return &oltService{
		timeout: 15 * time.Second,
	}
}

// GetSystemMetrics retrieves system metrics from the OLT via SNMP.
func (s *oltService) GetSystemMetrics(ctx context.Context, target SNMPTarget) (*SystemMetricsResponse, error) {
	client, err := s.connectToOLT(ctx, target)
	if err != nil {
		return nil, err
	}
	defer client.Disconnect()

	metrics, err := client.GetSystemMetrics(ctx)
	if err != nil {
		return nil, fmt.Errorf("failed to get system metrics from OLT %s: %w", target.IP, err)
	}

	return mapSystemMetrics(target.IP, metrics), nil
}

// GetPONPorts retrieves PON port metrics from the OLT via SNMP.
func (s *oltService) GetPONPorts(ctx context.Context, target SNMPTarget) (*PONPortListResponse, error) {
	client, err := s.connectToOLT(ctx, target)
	if err != nil {
		return nil, err
	}
	defer client.Disconnect()

	ports, err := client.GetPONPortMetrics(ctx)
	if err != nil {
		return nil, fmt.Errorf("failed to get PON port metrics from OLT %s: %w", target.IP, err)
	}

	responses := make([]PONPortResponse, 0, len(ports))
	for _, p := range ports {
		responses = append(responses, mapPONPort(target.IP, p))
	}

	return &PONPortListResponse{
		IPAddress: target.IP,
		Count:     len(responses),
		PonPorts:  responses,
	}, nil
}

// GetONTs retrieves ONT metrics from the OLT via SNMP.
func (s *oltService) GetONTs(ctx context.Context, target SNMPTarget, ponPortIndex int) (*ONTListResponse, error) {
	client, err := s.connectToOLT(ctx, target)
	if err != nil {
		return nil, err
	}
	defer client.Disconnect()

	onts, err := client.GetONTMetrics(ctx, ponPortIndex)
	if err != nil {
		return nil, fmt.Errorf("failed to get ONT metrics from OLT %s: %w", target.IP, err)
	}

	responses := make([]ONTResponse, 0, len(onts))
	for _, o := range onts {
		responses = append(responses, mapONT(target.IP, o))
	}

	return &ONTListResponse{
		IPAddress: target.IP,
		Total:     len(responses),
		ONTs:      responses,
	}, nil
}

// connectToOLT builds a synthetic device model from the SNMPTarget and
// establishes an SNMP session. No database lookup is required.
func (s *oltService) connectToOLT(ctx context.Context, target SNMPTarget) (*zte.ZTEOLTClient, error) {
	community := target.Community
	if community == "" {
		community = "public"
	}

	// Build a lightweight synthetic device so we can reuse the ZTE client as-is.
	device := &devicemodel.Device{
		ID:         target.IP, // use IP as identifier for metrics labelling
		IPAddress:  target.IP,
		DeviceType: devicemodel.DeviceTypeOLT,
		Protocol:   devicemodel.ProtocolSNMP,
		Credentials: &devicemodel.DeviceCredentials{
			SNMPCommunity: community,
		},
	}

	client := zte.NewZTEOLTClient(s.timeout)
	if err := client.Connect(ctx, device); err != nil {
		return nil, fmt.Errorf("failed to connect to OLT %s via SNMP: %w", target.IP, err)
	}

	return client, nil
}

// ── Mapping helpers ───────────────────────────────────────────────────────────

func mapSystemMetrics(ip string, m *zte.OLTSystemMetrics) *SystemMetricsResponse {
	return &SystemMetricsResponse{
		IPAddress:          ip,
		Timestamp:          m.Timestamp,
		SysDescr:           m.SysDescr,
		SysName:            m.SysName,
		UptimeSeconds:      m.UptimeSeconds,
		CPUUsagePercent:    m.CPUUsagePercent,
		MemoryTotalKB:      m.MemoryTotalKB,
		MemoryUsedKB:       m.MemoryUsedKB,
		MemoryUsagePercent: m.MemoryUsagePercent,
		TemperatureCelsius: m.TemperatureCelsius,
	}
}

func mapPONPort(ip string, p *zte.PONPortMetrics) PONPortResponse {
	return PONPortResponse{
		IPAddress:   ip,
		Timestamp:   p.Timestamp,
		PortIndex:   p.PortIndex,
		AdminStatus: p.AdminStatus.String(),
		OperStatus:  p.OperStatus.String(),
		TxPowerDBm:  p.TxPowerDBm,
		RxPowerDBm:  p.RxPowerDBm,
		ONTCount:    p.ONTCount,
	}
}

func mapONT(ip string, o *zte.ONTMetrics) ONTResponse {
	return ONTResponse{
		IPAddress:      ip,
		Timestamp:      o.Timestamp,
		PONPortIndex:   o.PONPortIndex,
		ONTIndex:       o.ONTIndex,
		SerialNumber:   o.SerialNumber,
		OperStatus:     o.OperStatus.String(),
		RxPowerDBm:     o.RxPowerDBm,
		TxPowerDBm:     o.TxPowerDBm,
		DistanceMeters: o.DistanceMeters,
		Description:    o.Description,
	}
}
