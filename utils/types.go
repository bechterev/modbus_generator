package utils

type Register struct {
	Address uint16
	Value   uint8
}
type DoubleRegister struct {
	Address uint16
	Value   uint16
}
type ModbusDevice struct {
	address          int
	DiscreteInput    []Register
	Coils            []Register
	InputRegisters   []DoubleRegister
	HoldingRegisters []DoubleRegister
}

const (
	DiscreteInput = iota
	Coils
	InputRegisters
	HoldingRegisters
)

const (
	MaxCoils              = 10000
	MaxDiscrete           = 20000
	MaxInput              = 10000
	MaxHolding            = 25535
	BatchSize             = 1000
	CountRegisters        = 1000
	StartCoils            = 0
	StartDiscreteInput    = 10000
	StartInputRegisters   = 30000
	StartHoldingRegisters = 40000
)

var RegisterStartMap = map[int]int{
	Coils:            StartCoils,
	DiscreteInput:    StartDiscreteInput,
	InputRegisters:   StartInputRegisters,
	HoldingRegisters: StartHoldingRegisters,
}
