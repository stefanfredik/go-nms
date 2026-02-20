# 📘 Panduan Development — go-nms

Panduan ini menjelaskan cara menyiapkan dan menjalankan lingkungan pengembangan (development environment)
untuk proyek **go-nms** secara lokal.

---

## 🗂️ Arsitektur Development

```
┌─────────────────────────────────────────────────────────┐
│                    Komputer Lokal Kamu                  │
│                                                         │
│  ┌──────────────────────────────────────┐               │
│  │  Go Services (dijalankan langsung)   │               │
│  │                                      │               │
│  │  api-gateway :9090  ◄── air (reload) │               │
│  └──────────┬───────────────────────────┘               │
│             │ koneksi ke localhost:xxxx                  │
│  ┌──────────▼───────────────────────────┐               │
│  │  Docker (infrastruktur)              │               │
│  │                                      │               │
│  │  PostgreSQL  :5499                   │               │
│  │  Redis       :6399                   │               │
│  │  NATS        :4299                   │               │
│  │  InfluxDB    :8099                   │               │
│  └──────────────────────────────────────┘               │
└─────────────────────────────────────────────────────────┘
```

**Pendekatan:** Infrastruktur berjalan di Docker, kode Go berjalan langsung di komputermu.
Ini membuat development lebih cepat karena tidak perlu rebuild Docker image setiap ada perubahan kode.

---

## ✅ Prasyarat

Pastikan semua ini sudah terinstall:

