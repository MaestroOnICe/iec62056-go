package main

import (
	"fmt"
	"log"
	"time"

	"github.com/MaestroOnICe/iec62056-go"
)

func main() {
	config := iec62056.SerialConfig{
		PortName: "/dev/ttyUSB0",
		DeviceConfig: iec62056.DeviceConfig{
			InitialBaudRate: iec62056.Baud300,
			Mode:            iec62056.ModeC,
			Timeout:         time.Second * 5,
			ReadOnly:        false,
		},
	}

	client, err := iec62056.NewSerialClient(config)
	if err != nil {
		log.Fatalf("Failed to create client: %v", err)
	}
	defer client.Close()

	msg, err := client.ReadMeter()
	if err != nil {
		log.Fatalf("Failed to read meter: %v", err)
	}

	fmt.Printf("Manufacturer: %s\n", msg.Manufacturer)
	fmt.Printf("Mode: %c\n", msg.Mode)
	fmt.Printf("Baud Rate: %d\n", msg.BaudRate)

	for _, block := range msg.DataBlocks {
		fmt.Printf("Address: %s, Value: %s, Unit: %s\n",
			block.Address, block.Value, block.Unit)
	}
}
