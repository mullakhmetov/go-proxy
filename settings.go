package main

import (
	"encoding/json"
	"io"
)

type SettingsError struct {
	Detail string
	Err    error
}

func (e *SettingsError) Error() string {
	return e.Detail + e.Err.Error()
}

// used to parse settings file
type Config struct {
	Services map[string]ServiceConfig `json:"services"`
}

type ServiceConfig struct {
	Backends []string `json:"backends"`
	Balance  string   `json:"balance"`
	Bind     string   `json:bind`
}

// open config file by path and decode it to Config struct
func (c *Config) ParseConfig(configFile io.Reader) error {
	jsonParser := json.NewDecoder(configFile)
	if err := jsonParser.Decode(c); err != nil {
		return &SettingsError{"Failed to parse config file", err}
	}
	return nil
}

func NewConfig(configFile io.Reader) (*Config, error) {
	c := new(Config)
	if err := c.ParseConfig(configFile); err != nil {
		return nil, err
	}
	return c, nil
}
