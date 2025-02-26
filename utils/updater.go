package utils

import (
	"fmt"
	"time"
)

var Devices map[int]*ModbusDevice

func InitDevices(count int) {
	Devices = CreatorDevices(count)
}

// UpdateValues updates the values of the devices
func UpdateValues(timeUpdate time.Duration) {
	for {
		for _, v := range Devices {
			go GenerateValues(v)
		}
		time.Sleep(timeUpdate * time.Millisecond)
	}
}

func PrintDevice(key int) {
	device, exists := Devices[key]
	if !exists {
		fmt.Printf("Device with key %d not found\n", key)
		return
	}

	fmt.Printf("Device at address %d:\n", device.address)
	fmt.Println("Discrete Inputs:", device.DiscreteInput)
	fmt.Println("Coils:", device.Coils)
	fmt.Println("Input Registers:", device.InputRegisters)
	fmt.Println("Holding Registers:", device.HoldingRegisters)
}
