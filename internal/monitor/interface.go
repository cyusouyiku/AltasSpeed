package main

import (
	"fmt"
	"log"

	netio "github.com/shirou/gopsutil/v4/net"
)

func main() {
	stats, err := netio.IOCounters(true)
	if err != nil {
		log.Fatalf("failed to read network stats: %v", err)
	}

	if len(stats) == 0 {
		fmt.Println("No network interfaces found.")
		return
	}

	for _, iface := range stats {
		if iface.BytesRecv == 0 && iface.BytesSent == 0 {
			continue
		}

		fmt.Printf("%s -> RX: %d bytes, TX: %d bytes\n", iface.Name, iface.BytesRecv, iface.BytesSent)
	}
}