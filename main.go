package main

import (
	"flag"
	"log"
	"strings"
	"time"

	"net/http"
	_ "net/http/pprof"

	"git.xelasys.ro/sigxcpu/sofar/adapters/comms/serial"
	"git.xelasys.ro/sigxcpu/sofar/adapters/comms/tcpip"
	"git.xelasys.ro/sigxcpu/sofar/adapters/databases/httpnow"
	"git.xelasys.ro/sigxcpu/sofar/adapters/databases/influx"
	"git.xelasys.ro/sigxcpu/sofar/adapters/databases/mosquitto"
	"git.xelasys.ro/sigxcpu/sofar/adapters/devices/sofar"
	gser "go.bug.st/serial"

	"git.xelasys.ro/sigxcpu/sofar/ports"
)

var (
	port               ports.CommunicationPort
	db                 ports.Database
	mqtt               ports.DatabaseWithListener
	nowDB              ports.Database
	device             ports.Device
	lastDateTimeUpdate time.Time
	loggerSerial       string

	hasInflux bool
	hasMQTT   bool
)

func init() {
	var portName, dbURL, dbName, mqttURL string

	flag.StringVar(&portName, "port", "", "port name (e.g. /dev/ttyUSB0 for serial or 1.2.3.4:23 for TCP/IP")
	flag.StringVar(&dbURL, "influx-url", "", "Influx DB URL for push (e.g. http://localhost:8086")
	flag.StringVar(&dbName, "influx-db", "", "Influx DB database name")
	flag.StringVar(&mqttURL, "mqtt-url", "", "MQTT broker URL (e.g. tcp://1.2.3.4:5678")
	flag.StringVar(&loggerSerial, "logger-serial", "", "Logger serial number")
	flag.Parse()

	hasInflux = dbURL != "" && dbName != ""
	hasMQTT = mqttURL != ""

	// check config
	if portName == "" {
		flag.Usage()
		log.Fatalf("please specify port")
	}

	if isSerialPort(portName) {
		port = serial.New(portName, 2400, 8, gser.NoParity, gser.OneStopBit)
		log.Printf("using serial communcations port %s", portName)
	} else {
		port = tcpip.New(portName)
		log.Printf("using TCP/IP communications port %s", portName)
	}

	if hasInflux {
		log.Printf("starting with Influx URL: %s, Database: %s", dbURL, dbName)

		db = influx.New(dbURL, dbName, map[string]string{"device": "sofar"})
	}

	var err error
	if hasMQTT {
		mqtt, err = mosquitto.New(mqttURL, "/sensors/energy/inverter2")
		if err != nil {
			log.Fatalf("MQTT connection failed: %s", err)
		}

	}

	device = sofar.NewSofarLogger(loggerSerial, port)

	nowDB = httpnow.NewHttpNow(8081)

	go func() {
		log.Println(http.ListenAndServe("localhost:6060", nil))
	}()

}

func main() {

	for {
		log.Printf("performing measurements")
		timeStart := time.Now()

		measurements, err := device.Query()
		if err != nil {
			log.Printf("failed to perform measurements: %s", err)
			time.Sleep(5 * time.Second)
			continue
		}

		if hasInflux {
			err = db.InsertRecord(measurements)

			if err != nil {
				log.Printf("failed to insert into database: %s", err)
			}
		}

		if hasMQTT {
			err = mqtt.InsertRecord(measurements)
			if err != nil {
				log.Printf("failed to insert record to MQTT: %s", err)
			}
		}

		nowDB.InsertRecord(measurements)

		duration := time.Since(timeStart)

		delay := 10*time.Second - duration
		if delay <= 0 {
			delay = 1 * time.Second
		}

		time.Sleep(delay)
	}

}

func isSerialPort(portName string) bool {
	return strings.HasPrefix(portName, "/")
}
