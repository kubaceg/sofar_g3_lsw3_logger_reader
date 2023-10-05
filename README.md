# Sofar g3 LSW-3 logger reader
Tool written in GO for reading metrics from Sofar LSW-3 and writing results into MQTT topics. 
Program queries logger modbus port in infinite loop and sends data into MQTT topics (e.g. mosquito in HomeAssistant).

## Installation and setup
1. Download go 1.19
2. Clone this repo `git clone git@github.com:kubaceg/sofar_g3_lsw3_logger_reader.git`
3. Go into project directory `cd sofar_g3_lsw3_logger_reader`
4. Copy example config `cp config-example.yaml config.yaml`
5. Edit `config.yaml` in Your favorite editor, fill all required stuff
6. Build program `make build` or build for ARM machines e.g. raspberryPi `make build-arm`
7. Run `./sofar` or `sofar-arm`

## Output data format
### MQTT
Data will be sent into MQTT topic with name `{mqttPrefix}/{fieldName}` where:
* mqttPrefix is value defined in `config.yaml` e.g. `/sensors/energy/inverter`
* fieldName is measurement name, all available measurements are described in `adapters/devices/sofar/sofar_protocol.go`, e.g. `PV_Generation_Today`

Full topic name for given example values is `/sensors/energy/inverter/PV_Generation_Today`.
Additional field is `All` which contains all measurements and their values marshalled into one json.

### OTLP
Data can also be sent over OTLP protocol to a gRPC or http server. Typically, this would be received by the 
[OTel-Collector](https://opentelemetry.io/docs/collector/) for further export to any required platform. 

Metrics are all captured as gauges and recorded and exported at the same frequency that measurements are taken. 
Metric names follow the convention `sofar.logger.<fieldName>` by default. This can be updated in the configuration file.

### Grafana dashboard
You can monitor Your solar instalation using [grafana dashboard](grafana)

## Origin
This is based on program written by @sigxcpu76 https://github.com/XtheOne/Inverter-Data-Logger/issues/37#issuecomment-1303091265.

## Contributing
Feel free if You want to extend this tool with new features. Just open issue or make PR.
