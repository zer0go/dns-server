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

func (c *AppConfig) InitialLoadFromEnv() {
	c.Records = make(map[string]string)

	recordsJSON := os.Getenv("RECORDS_JSON")
	if recordsJSON != "" {
		c.LoadFromJSON([]byte(recordsJSON))
	}

	recordsJSONFile := os.Getenv("RECORDS_JSON_FILE")
	if recordsJSONFile != "" {
		c.LoadFromFile(recordsJSONFile)
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