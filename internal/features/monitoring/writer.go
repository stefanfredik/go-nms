package monitoring

import (
	"time"

	influxdb2 "github.com/influxdata/influxdb-client-go/v2"
	"github.com/influxdata/influxdb-client-go/v2/api"
	"github.com/yourorg/nms-go/internal/worker/protocols/mikrotik"
)

// MetricWriter defines how metrics are stored
type MetricWriter interface {
	WriteSystemMetrics(metrics *mikrotik.SystemMetrics)
	WriteInterfaceMetrics(metrics []*mikrotik.InterfaceMetrics)
	Close()
}

type InfluxDBWriter struct {
	client   influxdb2.Client
	writeAPI api.WriteAPI
}

func NewInfluxDBWriter(url, token, org, bucket string) *InfluxDBWriter {
	client := influxdb2.NewClient(url, token)
	writeAPI := client.WriteAPI(org, bucket)

	return &InfluxDBWriter{
		client:   client,
		writeAPI: writeAPI,
	}
}

func (w *InfluxDBWriter) WriteSystemMetrics(m *mikrotik.SystemMetrics) {
	p := influxdb2.NewPointWithMeasurement("system_metrics").
		AddTag("device_id", m.DeviceID).
		AddField("cpu_usage", m.CPUUsage).
		AddField("memory_usage", m.MemoryUsage).
		AddField("uptime", m.Uptime).
		SetTime(time.Now())

	w.writeAPI.WritePoint(p)
}

func (w *InfluxDBWriter) WriteInterfaceMetrics(metrics []*mikrotik.InterfaceMetrics) {
	for _, m := range metrics {
		p := influxdb2.NewPointWithMeasurement("interface_metrics").
			AddTag("device_id", m.DeviceID).
			AddTag("interface", m.InterfaceName).
			AddField("bytes_in", m.BytesIn).
			AddField("bytes_out", m.BytesOut).
			SetTime(time.Now())

		w.writeAPI.WritePoint(p)
	}
}

func (w *InfluxDBWriter) Close() {
	w.client.Close()
}
