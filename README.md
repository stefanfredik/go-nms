# NMS-Go - Network Management System

![Go Version](https://img.shields.io/badge/Go-1.21%2B-blue)
![License](https://img.shields.io/badge/license-MIT-green)

Network Management System yang scalable untuk mengelola ribuan perangkat jaringan menggunakan Go.

## 🌟 Features

- **Multi-Protocol Support**: Mikrotik API, SSH, Telnet, TR-069, SNMP
- **Device Management**: Register, group, dan monitor berbagai jenis perangkat
- **Real-time Monitoring**: Dashboard dengan live updates
- **Configuration Management**: Remote config, backup, dan versioning
- **Alert System**: Threshold-based alerting dengan multiple notification channels
- **Scalable Architecture**: Microservices dengan horizontal scaling capability
- **Time-Series Data**: Efficient metrics storage dengan InfluxDB

## 🏗️ Architecture

```
┌─────────────┐     ┌──────────────┐     ┌─────────────┐
│   Web UI    │────▶│ API Gateway  │────▶│   Device    │
└─────────────┘     └──────────────┘     │   Service   │
                            │             └─────────────┘
                            │
                    ┌───────┴────────┐
                    │                │
              ┌─────▼──────┐   ┌────▼─────┐
              │ Collector  │   │  Config  │
              │  Service   │   │ Service  │
              └─────┬──────┘   └──────────┘
                    │
              ┌─────▼──────┐
              │   NATS/    │
              │  RabbitMQ  │
              └─────┬──────┘
                    │
         ┌──────────┼──────────┐
         │          │          │
    ┌────▼───┐ ┌───▼────┐ ┌───▼────┐
    │Worker 1│ │Worker 2│ │Worker N│
    └────┬───┘ └───┬────┘ └───┬────┘
         │         │          │
         └─────────┼──────────┘
                   │
         ┌─────────▼──────────┐
         │     InfluxDB       │
         └────────────────────┘
```

## 📋 Prerequisites

- Go 1.21 atau lebih baru
- Docker & Docker Compose
- PostgreSQL 14+
- InfluxDB 2.x
- Redis 7+
- NATS/RabbitMQ

## 🚀 Quick Start

### 1. Clone Repository

```bash
git clone https://github.com/yourorg/nms-go.git
cd nms-go
```

### 2. Setup Environment

```bash
cp configs/env/development.yaml.example configs/env/development.yaml
# Edit configs/env/development.yaml sesuai kebutuhan
```

### 3. Start Dependencies

```bash
docker-compose -f deployments/docker-compose/docker-compose.dev.yml up -d
```

### 4. Run Migrations

```bash
make migrate-up
```

### 5. Start Services

```bash
# Terminal 1: API Gateway
go run cmd/api-gateway/main.go

# Terminal 2: Device Service
go run cmd/device-service/main.go

# Terminal 3: Collector Service
go run cmd/collector-service/main.go

# Terminal 4: Workers
go run cmd/worker/main.go
```

### 6. Access Dashboard

```
http://localhost:8080
```

## 🛠️ Development

### Build All Services

```bash
make build
```

### Run Tests

```bash
make test
```

### Run with Coverage

```bash
make test-coverage
```

### Lint Code

```bash
make lint
```

### Generate Mocks

```bash
make generate-mocks
```

## 📁 Project Structure

Lihat [PROJECT_STRUCTURE.md](PROJECT_STRUCTURE.md) untuk detail lengkap struktur project.

```
nms-go/
├── cmd/                    # Main applications
├── internal/               # Private application code
├── pkg/                    # Public reusable libraries
├── api/                    # API definitions
├── configs/                # Configuration files
├── scripts/                # Build & deployment scripts
├── deployments/            # Docker & K8s configs
├── docs/                   # Documentation
└── test/                   # Test files
```

## 🔧 Configuration

Configuration files berada di `configs/env/`:

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
    user: nms_user
    password: password
    
  influxdb:
    url: http://localhost:8086
    token: your-token
    org: nms
    bucket: metrics
    
  redis:
    addr: localhost:6379
    password: ""
    db: 0

messagequeue:
  nats:
    url: nats://localhost:4222

logging:
  level: debug
  format: json
```

## 📊 Database Schema

### PostgreSQL Tables

- **devices**: Device inventory
- **device_credentials**: Encrypted credentials
- **device_groups**: Device grouping
- **alert_rules**: Alert configurations
- **config_backups**: Configuration backups
- **users**: User management

### InfluxDB Measurements

- **device_system**: CPU, memory, disk metrics
- **interface_traffic**: Interface statistics
- **wireless_metrics**: Wireless-specific metrics

## 🔌 Supported Protocols

### Mikrotik RouterOS API
```go
// internal/worker/protocols/mikrotik/mikrotik.go
client, err := mikrotik.Connect(device.IPAddress, device.Username, device.Password)
metrics, err := client.GetSystemMetrics(ctx)
```

### SSH
```go
// internal/worker/protocols/ssh/ssh.go
client, err := ssh.Connect(device.IPAddress, device.Username, device.Password)
output, err := client.ExecuteCommand(ctx, "show version")
```

### Telnet
```go
// internal/worker/protocols/telnet/telnet.go
client, err := telnet.Connect(device.IPAddress)
err := client.Login(device.Username, device.Password)
```

### TR-069
```go
// internal/worker/protocols/tr069/tr069.go
client := tr069.NewClient(acsURL)
params, err := client.GetParameterValues(deviceID, []string{"Device.DeviceInfo."})
```

### SNMP
```go
// internal/worker/protocols/snmp/snmp.go
client, err := snmp.Connect(device.IPAddress, community)
result, err := client.Get([]string{"1.3.6.1.2.1.1.1.0"})
```

## 📡 API Endpoints

### Device Management

```bash
# Register device
POST /api/v1/devices
{
  "name": "Router-01",
  "ip_address": "192.168.1.1",
  "device_type": "router",
  "protocol": "mikrotik_api",
  "username": "admin",
  "password": "password"
}

# List devices
GET /api/v1/devices

# Get device details
GET /api/v1/devices/:id

# Update device
PUT /api/v1/devices/:id

# Delete device
DELETE /api/v1/devices/:id
```

### Metrics

```bash
# Get device metrics
GET /api/v1/metrics/:deviceId?start=2024-01-01&end=2024-01-31

# Query metrics
GET /api/v1/metrics/query?metric=cpu_usage&deviceId=123&interval=5m
```

### Configuration

```bash
# Execute command
POST /api/v1/config/execute
{
  "device_ids": ["uuid1", "uuid2"],
  "commands": ["show version", "show ip interface brief"]
}

# Backup config
POST /api/v1/config/backup/:deviceId

# List backups
GET /api/v1/config/backups/:deviceId

# Restore config
POST /api/v1/config/restore
{
  "device_id": "uuid",
  "backup_id": "backup-uuid"
}
```

## 🚨 Alert Rules

Example alert rule definition:

```json
{
  "name": "High CPU Usage",
  "metric": "cpu_usage",
  "condition": "value > 80",
  "duration": "5m",
  "severity": "warning",
  "notification_channels": ["email", "telegram"]
}
```

## 📈 Monitoring

System mengexpose Prometheus metrics di `/metrics`:

```
# Collector metrics
nms_devices_total{type="router"}
nms_polling_duration_seconds
nms_polling_errors_total

# Worker metrics
nms_tasks_processed_total
nms_task_duration_seconds
nms_worker_active_connections
```

## 🧪 Testing

### Unit Tests
```bash
go test ./internal/device/service/...
```

### Integration Tests
```bash
go test ./test/integration/...
```

### E2E Tests
```bash
go test ./test/e2e/...
```

## 🐳 Docker Deployment

### Build Images

```bash
make docker-build
```

### Run with Docker Compose

```bash
docker-compose -f deployments/docker-compose/docker-compose.prod.yml up -d
```

## ☸️ Kubernetes Deployment

```bash
kubectl apply -f deployments/kubernetes/
```

## 🤝 Contributing

1. Fork the repository
2. Create feature branch (`git checkout -b feature/amazing-feature`)
3. Commit changes (`git commit -m 'Add amazing feature'`)
4. Push to branch (`git push origin feature/amazing-feature`)
5. Open Pull Request

## 📝 License

This project is licensed under the MIT License - see [LICENSE](LICENSE) file for details.

## 📧 Contact

Project Link: [https://github.com/yourorg/nms-go](https://github.com/yourorg/nms-go)

## 🙏 Acknowledgments

- [Go Project Layout](https://github.com/golang-standards/project-layout)
- [Gin Web Framework](https://gin-gonic.com/)
- [GORM](https://gorm.io/)
- [InfluxDB](https://www.influxdata.com/)
