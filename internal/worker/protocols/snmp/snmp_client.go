// Package snmp provides a generic SNMP client interface and implementation
// for use by device-specific adapters in go-nms.
package snmp

import (
	"context"
	"fmt"
	"time"

	"github.com/gosnmp/gosnmp"
)

// SNMPClient defines the interface for SNMP operations.
// Device-specific adapters depend on this interface, enabling easy mocking in tests.
type SNMPClient interface {
	// Connect establishes an SNMP session to the target host.
	Connect(ctx context.Context, host, community string, version gosnmp.SnmpVersion, timeout time.Duration) error

	// Disconnect closes the SNMP session.
	Disconnect() error

	// Get retrieves the values for the given OIDs.
	Get(oids []string) (*gosnmp.SnmpPacket, error)

	// Walk performs an SNMP walk starting from the given OID,
	// calling fn for each returned PDU.
	Walk(oid string, fn gosnmp.WalkFunc) error

	// GetBulk performs an SNMP GETBULK request for the given OIDs.
	// nonRepeaters is uint8 and maxRepetitions is uint32 to match the gosnmp API.
	GetBulk(oids []string, nonRepeaters uint8, maxRepetitions uint32) (*gosnmp.SnmpPacket, error)
}

// GoSNMPClient is the production implementation of SNMPClient backed by gosnmp.
type GoSNMPClient struct {
	snmp *gosnmp.GoSNMP
}

// NewGoSNMPClient creates a new GoSNMPClient with sensible defaults.
func NewGoSNMPClient() *GoSNMPClient {
	return &GoSNMPClient{}
}

// Connect establishes an SNMP session.
func (c *GoSNMPClient) Connect(ctx context.Context, host, community string, version gosnmp.SnmpVersion, timeout time.Duration) error {
	c.snmp = &gosnmp.GoSNMP{
		Target:             host,
		Port:               161,
		Community:          community,
		Version:            version,
		Timeout:            timeout,
		Retries:            2,
		ExponentialTimeout: true,
		MaxOids:            gosnmp.MaxOids,
	}

	if err := c.snmp.ConnectIPv4(); err != nil {
		return fmt.Errorf("snmp connect to %s failed: %w", host, err)
	}

	return nil
}

// Disconnect closes the SNMP session.
func (c *GoSNMPClient) Disconnect() error {
	if c.snmp != nil && c.snmp.Conn != nil {
		return c.snmp.Conn.Close()
	}

	return nil
}

// Get retrieves the values for the given OIDs.
func (c *GoSNMPClient) Get(oids []string) (*gosnmp.SnmpPacket, error) {
	if c.snmp == nil {
		return nil, fmt.Errorf("snmp client not connected")
	}

	packet, err := c.snmp.Get(oids)
	if err != nil {
		return nil, fmt.Errorf("snmp get failed: %w", err)
	}

	return packet, nil
}

// Walk performs an SNMP walk starting from the given OID.
func (c *GoSNMPClient) Walk(oid string, fn gosnmp.WalkFunc) error {
	if c.snmp == nil {
		return fmt.Errorf("snmp client not connected")
	}

	if err := c.snmp.BulkWalk(oid, fn); err != nil {
		return fmt.Errorf("snmp walk on %s failed: %w", oid, err)
	}

	return nil
}

// GetBulk performs an SNMP GETBULK request.
func (c *GoSNMPClient) GetBulk(oids []string, nonRepeaters uint8, maxRepetitions uint32) (*gosnmp.SnmpPacket, error) {
	if c.snmp == nil {
		return nil, fmt.Errorf("snmp client not connected")
	}

	packet, err := c.snmp.GetBulk(oids, nonRepeaters, maxRepetitions)
	if err != nil {
		return nil, fmt.Errorf("snmp getbulk failed: %w", err)
	}

	return packet, nil
}
