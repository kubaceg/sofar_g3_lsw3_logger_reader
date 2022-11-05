package main

import (
	"errors"
	"os"

	"github.com/kubaceg/sofar_g3_lsw3_logger_reader/adapters/databases/mosquitto"
	"gopkg.in/yaml.v2"
)

type Config struct {
	Inverter struct {
		Port         string `yaml:"port"`
		LoggerSerial string `yaml:"loggerSerial"`
		ReadInterval int    `default:"60" yaml:"readInterval"`
	} `yaml:"inverter"`
	Mqtt mosquitto.MqttConfig `yaml:"mqtt"`
}

func (c *Config) validate() error {
	if c.Inverter.Port == "" {
		return errors.New("missing required inverter.port config")
	}

	if c.Inverter.LoggerSerial == "" {
		return errors.New("missing required inverter.loggerSerial config")
	}

	return nil
}

func NewConfig(configPath string) (*Config, error) {
	config := &Config{}
	file, err := os.Open(configPath)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	d := yaml.NewDecoder(file)

	if err := d.Decode(&config); err != nil {
		return nil, err
	}

	if err := config.validate(); err != nil {
		return nil, err
	}

	return config, nil
}
