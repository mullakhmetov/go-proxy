package main

import (
	"strings"
	"testing"
)

var (
	validConfigJSONStr = `
		{
			"services": {
				"app": {
					"bind": "9080",
					"backends": ["http://localhost:9081", "http://localhost:9082"],
					"balance": "roundrobin"
				},
				"db": {
					"bind": "5432",
					"backends": ["postgres://localhost:5432", "postgres://localhost:5433"],
					"balance": "random"
				}
			}
		}`
	invalidConfigJSONStr = "invalid json"
)

func TestParseConfig(t *testing.T) {
	validConfig := strings.NewReader(validConfigJSONStr)
	invalidConfig := strings.NewReader(invalidConfigJSONStr)
	c := new(Config)
	if err := c.ParseConfig(validConfig); err != nil {
		t.Error(err)
	}
	err := c.ParseConfig(invalidConfig)
	if _, ok := err.(*SettingsError); !ok {
		t.Errorf("ParseConfig should raise error if invalid config was passed")
	}
}
