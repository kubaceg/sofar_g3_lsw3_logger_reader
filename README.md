# Sofar g3 LSW-3 logger reader

`sofar_g3_lsw3_logger_reader` is a go program to integrate certain Sofar solar PV inverters with MQTT (e.g. mosquito in HomeAssistant).
It works with Sofar inverters fitted with a LSW-3 WiFi data logger with serial number 23XXXXXXXX.
The program:
 - loops forever, polling the modbus port on the LSW-3 and writing the data to MQTT (all local, no dependency on SolarMan cloud),
 - reads its `config.yaml` from the current working directory,
 - writes logging to stderr,
 - supports MQTT Discovery.

`custom_components/sofar_g3_lsw3_logger_reader` is a HomeAssistant custom component to run the above program.

## Installation and setup

1. Download go 1.19 - 1.21
2. Clone this repo `git clone git@github.com:neilbacon/sofar_g3_lsw3_logger_reader.git`
3. Go into project directory `cd sofar_g3_lsw3_logger_reader`
6. Build program `make` (creates binaries in `custom_components/sofar_g3_lsw3_logger_reader`)
4. Go into directory `custom_components/sofar_g3_lsw3_logger_reader` and copy example config `cp config-example.yaml config.yaml` and edit `config.yaml` to suit your needs
7. Run `./sofar-x86` on x86 or `./sofar-arm` on arm (including Raspberry Pi)

To run it as part of Home Assistant see the [custom component README](custom_components/sofar_g3_lsw3_logger_reader/README.md).

## Output data format

### MQTT

#### Attribute Filtering
The LSW-3 provides a large number of attributes that you are most likely not interested in, so `config.yaml` allows you to filter them using either a white list or a black list used only if the white list is empty. The white list contains the attribute names (in full) to include and the black list contains regular expressions for attributes to exclude.

#### MQTT Discovery
On startup a message is sent (with `retain=true`) on a configuration topic for each data attribute. Home Assistant uses this information to configure an entity to manage the data attribute, removing the need for much manual configuration. The topic used is `{mqtt.discovery}/{attribute}/config` where `{mqtt.discovery}` comes from `config.yaml` and `{attribute}` is the attribute name. The JSON payload is described in the MQTT Discovery documentation.

`retain=true` causes a newly connected client (such as a restarting Home Assistant) to receive a copy of these messages. MQTT Discovery documentation suggests this is a good idea, but in development it can be a pain, leaving old messages hanging around indefinitely. You can use the `sh/cleanmqtt.sh` utility to clean up these messages. It depends on a MQTT installation and can be run in the HA OS's mosquitto container `addon_core_mosquitto`.

#### MQTT Data

All attributes are sent in a single message (with retain=false) with JSON payload to the topic `{mqtt.state}` specified in `config.yaml`.

### OTLP
Data can also be sent over OTLP protocol to a gRPC or http server. Typically, this would be received by the 
[OTel-Collector](https://opentelemetry.io/docs/collector/) for further export to any required platform. 

Metrics are all captured as gauges and recorded and exported at the same frequency that measurements are taken. 
Metric names follow the convention `sofar.logger.<fieldName>` by default. This can be updated in the configuration file.

## Origin
This is based on program written by @sigxcpu76 https://github.com/XtheOne/Inverter-Data-Logger/issues/37#issuecomment-1303091265.

## Contributing
Feel free if You want to extend this tool with new features. Just open issue or make PR.