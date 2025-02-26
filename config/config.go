package config

import (
	"log"
	"os"
	"strconv"

	"generator/server"

	"github.com/joho/godotenv"
)

type Config struct {
	Port        string
	DeviceCount int
	TimeSleep   int
}

var Cfg Config

func LoadConfig() {
	err := godotenv.Load()
	if err != nil {
		log.Println(server.ErrMsgFailLoadConfig)
	}

	Cfg.Port = os.Getenv("PORT")
	if Cfg.Port == "" {
		Cfg.Port = "1502" // default port
	}

	deviceCountStr := os.Getenv("DEVICE_COUNT")
	deviceCount, err := strconv.Atoi(deviceCountStr)
	if err != nil {
		log.Fatalf(server.ErrMsgCovertConfig, err)
	}
	Cfg.DeviceCount = deviceCount

	timeSleepStr := os.Getenv("TIME_SLEEP")
	timeSleep, err := strconv.Atoi(timeSleepStr)
	if err != nil {
		log.Fatalf(server.ErrMsgCovertConfig, err)
	}
	Cfg.TimeSleep = timeSleep
}
