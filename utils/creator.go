package utils

func CreatorDevices(count int) map[int]*ModbusDevice {
	devices := make(map[int]*ModbusDevice)
	for i := 0; i < count; i++ {
		devices[i] = &ModbusDevice{
			address:          i,
			DiscreteInput:    make([]Register, MaxDiscrete),
			Coils:            make([]Register, MaxCoils),
			InputRegisters:   make([]DoubleRegister, MaxInput),
			HoldingRegisters: make([]DoubleRegister, MaxHolding),
		}
	}
	return devices
}
