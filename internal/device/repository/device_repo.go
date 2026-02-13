package repository

import (
	"context"
	"fmt"

	"github.com/yourorg/nms-go/internal/device/model"
	"gorm.io/gorm"
)

// DeviceRepository defines the interface for device data access
type DeviceRepository interface {
	Create(ctx context.Context, device *model.Device) error
	GetByID(ctx context.Context, id string) (*model.Device, error)
	GetByIPAddress(ctx context.Context, ipAddress string) (*model.Device, error)
	List(ctx context.Context, filter *DeviceFilter) ([]*model.Device, error)
	Update(ctx context.Context, device *model.Device) error
	UpdateStatus(ctx context.Context, id string, status model.DeviceStatus) error
	Delete(ctx context.Context, id string) error
	Count(ctx context.Context, filter *DeviceFilter) (int64, error)
	GetByGroup(ctx context.Context, groupID string) ([]*model.Device, error)
	ListForPolling(ctx context.Context, limit int) ([]*model.Device, error)
}

// DeviceFilter represents filtering options for device queries
type DeviceFilter struct {
	DeviceType *model.DeviceType
	Protocol   *model.Protocol
	Status     *model.DeviceStatus
	GroupID    *string
	Tags       []string
	Enabled    *bool
	Search     string // Search in name, IP, description
	Limit      int
	Offset     int
}

type deviceRepository struct {
	db *gorm.DB
}

// NewDeviceRepository creates a new instance of DeviceRepository
func NewDeviceRepository(db *gorm.DB) DeviceRepository {
	return &deviceRepository{db: db}
}

// Create creates a new device
func (r *deviceRepository) Create(ctx context.Context, device *model.Device) error {
	return r.db.WithContext(ctx).Create(device).Error
}

// GetByID retrieves a device by ID with related data
func (r *deviceRepository) GetByID(ctx context.Context, id string) (*model.Device, error) {
	var device model.Device
	err := r.db.WithContext(ctx).
		Preload("Credentials").
		Preload("Group").
		First(&device, "id = ?", id).Error
	
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, fmt.Errorf("device not found: %s", id)
		}
		return nil, err
	}
	
	return &device, nil
}

// GetByIPAddress retrieves a device by IP address
func (r *deviceRepository) GetByIPAddress(ctx context.Context, ipAddress string) (*model.Device, error) {
	var device model.Device
	err := r.db.WithContext(ctx).
		Preload("Credentials").
		First(&device, "ip_address = ?", ipAddress).Error
	
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, fmt.Errorf("device not found with IP: %s", ipAddress)
		}
		return nil, err
	}
	
	return &device, nil
}

// List retrieves devices based on filter criteria
func (r *deviceRepository) List(ctx context.Context, filter *DeviceFilter) ([]*model.Device, error) {
	var devices []*model.Device
	
	query := r.db.WithContext(ctx).
		Preload("Credentials").
		Preload("Group")
	
	query = r.applyFilter(query, filter)
	
	if filter != nil {
		if filter.Limit > 0 {
			query = query.Limit(filter.Limit)
		}
		if filter.Offset > 0 {
			query = query.Offset(filter.Offset)
		}
	}
	
	err := query.Find(&devices).Error
	return devices, err
}

// Update updates an existing device
func (r *deviceRepository) Update(ctx context.Context, device *model.Device) error {
	return r.db.WithContext(ctx).
		Model(device).
		Updates(device).Error
}

// UpdateStatus updates only the status of a device
func (r *deviceRepository) UpdateStatus(ctx context.Context, id string, status model.DeviceStatus) error {
	return r.db.WithContext(ctx).
		Model(&model.Device{}).
		Where("id = ?", id).
		Update("status", status).Error
}

// Delete soft deletes a device
func (r *deviceRepository) Delete(ctx context.Context, id string) error {
	return r.db.WithContext(ctx).
		Delete(&model.Device{}, "id = ?", id).Error
}

// Count returns the total number of devices matching the filter
func (r *deviceRepository) Count(ctx context.Context, filter *DeviceFilter) (int64, error) {
	var count int64
	query := r.db.WithContext(ctx).Model(&model.Device{})
	query = r.applyFilter(query, filter)
	err := query.Count(&count).Error
	return count, err
}

// GetByGroup retrieves all devices in a specific group
func (r *deviceRepository) GetByGroup(ctx context.Context, groupID string) ([]*model.Device, error) {
	var devices []*model.Device
	err := r.db.WithContext(ctx).
		Preload("Credentials").
		Where("group_id = ?", groupID).
		Find(&devices).Error
	return devices, err
}

// ListForPolling retrieves enabled devices that are due for polling
func (r *deviceRepository) ListForPolling(ctx context.Context, limit int) ([]*model.Device, error) {
	var devices []*model.Device
	err := r.db.WithContext(ctx).
		Preload("Credentials").
		Where("enabled = ? AND status != ?", true, model.DeviceStatusError).
		Order("last_seen ASC NULLS FIRST").
		Limit(limit).
		Find(&devices).Error
	return devices, err
}

// applyFilter applies filter criteria to the query
func (r *deviceRepository) applyFilter(query *gorm.DB, filter *DeviceFilter) *gorm.DB {
	if filter == nil {
		return query
	}
	
	if filter.DeviceType != nil {
		query = query.Where("device_type = ?", *filter.DeviceType)
	}
	
	if filter.Protocol != nil {
		query = query.Where("protocol = ?", *filter.Protocol)
	}
	
	if filter.Status != nil {
		query = query.Where("status = ?", *filter.Status)
	}
	
	if filter.GroupID != nil {
		query = query.Where("group_id = ?", *filter.GroupID)
	}
	
	if filter.Enabled != nil {
		query = query.Where("enabled = ?", *filter.Enabled)
	}
	
	if len(filter.Tags) > 0 {
		query = query.Where("tags @> ?", filter.Tags)
	}
	
	if filter.Search != "" {
		searchPattern := "%" + filter.Search + "%"
		query = query.Where(
			"name ILIKE ? OR ip_address::text ILIKE ? OR description ILIKE ?",
			searchPattern, searchPattern, searchPattern,
		)
	}
	
	return query
}
