package main

import (
	"errors"
	"fmt"
	"github.com/BurntSushi/toml"
)

type Sensor struct {
	ID       int
	URL      string
	JSONPath string
}
type SensorConfig struct {
	Sensors []Sensor `toml:"sensor"`
}

func loadSensorConfig(filename string) (*SensorConfig, error) {
	var config SensorConfig
	if _, err := toml.DecodeFile(filename, &config); err != nil {
		return nil, err
	}
	if err := validateConfig(config); err != nil {
		return nil, err
	}
	return &config, nil
}

func validateConfig(config SensorConfig) error {
	errorList := []error{}
	if len(config.Sensors) == 0 {
		errorList = append(errorList, errors.New(fmt.Sprintf("No sensor configs found")))
	}
	ids := make(map[int]interface{})
	for _, sensor := range config.Sensors {
		if _, missing := ids[sensor.ID]; missing {
			errorList = append(errorList, errors.New(fmt.Sprintf("Sensor %d: Duplicate ID", sensor.ID)))
		}
		ids[sensor.ID] = nil

		if len(sensor.URL) == 0 {
			errorList = append(errorList, errors.New(fmt.Sprintf("Sensor %d: Empty URL", sensor.ID)))
		}
		if len(sensor.JSONPath) == 0 {
			errorList = append(errorList, errors.New(fmt.Sprintf("Sensor %d: Empty JSONPath", sensor.ID)))
		}
	}
	if len(errorList) > 0 {
		return errors.Join(errorList...)
	}
	return nil
}
