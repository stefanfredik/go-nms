package apigateway_test

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
	"github.com/yourorg/nms-go/internal/config_mgt"
	"github.com/yourorg/nms-go/internal/device/handler"
	"github.com/yourorg/nms-go/internal/device/model"
	"github.com/yourorg/nms-go/internal/device/service"
)

// MockDeviceService
type MockDeviceService struct {
	GetDeviceFunc      func(ctx context.Context, id string) (*model.Device, error)
	RegisterDeviceFunc func(ctx context.Context, req *service.RegisterDeviceRequest) (*model.Device, error)
	ListDevicesFunc    func(ctx context.Context, page, pageSize int) ([]*model.Device, int64, error)
}

func (m *MockDeviceService) GetDevice(ctx context.Context, id string) (*model.Device, error) {
	if m.GetDeviceFunc != nil {
		return m.GetDeviceFunc(ctx, id)
	}
	return nil, nil
}

func (m *MockDeviceService) RegisterDevice(ctx context.Context, req *service.RegisterDeviceRequest) (*model.Device, error) {
	if m.RegisterDeviceFunc != nil {
		return m.RegisterDeviceFunc(ctx, req)
	}
	return nil, nil
}

func (m *MockDeviceService) ListDevices(ctx context.Context, page, pageSize int) ([]*model.Device, int64, error) {
	if m.ListDevicesFunc != nil {
		return m.ListDevicesFunc(ctx, page, pageSize)
	}
	return nil, 0, nil
}

// MockConfigService
type MockConfigService struct {
	ExecuteCommandFunc func(ctx context.Context, deviceID, command string) (interface{}, error)
	BackupConfigFunc   func(ctx context.Context, deviceID string) (string, error)
}

func (m *MockConfigService) ExecuteCommand(ctx context.Context, deviceID, command string) (interface{}, error) {
	if m.ExecuteCommandFunc != nil {
		return m.ExecuteCommandFunc(ctx, deviceID, command)
	}
	return "mock output", nil
}

func (m *MockConfigService) BackupConfig(ctx context.Context, deviceID string) (string, error) {
	return "", nil
}

func setupRouter(deviceService service.DeviceService, configService config_mgt.ConfigService) *gin.Engine {
	gin.SetMode(gin.TestMode)
	r := gin.New()

	deviceHandler := handler.NewDeviceHandler(deviceService)
	configHandler := config_mgt.NewConfigHandler(configService)

	v1 := r.Group("/api/v1")
	{
		devices := v1.Group("/devices")
		{
			devices.GET("", deviceHandler.ListDevices)
			devices.POST("", deviceHandler.RegisterDevice)
			devices.GET("/:id", deviceHandler.GetDevice)
		}

		configGroup := v1.Group("/config")
		{
			configGroup.POST("/execute", configHandler.ExecuteCommand)
		}
	}
	return r
}

func TestRegisterDevice(t *testing.T) {
	mockService := &MockDeviceService{
		RegisterDeviceFunc: func(ctx context.Context, req *service.RegisterDeviceRequest) (*model.Device, error) {
			if req.Name == "BadDevice" {
				return nil, errors.New("registration failed")
			}
			return &model.Device{
				ID:        "uuid-123",
				Name:      req.Name,
				IPAddress: req.IPAddress,
			}, nil
		},
	}

	router := setupRouter(mockService, nil)

	tests := []struct {
		name         string
		payload      map[string]interface{}
		expectedCode int
	}{
		{
			name: "Success",
			payload: map[string]interface{}{
				"name":       "TestRouter",
				"ip_address": "192.168.1.1",
			},
			expectedCode: 201,
		},
		{
			name: "Failure - Service Error",
			payload: map[string]interface{}{
				"name":       "BadDevice",
				"ip_address": "192.168.1.1",
			},
			expectedCode: 500,
		},
		{
			name: "Failure - Invalid JSON",
			payload: map[string]interface{}{
				"name": 123, // Invalid type
			},
			expectedCode: 400,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			jsonBody, _ := json.Marshal(tt.payload)
			req, _ := http.NewRequest("POST", "/api/v1/devices", bytes.NewBuffer(jsonBody))
			req.Header.Set("Content-Type", "application/json")
			w := httptest.NewRecorder()
			router.ServeHTTP(w, req)

			assert.Equal(t, tt.expectedCode, w.Code)
		})
	}
}

func TestGetDevice(t *testing.T) {
	mockService := &MockDeviceService{
		GetDeviceFunc: func(ctx context.Context, id string) (*model.Device, error) {
			if id == "not-found" {
				return nil, errors.New("not found")
			}
			return &model.Device{ID: id, Name: "FoundDevice"}, nil
		},
	}

	router := setupRouter(mockService, nil)

	req, _ := http.NewRequest("GET", "/api/v1/devices/123", nil)
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)
	assert.Equal(t, 200, w.Code)

	req2, _ := http.NewRequest("GET", "/api/v1/devices/not-found", nil)
	w2 := httptest.NewRecorder()
	router.ServeHTTP(w2, req2)
	assert.Equal(t, 404, w2.Code)
}

func TestExecuteCommand(t *testing.T) {
	mockConfig := &MockConfigService{
		ExecuteCommandFunc: func(ctx context.Context, deviceID, command string) (interface{}, error) {
			if command == "fail" {
				return "", errors.New("command failed")
			}
			return "command output", nil
		},
	}

	router := setupRouter(nil, mockConfig)

	// Success case
	body, _ := json.Marshal(map[string]string{
		"device_id": "good-id",
		"command":   "/system resource print",
	})
	req, _ := http.NewRequest("POST", "/api/v1/config/execute", bytes.NewBuffer(body))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)
	assert.Equal(t, 200, w.Code)
	assert.Contains(t, w.Body.String(), "command output")

	// Failure case
	bodyFail, _ := json.Marshal(map[string]string{
		"device_id": "fail-id",
		"command":   "ls",
	})
	req2, _ := http.NewRequest("POST", "/api/v1/config/execute", bytes.NewBuffer(bodyFail))
	req2.Header.Set("Content-Type", "application/json")
	w2 := httptest.NewRecorder()
	router.ServeHTTP(w2, req2)
	assert.Equal(t, 500, w2.Code)
}
