package main

import (
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"os"

	"github.com/frodenas/helm-osb/broker"
	"github.com/frodenas/helm-osb/helm"
)

type Config struct {
	LogLevel     string        `json:"log_level"`
	BrokerConfig broker.Config `json:"broker"`
	HelmConfig   helm.Config   `json:"helm"`
}

func LoadConfig(configFilePath string) (config *Config, err error) {
	if configFilePath == "" {
		return config, errors.New("Must provide a non-empty configuration file")
	}

	configFile, err := os.Open(configFilePath)
	if err != nil {
		return config, err
	}
	defer configFile.Close()

	bytes, err := ioutil.ReadAll(configFile)
	if err != nil {
		return config, err
	}

	if err = json.Unmarshal(bytes, &config); err != nil {
		return config, err
	}

	if err = config.Validate(); err != nil {
		return config, fmt.Errorf("Validating configuration file contents: %s", err)
	}

	return config, nil
}

func (c Config) Validate() error {
	if c.LogLevel == "" {
		return errors.New("Must provide a non-empty Log Level")
	}

	if err := c.BrokerConfig.Validate(); err != nil {
		return fmt.Errorf("Validating Broker configuration: %s", err)
	}

	if err := c.HelmConfig.Validate(); err != nil {
		return fmt.Errorf("Validating Helm configuration: %s", err)
	}

	return nil
}
