package service

import (
	"context"
	"errors"
	"time"

	"github.com/yourorg/nms-go/internal/device/model"
	"github.com/yourorg/nms-go/internal/device/repository"
)

type DeviceService interface {
	RegisterDevice(ctx context.Context, req *RegisterDeviceRequest) (*model.Device, error)
	GetDevice(ctx context.Context, id string) (*model.Device, error)
	ListDevices(ctx context.Context, page, pageSize int) ([]*model.Device, int64, error)
}

type deviceService struct {
	repo repository.DeviceRepository
}

func NewDeviceService(repo repository.DeviceRepository) DeviceService {
	return &deviceService{repo: repo}
}

type RegisterDeviceRequest struct {
	Name            string             `json:"name"`
	IPAddress       string             `json:"ip_address"`
	DeviceType      model.DeviceType   `json:"device_type"`
	Protocol        model.Protocol     `json:"protocol"`
	PollingInterval int                `json:"polling_interval"`
	Tags            []string           `json:"tags"`
}

func (s *deviceService) RegisterDevice(ctx context.Context, req *RegisterDeviceRequest) (*model.Device, error) {
	// Check if device with same IP already exists
	existing, _ := s.repo.GetByIPAddress(ctx, req.IPAddress)
	if existing != nil {
		return nil, errors.New("device with this IP address already exists")
	}

	device := &model.Device{
		Name:            req.Name,
		IPAddress:       req.IPAddress,
		DeviceType:      req.DeviceType,
		Protocol:        req.Protocol,
		PollingInterval: req.PollingInterval,
		Tags:            req.Tags,
		Status:          model.DeviceStatusUnknown,
		CreatedAt:       time.Now(),
		UpdatedAt:       time.Now(),
		Enabled:         true,
	}

	if device.PollingInterval == 0 {
		device.PollingInterval = 300 // Default 5 mins
	}

	err := s.repo.Create(ctx, device)
	if err != nil {
		return nil, err
	}

	return device, nil
}

func (s *deviceService) GetDevice(ctx context.Context, id string) (*model.Device, error) {
	return s.repo.GetByID(ctx, id)
}

func (s *deviceService) ListDevices(ctx context.Context, page, pageSize int) ([]*model.Device, int64, error) {
	if page < 1 {
		page = 1
	}
	if pageSize < 1 || pageSize > 100 {
		pageSize = 20
	}
	offset := (page - 1) * pageSize
	
	filter := &repository.DeviceFilter{
		Limit:  pageSize,
		Offset: offset,
	}
	
	devices, err := s.repo.List(ctx, filter)
	if err != nil {
		return nil, 0, err
	}
	
	count, err := s.repo.Count(ctx, filter)
	if err != nil {
		return nil, 0, err
	}
	
	return devices, count, nil
}
