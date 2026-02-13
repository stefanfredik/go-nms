package model

import (
	"time"
)

// DeviceType represents the type of network device
type DeviceType string

const (
	DeviceTypeRouter   DeviceType = "router"
	DeviceTypeSwitch   DeviceType = "switch"
	DeviceTypeOLT      DeviceType = "olt"
	DeviceTypeONT      DeviceType = "ont"
	DeviceTypeAP       DeviceType = "access_point"
	DeviceTypeWireless DeviceType = "wireless"
)

// Protocol represents the communication protocol
type Protocol string

const (
	ProtocolMikrotikAPI Protocol = "mikrotik_api"
	ProtocolSSH         Protocol = "ssh"
	ProtocolTelnet      Protocol = "telnet"
	ProtocolTR069       Protocol = "tr069"
	ProtocolSNMP        Protocol = "snmp"
)

// DeviceStatus represents the current status of a device
type DeviceStatus string

const (
	DeviceStatusOnline  DeviceStatus = "online"
	DeviceStatusOffline DeviceStatus = "offline"
	DeviceStatusUnknown DeviceStatus = "unknown"
	DeviceStatusWarning DeviceStatus = "warning"
	DeviceStatusError   DeviceStatus = "error"
)

// Device represents a network device
type Device struct {
	ID              string       `json:"id" gorm:"primaryKey;type:uuid;default:gen_random_uuid()"`
	Name            string       `json:"name" gorm:"not null;size:255"`
	IPAddress       string       `json:"ip_address" gorm:"not null;type:inet"`
	DeviceType      DeviceType   `json:"device_type" gorm:"not null;size:50"`
	Protocol        Protocol     `json:"protocol" gorm:"not null;size:50"`
	Status          DeviceStatus `json:"status" gorm:"size:20;default:'unknown'"`
	PollingInterval int          `json:"polling_interval" gorm:"default:300"` // seconds
	CredentialsID   string       `json:"credentials_id" gorm:"type:uuid"`
	GroupID         *string      `json:"group_id,omitempty" gorm:"type:uuid"`
	Description     string       `json:"description" gorm:"type:text"`
	Tags            []string     `json:"tags" gorm:"type:text[]"`
	Metadata        JSONMap      `json:"metadata" gorm:"type:jsonb"`
	LastSeen        *time.Time   `json:"last_seen,omitempty"`
	LastError       string       `json:"last_error,omitempty" gorm:"type:text"`
	Enabled         bool         `json:"enabled" gorm:"default:true"`
	CreatedAt       time.Time    `json:"created_at"`
	UpdatedAt       time.Time    `json:"updated_at"`

	// Relationships
	Credentials *DeviceCredentials `json:"credentials,omitempty" gorm:"foreignKey:CredentialsID"`
	Group       *DeviceGroup       `json:"group,omitempty" gorm:"foreignKey:GroupID"`
}

// DeviceCredentials stores encrypted authentication credentials
type DeviceCredentials struct {
	ID                string    `json:"id" gorm:"primaryKey;type:uuid;default:gen_random_uuid()"`
	Name              string    `json:"name" gorm:"not null;size:255"`
	Username          string    `json:"username" gorm:"not null;size:255"`
	PasswordEncrypted string    `json:"-" gorm:"column:password_encrypted;type:text"` // Never expose in JSON
	SSHKeyEncrypted   string    `json:"-" gorm:"column:ssh_key_encrypted;type:text"`
	SNMPCommunity     string    `json:"-" gorm:"column:snmp_community;size:255"`
	SNMPVersion       string    `json:"snmp_version,omitempty" gorm:"size:10"`
	Description       string    `json:"description" gorm:"type:text"`
	CreatedAt         time.Time `json:"created_at"`
	UpdatedAt         time.Time `json:"updated_at"`
}

// DeviceGroup represents a logical grouping of devices
type DeviceGroup struct {
	ID          string     `json:"id" gorm:"primaryKey;type:uuid;default:gen_random_uuid()"`
	Name        string     `json:"name" gorm:"not null;size:255;uniqueIndex"`
	ParentID    *string    `json:"parent_id,omitempty" gorm:"type:uuid"`
	Description string     `json:"description" gorm:"type:text"`
	CreatedAt   time.Time  `json:"created_at"`
	UpdatedAt   time.Time  `json:"updated_at"`

	// Relationships
	Parent   *DeviceGroup   `json:"parent,omitempty" gorm:"foreignKey:ParentID"`
	Children []*DeviceGroup `json:"children,omitempty" gorm:"foreignKey:ParentID"`
	Devices  []*Device      `json:"devices,omitempty" gorm:"foreignKey:GroupID"`
}

// JSONMap is a custom type for JSONB fields
type JSONMap map[string]interface{}

// TableName specifies the table name for Device
func (Device) TableName() string {
	return "devices"
}

// TableName specifies the table name for DeviceCredentials
func (DeviceCredentials) TableName() string {
	return "device_credentials"
}

// TableName specifies the table name for DeviceGroup
func (DeviceGroup) TableName() string {
	return "device_groups"
}

// IsOnline checks if device is currently online
func (d *Device) IsOnline() bool {
	return d.Status == DeviceStatusOnline
}

// SupportsProtocol checks if device supports a specific protocol
func (d *Device) SupportsProtocol(protocol Protocol) bool {
	return d.Protocol == protocol
}

// GetPollingIntervalDuration returns polling interval as time.Duration
func (d *Device) GetPollingIntervalDuration() time.Duration {
	return time.Duration(d.PollingInterval) * time.Second
}
