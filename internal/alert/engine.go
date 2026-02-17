package alert

import (
	"encoding/json"
	"fmt"
	"log"

	"github.com/nats-io/nats.go"
	commonModel "github.com/yourorg/nms-go/internal/common/model"
	"github.com/yourorg/nms-go/internal/notification"
)

type Engine struct {
	natsConn       *nats.Conn
	notifier       notification.Service
	rules          []Rule
	stopChan       chan struct{}
}

func NewEngine(nc *nats.Conn, notifier notification.Service) *Engine {
	// Hardcoded rules for MVP
	rules := []Rule{
		{
			ID:          "rule-1",
			MetricName:  "rtt_ms",
			Operator:    ">",
			Threshold:   100.0,
			Description: "High Latency (>100ms)",
			Severity:    "warning",
		},
		{
			ID:          "rule-2",
			MetricName:  "success",
			Operator:    "=",
			Threshold:   0.0, // false becomes 0.0
			Description: "Device Down",
			Severity:    "critical",
		},
	}

	return &Engine{
		natsConn:       nc,
		notifier:       notifier,
		rules:          rules,
		stopChan:       make(chan struct{}),
	}
}

func (e *Engine) Start() {
	log.Println("Alert Engine started, subscribing to nms.metrics")

	sub, err := e.natsConn.Subscribe("nms.metrics", func(msg *nats.Msg) {
		var metric commonModel.Metric
		if err := json.Unmarshal(msg.Data, &metric); err != nil {
			log.Printf("Error unmarshalling metric: %v", err)
			return
		}

		e.evaluate(metric)
	})

	if err != nil {
		log.Fatalf("Error communicating with NATS: %v", err)
	}
	defer sub.Unsubscribe()

	<-e.stopChan
}

func (e *Engine) Stop() {
	close(e.stopChan)
}

func (e *Engine) evaluate(metric commonModel.Metric) {
	for _, rule := range e.rules {
		// specific device check (if rule has DeviceID)
		if rule.DeviceID != "" && rule.DeviceID != metric.DeviceID {
			continue
		}

		val, ok := metric.Values[rule.MetricName]
		if !ok {
			continue
		}

		// Convert boolean to float for comparison if needed
		floatVal, ok := toFloat(val)
		if !ok {
			continue
		}

		triggered := false
		switch rule.Operator {
		case ">":
			triggered = floatVal > rule.Threshold
		case "<":
			triggered = floatVal < rule.Threshold
		case "=":
			triggered = floatVal == rule.Threshold
		case ">=":
			triggered = floatVal >= rule.Threshold
		case "<=":
			triggered = floatVal <= rule.Threshold
		}

		if triggered {
			alertMsg := fmt.Sprintf("ALERT [%s]: Device %s (%s) - %s (Value: %.2f)", 
				rule.Severity, metric.DeviceName, metric.IPAddress, rule.Description, floatVal)
			
			log.Println("⚡ " + alertMsg)
			e.notifier.Send("admin@example.com", "NMS Alert: "+rule.Description, alertMsg)
		}
	}
}

func toFloat(unk interface{}) (float64, bool) {
	switch v := unk.(type) {
	case float64:
		return v, true
	case int:
		return float64(v), true
	case bool:
		if v {
			return 1.0, true
		}
		return 0.0, true
	default:
		return 0, false
	}
}
