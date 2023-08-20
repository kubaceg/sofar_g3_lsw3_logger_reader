package mosquitto

import (
	"encoding/json"
	"fmt"
	"log"
	"time"

	mqtt "github.com/eclipse/paho.mqtt.golang"
)

type MqttConfig struct {
	Url      string `yaml:"url"`
	User     string `yaml:"user"`
	Password string `yaml:"password"`
	Prefix   string `yaml:"prefix"`
}

type Connection struct {
	client mqtt.Client
	prefix string
}

var connectHandler mqtt.OnConnectHandler = func(client mqtt.Client) {
	log.Printf("MQTT Connected")
}

var connectLostHandler mqtt.ConnectionLostHandler = func(client mqtt.Client, err error) {
	log.Printf("Connect lost: %v", err)
}

func New(config *MqttConfig) (*Connection, error) {
	opts := mqtt.NewClientOptions()
	opts.AddBroker(config.Url)
	opts.SetClientID("sofar")
	opts.OnConnect = connectHandler
	opts.OnConnectionLost = connectLostHandler

	if config.User != "" {
		opts.SetUsername(config.User)
	}

	if config.Password != "" {
		opts.SetPassword(config.Password)
	}

	conn := &Connection{}
	conn.client = mqtt.NewClient(opts)
	conn.prefix = config.Prefix
	if token := conn.client.Connect(); token.Wait() && token.Error() != nil {
		return nil, token.Error()
	}

	return conn, nil

}

func publish(conn *Connection, k string, v interface{}) {
	token := conn.client.Publish(fmt.Sprintf("%s/%s", conn.prefix, k), 0, true, fmt.Sprintf("%v", v))
	res := token.WaitTimeout(1 * time.Second)
	if !res || token.Error() != nil {
		log.Printf("error inserting to MQTT: %s", token.Error())
	}
}

// next thing to do is add discovery
// func discovery() {
// MQTT Discovery: https://www.home-assistant.io/integrations/mqtt/#mqtt-discovery

// homeassistant/sensor/inverter/config
// payload
// unique_id: PV_Generation_Today01ad
// state_topic: "homeassistant/sensor/inverter/PV_Generation_Today"
// device_class: energy
// state_class: measurement
// unit_of_measurement: 'kWh'
// value_template: "{{ value|int * 0.01 }}"```

// {"name": null, "device_class": "motion", "state_topic": "homeassistant/binary_sensor/garden/state", "unique_id": "motion01ad", "device": {"identifiers": ["01ad"], "name": "Garden" }}
//}

func (conn *Connection) InsertRecord(measurement map[string]interface{}) error {
	m := make(map[string]interface{}, len(measurement))
	for k, v := range measurement {
		m[k] = v
	}
	m["LastTimestamp"] = time.Now().UnixNano() / int64(time.Millisecond)
	all, _ := json.Marshal(m)
	m["All"] = string(all)
	for k, v := range m {
		publish(conn, k, v)
	}
	return nil
}

func (conn *Connection) Subscribe(topic string, callback mqtt.MessageHandler) {
	conn.client.Subscribe(topic, 0, callback)
}
