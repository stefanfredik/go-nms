# go-nms API Documentation

**Base URL:** `http://localhost:8080/api/v1`

> **Architecture Note:** go-nms does **not** maintain its own device database.
> openaccess is the single source of truth for device inventory.
> All OLT endpoints accept SNMP connection details in the request body,
> allowing go-nms to connect directly to any OLT on demand.

---

## Table of Contents

- [Health Check](#health-check)
- [OLT Resources (ZTE C320 SNMP)](#olt-resources-zte-c320-snmp)
  - [POST /olt/system](#post-oltsystem)
  - [POST /olt/pon-ports](#post-oltpon-ports)
  - [POST /olt/onts](#post-oltonts)
- [Realtime Execution (Mikrotik)](#realtime-execution-mikrotik)
  - [POST /realtime/execute](#post-realtimeexecute)
  - [POST /realtime/stats](#post-realtimestats)
- [Device Registry](#device-registry)
  - [GET /devices](#get-devices)
  - [POST /devices](#post-devices)
  - [GET /devices/:id](#get-devicesid)
- [Config Management](#config-management)
  - [POST /config/execute](#post-configexecute)
- [Inventory Sync](#inventory-sync)
  - [POST /inventory/sync](#post-inventorysync)

---

## Health Check

### GET /health

Returns the server status.

**Response `200 OK`:**
```json
{ "status": "ok" }
```

---

## OLT Resources (ZTE C320 SNMP)

All OLT endpoints use **POST** with a JSON body containing an `SNMPTarget` object.
This allows openaccess to pass OLT connection details directly without any prior
device registration in go-nms.

### SNMPTarget Object

| Field       | Type   | Required | Default  | Description                        |
|-------------|--------|----------|----------|------------------------------------|
| `ip`        | string | ✅       | —        | Management IP address of the OLT   |
| `community` | string | ❌       | `public` | SNMP v2c community string          |
| `version`   | string | ❌       | `2c`     | SNMP version (`2c` only currently) |
| `port`      | uint16 | ❌       | `161`    | SNMP UDP port                      |

---

### POST /olt/system

Fetches system-level metrics from a ZTE C320 OLT via SNMP.

**Request Body:**
```json
{
  "target": {
    "ip": "192.168.1.100",
    "community": "public",
    "version": "2c"
  }
}
```

**Response `200 OK`:**
```json
{
  "ip_address": "192.168.1.100",
  "timestamp": "2026-02-18T02:50:00Z",
  "sys_descr": "ZTE Corporation ZXA10 C320",
  "sys_name": "OLT-CORE-A",
  "uptime_seconds": 1234567,
  "cpu_usage_percent": 23.5,
  "memory_total_kb": 524288,
  "memory_used_kb": 262144,
  "memory_usage_percent": 50.0,
  "temperature_celsius": 42.0
}
```

**Error `400 Bad Request`** — missing or invalid `target.ip`:
```json
{ "error": "invalid request body: Key: 'GetSystemMetricsRequest.Target.IP' Error:Field validation for 'IP' failed on the 'required' tag" }
```

**Error `500 Internal Server Error`** — SNMP connection failed:
```json
{ "error": "failed to connect to OLT 192.168.1.100 via SNMP: ..." }
```

---

### POST /olt/pon-ports

Fetches metrics for all PON ports on a ZTE C320 OLT.

**Request Body:**
```json
{
  "target": {
    "ip": "192.168.1.100",
    "community": "public"
  }
}
```

**Response `200 OK`:**
```json
{
  "ip_address": "192.168.1.100",
  "count": 8,
  "pon_ports": [
    {
      "ip_address": "192.168.1.100",
      "timestamp": "2026-02-18T02:50:00Z",
      "port_index": 1,
      "admin_status": "up",
      "oper_status": "up",
      "tx_power_dbm": 2.5,
      "rx_power_dbm": -18.3,
      "ont_count": 32
    }
  ]
}
```

**PON Port Status Values:** `up`, `down`, `testing`, `unknown`, `dormant`, `not-present`, `lower-layer-down`

---

### POST /olt/onts

Fetches metrics for all ONTs registered on a ZTE C320 OLT.
Optionally filter by PON port.

**Request Body:**
```json
{
  "target": {
    "ip": "192.168.1.100",
    "community": "public"
  },
  "pon_port": 1
}
```

> Set `pon_port` to `0` or omit it to return ONTs from **all** PON ports.

**Response `200 OK`:**
```json
{
  "ip_address": "192.168.1.100",
  "total": 2,
  "onts": [
    {
      "ip_address": "192.168.1.100",
      "timestamp": "2026-02-18T02:50:00Z",
      "pon_port_index": 1,
      "ont_index": 1,
      "serial_number": "ZTEG12345678",
      "oper_status": "online",
      "rx_power_dbm": -22.1,
      "tx_power_dbm": 2.0,
      "distance_meters": 1500,
      "description": "Pelanggan A"
    }
  ]
}
```

**ONT Status Values:** `online`, `offline`, `unregistered`, `unknown`

---

## Realtime Execution (Mikrotik)

### POST /realtime/execute

Executes a command on a Mikrotik device via the Mikrotik API protocol.

**Request Body:**
```json
{
  "target": {
    "ip": "10.0.0.1",
    "driver": "mikrotik",
    "auth": {
      "username": "admin",
      "password": "secret",
      "port": 8728
    }
  },
  "command": "/system/resource/print"
}
```

**Response `200 OK`:**
```json
{
  "output": "...",
  "success": true
}
```

---

### POST /realtime/stats

Fetches real-time statistics from a Mikrotik device.

**Request Body:**
```json
{
  "target": {
    "ip": "10.0.0.1",
    "driver": "mikrotik",
    "auth": {
      "username": "admin",
      "password": "secret",
      "port": 8728
    }
  }
}
```

---

## Device Registry

> The device registry is used internally by go-nms for background monitoring.
> For on-demand OLT queries, use the [OLT Resources](#olt-resources-zte-c320-snmp) endpoints instead.

### GET /devices

Returns all registered devices.

**Response `200 OK`:**
```json
[
  {
    "id": "550e8400-e29b-41d4-a716-446655440000",
    "name": "OLT Core A",
    "ip_address": "192.168.1.100",
    "device_type": "olt",
    "protocol": "snmp",
    "status": "online"
  }
]
```

### POST /devices

Registers a new device for background monitoring.

**Request Body:**
```json
{
  "name": "OLT Core A",
  "ip_address": "192.168.1.100",
  "device_type": "olt",
  "protocol": "snmp",
  "description": "ZTE C320 - POP Utara"
}
```

**Device Types:** `router`, `switch`, `olt`, `ont`, `access_point`, `wireless`

**Protocols:** `mikrotik_api`, `ssh`, `telnet`, `tr069`, `snmp`

### GET /devices/:id

Returns a single device by UUID.

---

## Config Management

### POST /config/execute

Executes a configuration command on a device via SSH.

**Request Body:**
```json
{
  "device_id": "550e8400-e29b-41d4-a716-446655440000",
  "command": "show version"
}
```

---

## Inventory Sync

### POST /inventory/sync

Triggers a full inventory sync from openaccess to go-nms's device registry.
Called by openaccess when devices are created or updated.

**Response `200 OK`:**
```json
{
  "synced": 42,
  "errors": []
}
```
