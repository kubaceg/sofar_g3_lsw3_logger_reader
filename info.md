# Integration for Sofar Solar PV Inverters with LSW-3 WiFi data logger with SN 23XXXXXXX

This integration regularly polls the LSW-3 and sends inverter data on a MQTT topic.
It supports MQTT discovery.
It is implemented by a binary executable `sofar-arm` (on arm) or `sofar-x86` (on x86), which is kicked off when the integration starts and runs forever. 

## Manual Installation

1. First get the go program `sofar_g3_lsw3_logger_reader` working on your dev machine, so that you have a working `config.yaml` and `sofar-x86` (x86) and `sofar-arm` (arm) binary executables.
2. Copy this directory `custom_components/sofar_g3_lsw3_logger_reader` to homeassistant's `/config/custom_components` 
(creating `/config/custom_components/sofar_g3_lsw3_logger_reader`). 
You can use the `Samba share` or `Advanced SSH & Web Terminal` add-ons to do this.
3. Copy the working `config.yaml` (from step 1) to the same directory.
4. Copy the executables `sofar` (x86) and `sofar-arm` (arm) (from step 1) to this same directory.
5. Enable the custom component by adding a line `sofar_g3_lsw3_logger_reader:` to homeassistant's `/config/configuration.yaml`.
6. Do a full restart of homeassistant: `Developer Tools` > `YAML` > `CHECK CONFIGURATION` then `RESTART` > `Restart Home Assistant`
7. Check the content of `err.log` in this same directory.
8. Add the `Inverter` device to your dashhboard: `Settings` > `Devices & Services` > `Integration` > `MQTT` > `1 device` > `ADD TO DASHBOARD`

## HACS Installation

1. Install using HACS
2. In directory `/config/custom_components/sofar_g3_lsw3_logger_reader` copy `config-example.yaml` to `config.yaml` and edit to match your requirements.
3. Restart Home Assistant to run the integration.

## To Do

1. Get logging going to homeassistant's log.
2. Add support for "config flow" 