| Tools | Cara Cek | Install |
|---|---|---|
| Go 1.21+ | `go version` | [go.dev/dl](https://go.dev/dl) |
| Docker | `docker --version` | [docs.docker.com](https://docs.docker.com/get-docker/) |
| Docker Compose | `docker compose version` | Sudah termasuk di Docker Desktop |
| `air` (hot reload) | `air -v` | `make install-air` |
| `make` | `make --version` | Sudah ada di Linux/macOS |

---

## 🚀 Setup Pertama Kali

Lakukan ini **sekali saja** saat pertama kali setup:

```bash
# 1. Masuk ke folder proyek
cd ~/dev/go-nms

# 2. Download semua dependensi Go
go mod download

# 3. Install air (hot-reload tool)
make install-air
```

---

## 🔄 Workflow Sehari-hari

### Langkah 1 — Jalankan Infrastruktur

Buka terminal pertama:

```bash
make dev-infra
```

Perintah ini menjalankan PostgreSQL, Redis, NATS, dan InfluxDB di Docker.
Kamu hanya perlu menjalankan ini **sekali** di awal hari, atau setelah komputer restart.

Output yang diharapkan:
```
🚀 Menjalankan infrastruktur development...
✅ Infrastruktur siap:
   PostgreSQL : localhost:5499
   Redis      : localhost:6399
   NATS       : localhost:4299
   InfluxDB   : localhost:8099
```

### Langkah 2 — Jalankan API Gateway (dengan Hot Reload)

Buka terminal kedua:

```bash
make dev
```

Perintah ini menjalankan `api-gateway` dengan **hot reload** menggunakan `air`.
Artinya: setiap kali kamu menyimpan file `.go`, aplikasi otomatis restart tanpa perlu perintah apapun.

Output yang diharapkan:
```
🔥 Menjalankan api-gateway dengan hot-reload...

  __    _   ___
 / /\  | | | |_)
/_/--\ |_| |_| \  v1.x.x

watching ...
building...
running...
[GIN-debug] Listening and serving HTTP on :9090
```

### Langkah 3 — Mulai Coding!

API Gateway berjalan di: **http://localhost:9090**

Coba akses:
```bash
curl http://localhost:9090/api/v1/devices
```

---

## 🛠️ Perintah Berguna

```bash
# Jalankan infrastruktur
make dev-infra

# Jalankan api-gateway dengan hot reload
make dev

# Cek status container
make dev-status

# Lihat log infrastruktur
make dev-logs

# Hentikan infrastruktur (akhir hari)
make dev-down

# Build binary tanpa menjalankan
make build-api-gateway

# Jalankan semua test
make test

# Format kode Go
make fmt
```

---

## 📁 Struktur File Penting

```
go-nms/
├── cmd/
│   ├── api-gateway/      # Entry point API Gateway
│   │   └── main.go       # ← Titik mulai program
│   ├── worker/           # Background worker
│   └── collector/        # SNMP collector
│
├── internal/
│   ├── common/
│   │   ├── config/       # Konfigurasi aplikasi
│   │   └── database/     # Koneksi database
│   ├── device/           # Modul device (CRUD)
│   ├── features/
│   │   ├── monitoring/   # Monitoring & metrics
│   │   └── olt/          # OLT (Optical Line Terminal)
│   └── api-gateway/      # Router & middleware
│
├── config.yaml           # Konfigurasi default (untuk go run lokal)
├── .env.dev              # Environment variables untuk development
├── .air.toml             # Konfigurasi hot-reload
├── docker-compose.yml    # Docker: semua service (production-like)
├── docker-compose.dev.yml # Docker: hanya infrastruktur (development)
└── Makefile              # Kumpulan perintah shortcut
```

---

## ⚙️ Konfigurasi

Aplikasi membaca konfigurasi dari dua sumber (urutan prioritas: env var > config.yaml):

### File `.env.dev` (untuk development)

```env
DATABASE_HOST=localhost
DATABASE_PORT=5499
DATABASE_USER=nms
DATABASE_PASSWORD=nms_password
DATABASE_DBNAME=nms_db
...
```

File ini dibaca otomatis oleh `air` saat `make dev` dijalankan.

### File `config.yaml` (nilai default)

```yaml
server:
  port: 9090
database:
  host: "localhost"
  port: 5499
  ...
```

---

## 🐛 Debugging

### Cek apakah infrastruktur berjalan

```bash
make dev-status
```

### Lihat log database / redis

```bash
make dev-logs
# atau spesifik satu service:
docker compose -f docker-compose.dev.yml logs postgres -f
docker compose -f docker-compose.dev.yml logs redis -f
```

### Konek ke database PostgreSQL

```bash
docker compose -f docker-compose.dev.yml exec postgres \
  psql -U nms -d nms_db
```

### Konek ke Redis

```bash
docker compose -f docker-compose.dev.yml exec redis redis-cli
```

### Reset database (hapus semua data)

```bash
make dev-down
docker volume rm go-nms_postgres_dev_data
make dev-infra
```

---

## 🏗️ Menjalankan Service Lain

Selain `api-gateway`, ada service lain yang bisa dijalankan:

```bash
# Terminal terpisah untuk setiap service:
go run ./cmd/worker/main.go
go run ./cmd/collector/main.go
go run ./cmd/alert/main.go
```

Atau gunakan script bawaan (berjalan di background):

```bash
./start_services.sh
```

> **Catatan:** Service-service tersebut membaca config dari `config.yaml` secara default.
> Pastikan infrastruktur sudah berjalan (`make dev-infra`) sebelum menjalankan service manapun.

---

## ❓ Troubleshooting

### Error: `air: command not found`
```bash
make install-air
# Atau tambahkan ke PATH:
export PATH=$PATH:$(go env GOPATH)/bin
```

### Error: `failed to connect to database`
- Pastikan infrastruktur berjalan: `make dev-status`
- Cek apakah port 5499 dipakai aplikasi lain: `lsof -i :5499`
- Restart infrastruktur: `make dev-down && make dev-infra`

### Error: `port already in use: :9090`
```bash
# Cari proses yang pakai port 9090
lsof -ti :9090 | xargs kill -9
```

### Hot reload tidak bekerja
- Pastikan file disimpan (Ctrl+S)
- Cek log `air` di terminal
- Coba restart: `Ctrl+C` lalu `make dev`
