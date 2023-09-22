# Sofar g3 LSW-3 logger reader Custom Component

This is a HomeAssistant custom component to run the [sofar_g3_lsw3_logger_reader](../../README.md) program.
It runs the program with:
  - the current working directory set to /config/custom_components/sofar_g3_lsw3_logger_reader
  - stdout/err going to files out/err.log (because I haven't succeeded in integrating it with homeassistant logging)

## Installation and setup

1. First get the go program `sofar_g3_lsw3_logger_reader` working on your dev machine, so that you have a working `config.yaml` and `sofar-x86` (x86) and `sofar-arm` (arm) binary executables (in this directory on your dev machine)
2. Copy this directory `custom_components/sofar_g3_lsw3_logger_reader` from your dev machine to homeassistant's `/config/custom_components` 
(creating `/config/custom_components/sofar_g3_lsw3_logger_reader`). 
You can use the `Samba share` or `Advanced SSH & Web Terminal` add-ons to do this.
3. Enable the custom component by adding a line `sofar_g3_lsw3_logger_reader:` to homeassistant's `/config/configuration.yaml`.
4. Do a full restart of homeassistant: `Developer Tools` > `YAML` > `CHECK CONFIGURATION` then `RESTART` > `Restart Home Assistant`
5. Check the content of homeassistant's `/config/custom_components/sofar_g3_lsw3_logger_reader/err.log`
6. Add the `Inverter` device to your dashhboard: `Settings` > `Devices & Services` > `Integration` > `MQTT` > `1 device` > `ADD TO DASHBOARD`

## To Do

1. Add to HACS
2. Get logging going to homeassistant's log.
3. Add support for "config flow" 
