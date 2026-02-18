package main

import (
	"fmt"
	"os"
	"time"

	"github.com/gosnmp/gosnmp"
)

func main() {
	ip := "10.70.126.102"
	community := "jinomro"
	if len(os.Args) > 1 {
		ip = os.Args[1]
	}
	if len(os.Args) > 2 {
		community = os.Args[2]
	}

	g := &gosnmp.GoSNMP{
		Target:    ip,
		Port:      161,
		Community: community,
		Version:   gosnmp.Version2c,
		Timeout:   5 * time.Second,
		Retries:   2,
	}

	if err := g.Connect(); err != nil {
		fmt.Printf("Connect error: %v\n", err)
		os.Exit(1)
	}
	defer g.Conn.Close()

	// Walk 1.3.6.1.2.1.1 (Standard System OIDs)
	fmt.Println("\n--- Walk 1.3.6.1.2.1.1 (Standard System) ---")
	err := g.Walk("1.3.6.1.2.1.1", func(pdu gosnmp.SnmpPDU) error {
		fmt.Printf("%-30s %-10s %v\n", pdu.Name, pdu.Type, pdu.Value)
		return nil
	})
	if err != nil {
		fmt.Printf("Walk error: %v\n", err)
	}
}
