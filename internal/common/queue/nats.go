package queue

import (
	"fmt"

	"github.com/nats-io/nats.go"
	"github.com/yourorg/nms-go/internal/common/config"
)

func NewNATSConnection(cfg config.NATSConfig) (*nats.Conn, error) {
	nc, err := nats.Connect(cfg.URL)
	if err != nil {
		return nil, fmt.Errorf("failed to connect to nats: %w", err)
	}

	return nc, nil
}
