package main

import (
	"generator/server"
	"generator/utils"
	"time"

	"generator/config"
)

func main() {
	config.LoadConfig()
	utils.InitRegisters()
	utils.InitDevices(config.Cfg.DeviceCount)
	go utils.UpdateValues(time.Duration(config.Cfg.TimeSleep))
	server.StartModbusTCPServer(":" + config.Cfg.Port)
}
