package alert

// Rule represents a condition to trigger an alert
type Rule struct {
	ID          string  `json:"id"`
	DeviceID    string  `json:"device_id"` // Empty for global rules
	MetricName  string  `json:"metric_name"`
	Operator    string  `json:"operator"` // >, <, =, >=, <=
	Threshold   float64 `json:"threshold"`
	Description string  `json:"description"`
	Severity    string  `json:"severity"` // info, warning, critical
}
