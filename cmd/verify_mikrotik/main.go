package main

import (
	"flag"
	"fmt"
	"log"
	
	"github.com/yourorg/nms-go/internal/worker"
)

func main() {
	ip := flag.String("ip", "", "Mikrotik IP address")
	user := flag.String("user", "admin", "Username")
	pass := flag.String("pass", "admin", "Password")
	flag.Parse()

	if *ip == "" {
		log.Fatal("Please provide an IP address using -ip")
	}

	adapter := worker.NewMikrotikAdapter()
	fmt.Printf("Connecting to %s as %s...\n", *ip, *user)

	metrics, ok := adapter.FetchSystemResources(*ip, *user, *pass)
	if !ok {
		log.Fatalf("Failed to fetch resources. Check connectivity and credentials.")
	}

	fmt.Println("Successfully fetched metrics:")
	for k, v := range metrics {
		fmt.Printf(" - %s: %v\n", k, v)
	}
}
