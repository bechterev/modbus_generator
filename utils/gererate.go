package utils

import (
	"encoding/binary"
	"math/rand"
	"sync"
)

var registerAddresses []uint16

func GenerateValues(device *ModbusDevice) {
	var (
		coilValues, discreteValues                                                                     []uint8
		inputValues, holdingValues, coilAddresses, discreteAddresses, holdingAddresses, inputAddresses []uint16
	)

	var wg sync.WaitGroup
	wg.Add(2)

	go func() {
		defer wg.Done()
		discreteValues = getValues[uint8](true, MaxDiscrete)
		coilValues = discreteValues[:]
		coilAddresses = registerAddresses[:StartDiscreteInput]
		discreteAddresses = registerAddresses[StartDiscreteInput:StartInputRegisters]
	}()

	go func() {
		defer wg.Done()
		holdingValues = getValues[uint16](false, MaxHolding)
		inputValues = holdingValues[:]
		inputAddresses = registerAddresses[StartInputRegisters:StartHoldingRegisters]
		holdingAddresses = registerAddresses[StartHoldingRegisters:]
	}()

	wg.Wait()

	wg.Add(2)
	go func() {
		defer wg.Done()
		for i := 0; i < MaxInput; i++ {
			device.Coils[i] = Register{coilAddresses[i], coilValues[i]}
			device.InputRegisters[i] = DoubleRegister{inputAddresses[i], inputValues[i]}
		}
	}()

	go func() {
		defer wg.Done()
		for i := 0; i < MaxHolding; i++ {
			if i < MaxDiscrete {
				device.DiscreteInput[i] = Register{discreteAddresses[i], discreteValues[i]}
			}
			device.HoldingRegisters[i] = DoubleRegister{holdingAddresses[i], holdingValues[i]}
		}
	}()

	wg.Wait()
}

func InitRegisters() {
	registerAddresses = make([]uint16, 65535)
	for i := 0; i < 65535; i++ {
		registerAddresses[i] = uint16(i)
	}
}

func getValues[T uint8 | uint16](isBool bool, sizeRegister int) []T {
	values := make([]T, sizeRegister)
	buf := make([]byte, sizeRegister*2)

	// bulk-generate random data
	_, err := rand.Read(buf)
	if err != nil {
		panic(err)
	}

	if isBool {
		for i := range values {
			values[i] = T(buf[i] & 1)
		}
	} else {
		for i := range values {
			values[i] = T(binary.LittleEndian.Uint16(buf[i*2 : i*2+2]))
		}
	}

	return values
}
