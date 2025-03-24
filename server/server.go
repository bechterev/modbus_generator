package server

import (
	"encoding/binary"
	"fmt"
	"generator/utils"
	"log"
	"net"
)

func handleRequest(data []byte) ([]byte, error) {
	if len(data) < 12 {
		return nil, ErrInvalidRequest // ignore short messages
	}

	// Parse the request
	transactionID := data[0:2]                          // ID transaction
	protocolID := data[2:4]                             // Modbus TCP has 0000
	length := data[4:6]                                 // Length of the message
	unitID := data[6]                                   // Device address
	functionCode := data[7]                             // Code of function (0x03 - read registers)
	startAddress := binary.BigEndian.Uint16(data[8:10]) // Start address
	count := binary.BigEndian.Uint16(data[10:12])       // Count of registers

	fmt.Printf("Got request: TransactionID=%x, ProtocolID=%x, Length=%x, UnitID=%x, FunctionCode=%x\n StartAddress=%d, Count=%d\n",
		transactionID, protocolID, length, unitID, functionCode, startAddress, count)
	device, exists := utils.Devices[int(unitID)]
	if !exists {
		return nil, ErrDeviceNotFound
	}

	response := append(transactionID, protocolID...)
	var byteCount uint16
	var responseData []byte
	var err error

	switch functionCode {
	case 0x03:
		byteCount = uint16(2 * count)
		responseData, err = extractDoubleRegisters(device.HoldingRegisters, startAddress, count)
	case 0x04:
		byteCount = uint16(2 * count)
		responseData, err = extractDoubleRegisters(device.InputRegisters, startAddress, count)
	case 0x01:
		byteCount = uint16(count)
		responseData, err = extractRegisters(device.Coils, startAddress, count)
	case 0x02:
		byteCount = uint16(count)
		responseData, err = extractRegisters(device.DiscreteInput, startAddress, count)
	default:
		return nil, ErrUnknownFunction
	}
	if err == ErrInvalidAddress || err == ErrInvalidCount {
		fmt.Println(err)
		response := make([]byte, 9)
		copy(response[0:2], transactionID) // Transaction ID
		copy(response[2:4], protocolID)    // Protocol ID
		response[4] = 0x00
		response[5] = 0x03 // Length (Unit ID + Function Code + Exception Code)
		response[6] = unitID
		response[7] = functionCode //  Function Code
		response[8] = 0x03         // For example Exception Code
		return response, nil
	}
	if err != nil {
		return nil, err
	}

	ln := uint16(3) + byteCount
	lengthField := make([]byte, 2)
	binary.BigEndian.PutUint16(lengthField, ln)
	response = append(response, lengthField...)
	response = append(response, unitID)
	response = append(response, functionCode)
	binary.BigEndian.PutUint16(lengthField, byteCount)
	response = append(response, lengthField[1])
	response = append(response, responseData...)
	return response, nil
}

func extractRegisters(registers []utils.Register, startAddress, count uint16) ([]byte, error) {
	if count > 250 {
		return nil, ErrInvalidCount
	}
	if int(startAddress)+int(count) > len(registers) {
		return nil, ErrInvalidAddress
	}

	data := make([]byte, count)
	for i := uint16(0); i < count; i++ {
		data[i] = registers[startAddress+uint16(i)].Value
	}
	return data, nil
}

func extractDoubleRegisters(registers []utils.DoubleRegister, startAddress, count uint16) ([]byte, error) {
	if count > 250 {
		return nil, ErrInvalidCount
	}
	if int(startAddress)+int(count) > len(registers) {
		return nil, ErrInvalidAddress
	}

	data := make([]byte, count*2)
	for i := uint16(0); i < count; i++ {
		binary.BigEndian.PutUint16(data[i*2:], registers[startAddress+uint16(i)].Value)
	}
	return data, nil
}

func StartModbusTCPServer(port string) {
	listener, err := net.Listen("tcp", port)
	if err != nil {
		log.Fatalf(ErrMsgServerFail, err)
	}
	defer listener.Close()
	fmt.Println(ErrMsgServerNotStarted, port)

	for {
		conn, err := listener.Accept()
		if err != nil {
			log.Println(ErrMsgServerFailConnect, err)
			continue
		}

		go func(c net.Conn) {
			defer c.Close()
			buffer := make([]byte, 1024)

			for {
				n, err := c.Read(buffer)
				if err != nil {
					log.Println(ErrMsgServerFailRead, err)
					break
				}

				response, err := handleRequest(buffer[:n])
				if err != nil {
					log.Println("Error handling request:", err)
					continue
				}

				if response != nil {
					_, err := c.Write(response)
					if err != nil {
						log.Println("Error writing response:", err)
						break
					}
				}
			}
		}(conn)
	}
}
