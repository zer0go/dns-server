package config

import (
	"encoding/json"
	"os"
	"log"
)

type AppConfig struct {
	Records map[string]string
}

func NewAppConfig() *AppConfig {
	return &AppConfig{
		Records: make(map[string]string),
	}
}

func (c *AppConfig) LoadFromJSON(b []byte) {
	r := make(map[string]string)
	_ = json.Unmarshal(b, &r)

	for k, v := range r {
		c.Records[k] = v
	}
}

func (c *AppConfig) LoadFromFile(f string) {
	b, err := os.ReadFile(f)
    if err != nil {
        log.Print(err)
    }

    c.LoadFromJSON(b)
}