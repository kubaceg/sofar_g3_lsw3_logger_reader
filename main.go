package main

import (
	"log"
	_ "net/http/pprof"
	"strings"
	"time"

	gser "go.bug.st/serial"

	"github.com/kubaceg/sofar_g3_lsw3_logger_reader/adapters/comms/serial"
	"github.com/kubaceg/sofar_g3_lsw3_logger_reader/adapters/comms/tcpip"
	"github.com/kubaceg/sofar_g3_lsw3_logger_reader/adapters/databases/mosquitto"
	"github.com/kubaceg/sofar_g3_lsw3_logger_reader/adapters/devices/sofar"

	"github.com/kubaceg/sofar_g3_lsw3_logger_reader/ports"
)

var (
	config *Config
	port   ports.CommunicationPort
	mqtt   ports.DatabaseWithListener
	device ports.Device

	hasMQTT bool
)

func initialize() {
	var err error
	config, err = NewConfig("config.yaml")
	if err != nil {
		log.Fatalln(err)
	}

	hasMQTT = config.Mqtt.Url != "" && config.Mqtt.Prefix != ""

	if isSerialPort(config.Inverter.Port) {
		port = serial.New(config.Inverter.Port, 2400, 8, gser.NoParity, gser.OneStopBit)
		log.Printf("using serial communcations port %s", config.Inverter.Port)
	} else {
		port = tcpip.New(config.Inverter.Port)
		log.Printf("using TCP/IP communications port %s", config.Inverter.Port)
	}

	if hasMQTT {
		mqtt, err = mosquitto.New(&config.Mqtt)
		if err != nil {
			log.Fatalf("MQTT connection failed: %s", err)
		}

	}

	device = sofar.NewSofarLogger(config.Inverter.LoggerSerial, port)
}

func main() {
	initialize()

	for {
		log.Printf("performing measurements")
		timeStart := time.Now()

		measurements, err := device.Query()
		if err != nil {
			log.Printf("failed to perform measurements: %s", err)
			time.Sleep(5 * time.Second)
			continue
		}

		if hasMQTT {
			err = mqtt.InsertRecord(measurements)
			if err != nil {
				log.Printf("failed to insert record to MQTT: %s", err)
			}
		}

		duration := time.Since(timeStart)

		delay := time.Duration(config.Inverter.ReadInterval)*time.Second - duration
		if delay <= 0 {
			delay = 1 * time.Second
		}

		time.Sleep(delay)
	}

}

func isSerialPort(portName string) bool {
	return strings.HasPrefix(portName, "/")
}
