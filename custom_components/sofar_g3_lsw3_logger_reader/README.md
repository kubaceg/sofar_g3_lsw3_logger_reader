# Custom Component to run sofar_g3_lsw3_logger_reader

## Introduction

`sofar_g3_lsw3_logger_reader` is a program written in go that:
 - loops forever, polling the LSW3 data logger and writing the data to a MQTT topic
 - reads its `config.yaml` from the current working directory
 - writes logging to stderr
 - supports MQTT Discovery

 This component runs `sofar_g3_lsw3_logger_reader` with:
  - the current working directory set to /config/custom_components/sofar_g3_lsw3_logger_reader
  - stdout/err going to files out/err.log (because I haven't succeeded in integrating it with homeassistant logging)

## Installation

1. First get `sofar_g3_lsw3_logger_reader` working on your dev machine, so that you have a working `config.yaml` and `sofar` (x86) and `sofar-arm` (arm) binary executables.
2. Copy this directory `custom_components/sofar_g3_lsw3_logger_reader` to homeassistant's `/config/custom_components` 
(creating `/config/custom_components/sofar_g3_lsw3_logger_reader`). 
You can use the `Samba share` or `Advanced SSH & Web Terminal` add-ons to do this.
3. Copy the working `config.yaml` (from step 1) to the same directory.
4. Copy the `sofar` (x86) or `sofar-arm` (arm) binary executable (from step 1) to `sofar` in this same directory.
5. Enable the custom component by adding a line `sofar_g3_lsw3_logger_reader:` to homeassistant's `/config/configuration.yaml`.
6. Do a full restart of homeassistant: `Developer Tools` > `YAML` > `CHECK CONFIGURATION` then `RESTART` > `Restart Home Assistant`
7. Check the content of `err.log` in this same directory.
8. Add the `Inverter` device to your dashhboard: `Settings` > `Devices & Services` > `Integration` > `MQTT` > `1 device` > `ADD TO DASHBOARD`

## To Do

1. Add to HACS
2. Get logging going to homeassistant's log.
3. Add support for "config flow" 
