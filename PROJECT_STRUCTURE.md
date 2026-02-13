# NMS-Go Project Structure

## Overview
Struktur project ini mengikuti best practices Go project layout (https://github.com/golang-standards/project-layout) dengan pendekatan microservices architecture.

## Directory Structure

```
nms-go/
├── cmd/                          # Main applications
├── internal/                     # Private application code
├── pkg/                          # Public reusable libraries
├── api/                          # API contracts & definitions
├── configs/                      # Configuration files
├── scripts/                      # Build, deployment, utility scripts
├── deployments/                  # Deployment configurations
├── docs/                         # Documentation
├── test/                         # Additional test files
├── go.mod                        # Go module definition
├── go.sum                        # Go dependencies checksum
├── Makefile                      # Build automation
└── README.md                     # Project README
```

---

## 1. `/cmd` - Main Applications

Setiap subdirectory adalah aplikasi yang dapat dijalankan. File `main.go` berada di setiap direktori.

### `/cmd/api-gateway`
**Purpose**: HTTP API Gateway - entry point untuk semua REST API requests
- Authentication & authorization
- Rate limiting
- Request routing
- Load balancing

**Example structure**:
```go
// cmd/api-gateway/main.go
package main

func main() {
    // Initialize gateway
    // Setup routes
    // Start server
}
```

### `/cmd/device-service`
**Purpose**: Device management service
- Device CRUD operations
- Device inventory management
- Device grouping

### `/cmd/collector-service`
**Purpose**: Orchestrate data collection
- Manage polling schedules
- Dispatch polling tasks to message queue
- Monitor worker health

### `/cmd/worker`
**Purpose**: Execute actual device polling
- Consume tasks from message queue
- Connect to devices via protocols
- Store collected data
- Multiple workers can run in parallel

### `/cmd/config-service`
**Purpose**: Configuration management service
- Execute remote configurations
- Backup management
- Version control
- Rollback functionality

### `/cmd/alert-service`
**Purpose**: Alert evaluation service
- Evaluate alert rules
- Trigger notifications
- Alert state management

### `/cmd/notification-service`
**Purpose**: Send notifications
- Email notifications
- Telegram bot
- Webhook delivery
- SMS gateway integration

### `/cmd/analytics-service`
**Purpose**: Data processing & analytics
- Aggregation
- Report generation
- Trend analysis

### `/cmd/migrate`
**Purpose**: Database migration tool
- Run database migrations
- Seed initial data

---

## 2. `/internal` - Private Application Code

Code yang tidak boleh di-import oleh aplikasi lain. Go compiler akan enforce hal ini.

### `/internal/api-gateway/`
```
api-gateway/
├── handler/          # HTTP request handlers
├── middleware/       # Auth, logging, CORS middlewares
└── router/          # Route definitions
```

### `/internal/device/`
```
device/
├── handler/         # HTTP handlers for device endpoints
├── service/         # Business logic layer
├── repository/      # Data access layer
└── model/          # Domain models & entities
```

**Example**:
```go
// internal/device/model/device.go
type Device struct {
    ID             string
    Name           string
    IPAddress      string
    DeviceType     DeviceType
    Protocol       Protocol
    Status         DeviceStatus
    PollingInterval int
}

// internal/device/repository/device_repo.go
type DeviceRepository interface {
    Create(ctx context.Context, device *Device) error
    GetByID(ctx context.Context, id string) (*Device, error)
    List(ctx context.Context, filter Filter) ([]*Device, error)
    Update(ctx context.Context, device *Device) error
    Delete(ctx context.Context, id string) error
}

// internal/device/service/device_service.go
type DeviceService interface {
    RegisterDevice(ctx context.Context, req RegisterRequest) error
    ValidateConnection(ctx context.Context, device *Device) error
    GetDeviceMetrics(ctx context.Context, deviceID string) (*Metrics, error)
}
```

### `/internal/worker/`
```
worker/
├── handler/         # Task handler
├── parser/          # Parse device responses
└── protocols/       # Protocol implementations
    ├── mikrotik/    # Mikrotik RouterOS API
    ├── ssh/         # SSH client
    ├── telnet/      # Telnet client
    ├── tr069/       # TR-069/CWMP
    └── snmp/        # SNMP v2c/v3
```

**Example protocol interface**:
```go
// internal/worker/protocols/protocol.go
type DeviceProtocol interface {
    Connect(ctx context.Context, device *Device) error
    Disconnect() error
    ExecuteCommand(ctx context.Context, cmd string) (string, error)
    GetSystemMetrics(ctx context.Context) (*SystemMetrics, error)
    GetInterfaceMetrics(ctx context.Context) ([]*InterfaceMetrics, error)
}

// internal/worker/protocols/mikrotik/mikrotik.go
type MikrotikClient struct {
    client *routeros.Client
}

func (m *MikrotikClient) Connect(ctx context.Context, device *Device) error {
    // Implementation
}
```

### `/internal/collector/`
```
collector/
├── scheduler/       # Job scheduling logic
├── service/         # Collection orchestration
└── repository/      # Polling config storage
```

### `/internal/config/`
```
config/
├── handler/         # Config API handlers
├── service/         # Config management logic
├── repository/      # Config storage
└── versioning/      # Version control logic
```

### `/internal/alert/`
```
alert/
├── engine/          # Rule evaluation engine
├── handler/         # Alert API handlers
├── service/         # Alert management
└── repository/      # Alert rules storage
```

### `/internal/notification/`
```
notification/
├── sender/          # Notification senders
│   ├── email.go
│   ├── telegram.go
│   ├── webhook.go
│   └── sms.go
├── handler/         # Notification API handlers
└── service/         # Notification logic
```

### `/internal/common/`
Shared internal code across services:
```
common/
├── auth/            # JWT, session management
├── validator/       # Input validation
├── logger/          # Structured logging wrapper
├── errors/          # Custom error types
└── config/          # Configuration loader
```

---

## 3. `/pkg` - Public Libraries

Code yang boleh di-import oleh external projects.

### `/pkg/database/`
Database clients dan helpers:
```go
// pkg/database/postgres/postgres.go
func NewPostgresDB(config Config) (*gorm.DB, error) {
    // Initialize PostgreSQL connection
}

// pkg/database/influxdb/influxdb.go
func NewInfluxDBClient(config Config) (influxdb2.Client, error) {
    // Initialize InfluxDB connection
}

// pkg/database/redis/redis.go
func NewRedisClient(config Config) (*redis.Client, error) {
    // Initialize Redis connection
}
```

### `/pkg/messagequeue/`
Message queue clients:
```go
// pkg/messagequeue/nats/nats.go
func NewNATSClient(url string) (*nats.Conn, error)

// pkg/messagequeue/rabbitmq/rabbitmq.go
func NewRabbitMQClient(config Config) (*amqp.Connection, error)
```

### `/pkg/crypto/`
Encryption utilities:
```go
// pkg/crypto/aes.go
func Encrypt(plaintext string, key []byte) (string, error)
func Decrypt(ciphertext string, key []byte) (string, error)
```

### `/pkg/utils/`
Common utilities:
```go
// pkg/utils/network.go
func ValidateIP(ip string) bool
func Ping(ip string) (time.Duration, error)

// pkg/utils/retry.go
func RetryWithBackoff(fn func() error, maxRetries int) error
```

### `/pkg/metrics/`
Prometheus metrics helpers:
```go
// pkg/metrics/metrics.go
func IncrementCounter(name string, labels map[string]string)
func RecordHistogram(name string, value float64)
```

### `/pkg/logging/`
Structured logging utilities:
```go
// pkg/logging/logger.go
func NewLogger(config Config) (*zap.Logger, error)
func WithContext(ctx context.Context, logger *zap.Logger) *zap.Logger
```

---

## 4. `/api` - API Definitions

### `/api/proto/`
Protocol Buffer definitions untuk gRPC (jika digunakan):
```protobuf
// api/proto/device/v1/device.proto
syntax = "proto3";

package device.v1;

service DeviceService {
  rpc CreateDevice(CreateDeviceRequest) returns (CreateDeviceResponse);
  rpc GetDevice(GetDeviceRequest) returns (GetDeviceResponse);
}
```

### `/api/openapi/`
OpenAPI/Swagger specifications:
```yaml
# api/openapi/device-service.yaml
openapi: 3.0.0
info:
  title: Device Service API
  version: 1.0.0
paths:
  /devices:
    get:
      summary: List all devices
    post:
      summary: Register new device
```

### `/api/graphql/`
GraphQL schemas (jika digunakan):
```graphql
# api/graphql/schema.graphql
type Device {
  id: ID!
  name: String!
  ipAddress: String!
  status: DeviceStatus!
}
```

---

## 5. `/configs` - Configuration Files

### `/configs/env/`
Environment-specific configs:
```
env/
├── development.yaml
├── staging.yaml
└── production.yaml
```

**Example config**:
```yaml
# configs/env/development.yaml
server:
  port: 8080
  host: 0.0.0.0

database:
  postgres:
    host: localhost
    port: 5432
    database: nms_dev
  influxdb:
    url: http://localhost:8086
    token: dev-token
  redis:
    addr: localhost:6379

messagequeue:
  nats:
    url: nats://localhost:4222

logging:
  level: debug
  format: json
```

### `/configs/database/`
Database specific configs:
```sql
-- configs/database/init.sql
CREATE DATABASE nms_production;
CREATE USER nms_user WITH ENCRYPTED PASSWORD 'secure_password';
GRANT ALL PRIVILEGES ON DATABASE nms_production TO nms_user;
```

---

## 6. `/scripts` - Scripts

### `/scripts/migrations/`
Database migration files:
```sql
-- scripts/migrations/000001_create_devices_table.up.sql
CREATE TABLE devices (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    name VARCHAR(255) NOT NULL,
    ip_address INET NOT NULL,
    device_type VARCHAR(50) NOT NULL,
    protocol VARCHAR(50) NOT NULL,
    status VARCHAR(20) DEFAULT 'unknown',
    polling_interval INTEGER DEFAULT 300,
    created_at TIMESTAMP DEFAULT NOW(),
    updated_at TIMESTAMP DEFAULT NOW()
);

-- scripts/migrations/000001_create_devices_table.down.sql
DROP TABLE devices;
```

### `/scripts/seed/`
Seed data scripts:
```go
// scripts/seed/devices.go
package main

func seedDevices() {
    // Insert sample devices
}
```

**Utility scripts**:
```bash
# scripts/build.sh
#!/bin/bash
# Build all services

# scripts/run-dev.sh
#!/bin/bash
# Run services in development mode

# scripts/generate-mocks.sh
#!/bin/bash
# Generate test mocks
```

---

## 7. `/deployments` - Deployment Configs

### `/deployments/docker/`
Individual service Dockerfiles:
```dockerfile
# deployments/docker/Dockerfile.api-gateway
FROM golang:1.21-alpine AS builder
WORKDIR /app
COPY . .
RUN go build -o api-gateway ./cmd/api-gateway

FROM alpine:latest
COPY --from=builder /app/api-gateway /api-gateway
EXPOSE 8080
CMD ["/api-gateway"]
```

### `/deployments/docker-compose/`
Docker Compose files:
```yaml
# deployments/docker-compose/docker-compose.dev.yml
version: '3.8'

services:
  postgres:
    image: postgres:14
    environment:
      POSTGRES_DB: nms_dev
      POSTGRES_USER: nms_user
      POSTGRES_PASSWORD: password
    ports:
      - "5432:5432"

  influxdb:
    image: influxdb:2.7
    ports:
      - "8086:8086"

  redis:
    image: redis:7-alpine
    ports:
      - "6379:6379"

  nats:
    image: nats:latest
    ports:
      - "4222:4222"

  api-gateway:
    build:
      context: ../..
      dockerfile: deployments/docker/Dockerfile.api-gateway
    ports:
      - "8080:8080"
    depends_on:
      - postgres
      - redis
```

### `/deployments/kubernetes/`
Kubernetes manifests:
```yaml
# deployments/kubernetes/api-gateway-deployment.yaml
apiVersion: apps/v1
kind: Deployment
metadata:
  name: api-gateway
spec:
  replicas: 3
  selector:
    matchLabels:
      app: api-gateway
  template:
    metadata:
      labels:
        app: api-gateway
    spec:
      containers:
      - name: api-gateway
        image: nms/api-gateway:latest
        ports:
        - containerPort: 8080
```

---

## 8. `/test` - Test Files

### `/test/integration/`
Integration tests:
```go
// test/integration/device_test.go
func TestDeviceEndToEnd(t *testing.T) {
    // Test full device registration flow
}
```

### `/test/e2e/`
End-to-end tests:
```go
// test/e2e/polling_flow_test.go
func TestCompletePollingFlow(t *testing.T) {
    // Test from device registration to data collection
}
```

### `/test/fixtures/`
Test data fixtures:
```json
// test/fixtures/devices.json
[
  {
    "name": "Router-01",
    "ip_address": "192.168.1.1",
    "device_type": "router",
    "protocol": "mikrotik_api"
  }
]
```

---

## 9. Root Level Files

### `go.mod`
```go
module github.com/yourorg/nms-go

go 1.21

require (
    github.com/gin-gonic/gin v1.9.1
    gorm.io/gorm v1.25.5
    // ... dependencies
)
```

### `Makefile`
```makefile
.PHONY: build test run migrate

build:
	@echo "Building services..."
	go build -o bin/api-gateway ./cmd/api-gateway
	go build -o bin/worker ./cmd/worker

test:
	go test -v ./...

run-dev:
	docker-compose -f deployments/docker-compose/docker-compose.dev.yml up

migrate-up:
	migrate -path scripts/migrations -database "postgres://localhost:5432/nms?sslmode=disable" up

lint:
	golangci-lint run

generate-mocks:
	mockgen -source=internal/device/repository/device_repo.go -destination=internal/device/repository/mock/device_repo_mock.go
```

### `README.md`
```markdown
# NMS-Go - Network Management System

## Prerequisites
- Go 1.21+
- Docker & Docker Compose
- PostgreSQL 14+
- InfluxDB 2.x
- Redis 7+

## Quick Start
1. Clone repository
2. Run `make run-dev`
3. Access API at http://localhost:8080

## Development
- `make build` - Build all services
- `make test` - Run tests
- `make migrate-up` - Run migrations
```

---

## Design Principles

### 1. **Separation of Concerns**
- Each service has single responsibility
- Clear boundaries between layers (handler → service → repository)

### 2. **Dependency Injection**
- Dependencies injected through constructors
- Easy to mock for testing

### 3. **Interface-Driven Design**
- Define interfaces for repositories and services
- Enables easy swapping of implementations

### 4. **Configuration Management**
- Environment-specific configs
- Secrets via environment variables
- Use Viper for config loading

### 5. **Error Handling**
- Custom error types in `/internal/common/errors`
- Consistent error responses
- Proper logging of errors

### 6. **Testing Strategy**
- Unit tests alongside code (`*_test.go`)
- Integration tests in `/test/integration`
- Use table-driven tests
- Mock external dependencies

---

## Example Service Implementation

### Complete Device Service Example

```go
// internal/device/model/device.go
package model

type Device struct {
    ID              string    `json:"id" gorm:"primaryKey"`
    Name            string    `json:"name" gorm:"not null"`
    IPAddress       string    `json:"ip_address" gorm:"not null"`
    DeviceType      string    `json:"device_type"`
    Protocol        string    `json:"protocol"`
    Status          string    `json:"status"`
    PollingInterval int       `json:"polling_interval"`
    CreatedAt       time.Time `json:"created_at"`
    UpdatedAt       time.Time `json:"updated_at"`
}

// internal/device/repository/device_repo.go
package repository

type DeviceRepository interface {
    Create(ctx context.Context, device *model.Device) error
    GetByID(ctx context.Context, id string) (*model.Device, error)
    List(ctx context.Context) ([]*model.Device, error)
    Update(ctx context.Context, device *model.Device) error
    Delete(ctx context.Context, id string) error
}

type deviceRepository struct {
    db *gorm.DB
}

func NewDeviceRepository(db *gorm.DB) DeviceRepository {
    return &deviceRepository{db: db}
}

func (r *deviceRepository) Create(ctx context.Context, device *model.Device) error {
    return r.db.WithContext(ctx).Create(device).Error
}

// internal/device/service/device_service.go
package service

type DeviceService interface {
    RegisterDevice(ctx context.Context, req RegisterDeviceRequest) (*model.Device, error)
    GetDevice(ctx context.Context, id string) (*model.Device, error)
}

type deviceService struct {
    repo repository.DeviceRepository
}

func NewDeviceService(repo repository.DeviceRepository) DeviceService {
    return &deviceService{repo: repo}
}

func (s *deviceService) RegisterDevice(ctx context.Context, req RegisterDeviceRequest) (*model.Device, error) {
    device := &model.Device{
        ID:              uuid.New().String(),
        Name:            req.Name,
        IPAddress:       req.IPAddress,
        DeviceType:      req.DeviceType,
        Protocol:        req.Protocol,
        Status:          "unknown",
        PollingInterval: req.PollingInterval,
    }
    
    if err := s.repo.Create(ctx, device); err != nil {
        return nil, err
    }
    
    return device, nil
}

// internal/device/handler/device_handler.go
package handler

type DeviceHandler struct {
    service service.DeviceService
}

func NewDeviceHandler(service service.DeviceService) *DeviceHandler {
    return &DeviceHandler{service: service}
}

func (h *DeviceHandler) RegisterDevice(c *gin.Context) {
    var req RegisterDeviceRequest
    if err := c.ShouldBindJSON(&req); err != nil {
        c.JSON(400, gin.H{"error": err.Error()})
        return
    }
    
    device, err := h.service.RegisterDevice(c.Request.Context(), req)
    if err != nil {
        c.JSON(500, gin.H{"error": err.Error()})
        return
    }
    
    c.JSON(201, device)
}

// cmd/device-service/main.go
package main

func main() {
    // Load config
    cfg := config.Load()
    
    // Initialize database
    db := postgres.NewPostgresDB(cfg.Database.Postgres)
    
    // Initialize repository
    deviceRepo := repository.NewDeviceRepository(db)
    
    // Initialize service
    deviceService := service.NewDeviceService(deviceRepo)
    
    // Initialize handler
    deviceHandler := handler.NewDeviceHandler(deviceService)
    
    // Setup router
    r := gin.Default()
    r.POST("/devices", deviceHandler.RegisterDevice)
    
    // Start server
    r.Run(":8080")
}
```

---

## Migration Strategy

### Phase 1: Setup & Foundation
1. Setup project structure
2. Implement common packages (`pkg/`)
3. Setup databases and message queue
4. Implement authentication

### Phase 2: Core Services
1. Device Service
2. Basic polling (Collector + Worker)
3. Simple metrics storage

### Phase 3: Protocol Support
1. Implement all protocol adapters
2. Advanced polling strategies
3. Error handling & retry logic

### Phase 4: Advanced Features
1. Configuration management
2. Alert engine
3. Notification system
4. Analytics & reporting

---

## Best Practices

1. **Use contexts**: Pass `context.Context` for cancellation and timeouts
2. **Log consistently**: Use structured logging (zap)
3. **Handle errors properly**: Don't ignore errors, wrap with context
4. **Write tests**: Aim for >80% coverage
5. **Document code**: Use godoc comments
6. **Use interfaces**: For better testability and flexibility
7. **Avoid global state**: Use dependency injection
8. **Keep functions small**: Single responsibility principle
9. **Use meaningful names**: Clear, descriptive variable/function names
10. **Version your APIs**: Use `/api/v1/` prefixes

---

## Tools & Commands

```bash
# Install dependencies
go mod download

# Run tests
go test ./...

# Run tests with coverage
go test -cover ./...

# Generate coverage report
go test -coverprofile=coverage.out ./...
go tool cover -html=coverage.out

# Format code
go fmt ./...

# Lint code
golangci-lint run

# Generate mocks
go install github.com/golang/mock/mockgen@latest
mockgen -source=internal/device/repository/device_repo.go -destination=internal/device/repository/mock/device_repo_mock.go

# Run database migrations
migrate -path scripts/migrations -database "postgres://localhost:5432/nms?sslmode=disable" up
```

---

## Resources

- [Go Project Layout](https://github.com/golang-standards/project-layout)
- [Effective Go](https://golang.org/doc/effective_go)
- [Go Code Review Comments](https://github.com/golang/go/wiki/CodeReviewComments)
- [GORM Documentation](https://gorm.io/docs/)
- [Gin Documentation](https://gin-gonic.com/docs/)
