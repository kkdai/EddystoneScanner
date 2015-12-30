package main

import (
	"fmt"
	"log"
	"strings"

	"github.com/paypal/gatt"
	"github.com/paypal/gatt/examples/option"
)

var done = make(chan struct{})

const eddystoneServicesUUID string = "FEAA"

func validEddystone(serviceData []gatt.UUID) bool {
	if len(serviceData) > 0 && strings.ToUpper(serviceData[0].String()) == eddystoneServicesUUID {
		return true
	}
	return false
}

func onStateChanged(d gatt.Device, s gatt.State) {
	fmt.Println("State:", s)
	switch s {
	case gatt.StatePoweredOn:
		fmt.Println("Scanning for eddystone beacon...")
		d.Scan([]gatt.UUID{}, false)

		fmt.Println("=====================================")
		return
	default:
		d.StopScanning()
	}
}

func onPeriphDiscovered(p gatt.Peripheral, a *gatt.Advertisement, rssi int) {
	if len(a.Services) == 0 {
		return
	}

	if validEddystone(a.Services) {
		fmt.Println("Eddystone Beacon Found!.....")
		ed := NewEddystoneParser(a)
		ed.PrintBeacon()
		fmt.Println("-------------------------------------")
	}
	return
}

func main() {
	d, err := gatt.NewDevice(option.DefaultClientOptions...)
	if err != nil {
		log.Fatalf("Failed to open device, err: %s\n", err)
		return
	}

	// Register handlers.
	d.Handle(
		gatt.PeripheralDiscovered(onPeriphDiscovered),
	)

	d.Init(onStateChanged)
	<-done
	fmt.Println("Done")
}
