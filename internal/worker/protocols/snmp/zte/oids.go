// Package zte provides an SNMP adapter for ZTE C320 OLT devices.
package zte

// ZTE C320 SNMP OID constants.
// These OIDs are sourced from the ZTE C320 MIB documentation.
// All OIDs use the standard SNMP format (dot-separated integers).
const (
	// --- System OIDs (RFC 1213 / SNMPv2-MIB) ---

	// --- System OIDs (RFC 1213 / SNMPv2-MIB) ---

	// OIDSysDescr is the textual description of the entity (e.g., "ZTE C320 OLT").
	OIDSysDescr = "1.3.6.1.2.1.1.1.0"

	// OIDSysUpTime is the time (in hundredths of a second) since the network management
	// portion of the system was last re-initialized.
	OIDSysUpTime = "1.3.6.1.2.1.1.3.0"

	// OIDSysName is the administratively-assigned name for this managed node.
	OIDSysName = "1.3.6.1.2.1.1.5.0"

	// --- ZTE System OIDs (Discovered from Device Walk) ---

	// The 1015.2 tree contains card information (Shelf/Slot).
	// We will use the first valid card found for system-wide metrics or average them.

	// OIDZTECardTable - Base for card metrics
	OIDZTECardTable = "1.3.6.1.4.1.3902.1015.2.1.1.3.1"

	// OIDZTECardCPUUsage - .1.3.6.1.4.1.3902.1015.2.1.1.3.1.9.1.1 (Values: 22, 17, 23)
	// This is a table column, indexed by card position.
	OIDZTECardCPUUsage = "1.3.6.1.4.1.3902.1015.2.1.1.3.1.9.1.1"

	// OIDZTECardTemperature - .1.3.6.1.4.1.3902.1015.2.1.1.3.1.8.1.1 (Values: 43, 42)
	OIDZTECardTemperature = "1.3.6.1.4.1.3902.1015.2.1.1.3.1.8.1.1"

	// OIDZTECardMemoryUsage - .1.3.6.1.4.1.3902.1015.2.1.1.3.1.11.1.1 (Values: 41, 39, 29)
	// Assuming this is percentage based on values.
	OIDZTECardMemoryUsage = "1.3.6.1.4.1.3902.1015.2.1.1.3.1.11.1.1"

	// OIDZTECardMemoryTotal - Not explicitly found in %-like table.
	// We can try .1.3.6.1.4.1.3902.1015.2.1.1.3.1.19.1.1 (Values: 512, 512, 2048) in MB?
	OIDZTECardMemoryTotal = "1.3.6.1.4.1.3902.1015.2.1.1.3.1.19.1.1"

	// --- ZTE PON Port OIDs (1.3.6.1.4.1.3902.1015.3.1) ---

	// OIDZTEPONPortTable
	OIDZTEPONPortTable = "1.3.6.1.4.1.3902.1015.3.1.2.1"

	// OIDZTEPONPortIndex - .1.3.6.1.4.1.3902.1015.3.1.2.1.1 (1, 2, 1)
	OIDZTEPONPortIndex = "1.3.6.1.4.1.3902.1015.3.1.2.1.1"

	// OIDZTEPONPortType - .1.3.6.1.4.1.3902.1015.3.1.1.1.5 (1=GPON?)
	OIDZTEPONPortType = "1.3.6.1.4.1.3902.1015.3.1.1.1.5"

	// OIDZTEPONPortAdminStatus - .1.3.6.1.4.1.3902.1015.3.1.2.1.3 (3=testing? 4=?)
	OIDZTEPONPortAdminStatus = "1.3.6.1.4.1.3902.1015.3.1.2.1.3"

	// OIDZTEPONPortOperStatus - .1.3.6.1.4.1.3902.1015.3.1.2.1.4 (1=up, 2=down)
	OIDZTEPONPortOperStatus = "1.3.6.1.4.1.3902.1015.3.1.2.1.4"

	// OIDZTEPONPortTxPower - .1.3.6.1.4.1.3902.1015.3.1.3.1.12 (Values: 115263 -> 11.5 dBm)
	OIDZTEPONPortTxPower = "1.3.6.1.4.1.3902.1015.3.1.3.1.12"

	// OIDZTEPONPortRxPower - .1.3.6.1.4.1.3902.1015.3.1.3.1.10 (Values: 44078864?)
	// Or maybe .1.3.6.1.4.1.3902.1015.3.1.3.1.6 (Values: 12995588)?
	OIDZTEPONPortRxPower = "1.3.6.1.4.1.3902.1015.3.1.3.1.10"

	// OIDZTEPONPortONTCount - .1.3.6.1.4.1.3902.1015.3.1.3.1.13 (Values: 11, 24)
	OIDZTEPONPortONTCount = "1.3.6.1.4.1.3902.1015.3.1.3.1.13"

	// --- ZTE ONT OIDs (Discovered) ---
	// Found table: 1.3.6.1.4.1.3902.1015.3.1.13.1
	// Index appears to be complex (Circuit ID?), but contains metrics.

	OIDZTEONTTable = "1.3.6.1.4.1.3902.1015.3.1.13.1"

	// .3 = Status? (Value 30) - treating as OperStatus for now
	OIDZTEONTOperStatus = "1.3.6.1.4.1.3902.1015.3.1.13.1.3"

	// .4 = Distance? (Value ~7000)
	OIDZTEONTDistance = "1.3.6.1.4.1.3902.1015.3.1.13.1.4"

	// .5 = RxPower? (Value -140 -> -14.0dBm)
	OIDZTEONTRxPower = "1.3.6.1.4.1.3902.1015.3.1.13.1.5"

	// .6 = ? (Value 90) - maybe TxPower?
	OIDZTEONTTxPower = "1.3.6.1.4.1.3902.1015.3.1.13.1.6"

	// Serial Number not found in this table, using index as makeshift SN.
	// We'll update client.go to handle this.
	OIDZTEONTSerialNumber = "1.3.6.1.4.1.3902.1015.3.1.13.1.1" // Placeholder
	OIDZTEONTDescription  = "1.3.6.1.4.1.3902.1015.3.1.13.1.1" // Placeholder
)

// PONPortStatus represents the operational status of a PON port.
type PONPortStatus int

const (
	PONPortStatusUp      PONPortStatus = 1
	PONPortStatusDown    PONPortStatus = 2
	PONPortStatusUnknown PONPortStatus = 0
)

// String returns a human-readable representation of the PON port status.
func (s PONPortStatus) String() string {
	switch s {
	case PONPortStatusUp:
		return "up"
	case PONPortStatusDown:
		return "down"
	default:
		return "unknown"
	}
}

// ONTStatus represents the operational status of an ONT.
type ONTStatus int

const (
	ONTStatusOnline  ONTStatus = 1
	ONTStatusOffline ONTStatus = 2
	ONTStatusUnreg   ONTStatus = 3
	ONTStatusUnknown ONTStatus = 0
)

// String returns a human-readable representation of the ONT status.
func (s ONTStatus) String() string {
	switch s {
	case ONTStatusOnline:
		return "online"
	case ONTStatusOffline:
		return "offline"
	case ONTStatusUnreg:
		return "unregistered"
	default:
		return "unknown"
	}
}
