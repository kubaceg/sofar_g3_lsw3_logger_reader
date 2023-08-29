package mosquitto

import (
	"encoding/json"
	"fmt"
	"log"
	"time"

	mqtt "github.com/eclipse/paho.mqtt.golang"
	"github.com/kubaceg/sofar_g3_lsw3_logger_reader/ports"
)

type MqttConfig struct {
	Url       string `yaml:"url"`
	User      string `yaml:"user"`
	Password  string `yaml:"password"`
	Discovery string `yaml:"discovery"`
	State     string `yaml:"state"`
}

type Connection struct {
	client mqtt.Client
	state  string
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
	conn.state = config.State
	if token := conn.client.Connect(); token.Wait() && token.Error() != nil {
		return nil, token.Error()
	}

	return conn, nil

}

func (conn *Connection) publish(topic string, msg string, retain bool) {
	token := conn.client.Publish(topic, 0, retain, msg)
	res := token.WaitTimeout(1 * time.Second)
	if !res || token.Error() != nil {
		log.Printf("error inserting to MQTT: %s", token.Error())
	}
}

// return "power" for kW etc., "energy" for kWh etc.
// func unit2DeviceClass(unit string) string {
// 	if strings.HasSuffix(unit, "Wh") {
// 		return "energy"
// 	} else if strings.HasSuffix(unit, "W") {
// 		return "power"
// 	} else {
// 		return ""
// 	}
// }

// MQTT Discovery: https://www.home-assistant.io/integrations/mqtt/#mqtt-discovery
func (conn *Connection) InsertDiscoveryRecord(discovery string, state string, fields []ports.DiscoveryField) error {
	uniq := "01ad" // TODO: get from config?
	for _, f := range fields {
		topic := fmt.Sprintf("%s/%s/config", discovery, f.Name)
		json, _ := json.Marshal(map[string]interface{}{
			"name":      f.Name,
			"unique_id": fmt.Sprintf("%s_%s", f.Name, uniq),
			// "device_class": unit2DeviceClass(f.Unit),  // TODO: not working, always "energy"
			// "state_class": "measurement",
			"state_topic":         state,
			"unit_of_measurement": f.Unit,
			"value_template":      fmt.Sprintf("{{ value_json.%s|int * %s }}", f.Name, f.Factor),
			"device": map[string]interface{}{
				"identifiers": [...]string{fmt.Sprintf("Inverter_%s", uniq)},
				"name":        "Inverter",
			},
		})
		conn.publish(topic, string(json), true) // MQTT Discovery messages should be retained, but in dev it can become a pain
	}
	return nil
}

func (conn *Connection) InsertRecord(measurement map[string]interface{}) error {
	// make a copy
	m := make(map[string]interface{}, len(measurement))
	for k, v := range measurement {
		m[k] = v
	}
	// add LastTimestamp
	m["LastTimestamp"] = time.Now().UnixNano() / int64(time.Millisecond)
	json, _ := json.Marshal(m)
	conn.publish(conn.state, string(json), false) // state messages should not be retained
	return nil
}

func (conn *Connection) Subscribe(topic string, callback mqtt.MessageHandler) {
	conn.client.Subscribe(topic, 0, callback)
}
