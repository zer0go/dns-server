package config

import (
	"encoding/json"
	"os"
	"fmt"
	"path/filepath"
	"strconv"
	"time"

	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
	"github.com/rs/zerolog/pkgerrors"
	"github.com/zer0go/dns-server/internal/record"
	"github.com/caarlos0/env/v9"
)

type RecordsConfig struct {
	A   []record.ARecord   `json:"A"`
	SOA []record.SOARecord `json:"SOA"`
}

type DomainRecordsConfig struct {
	Domain string            `json:"domain"`
	TTL    int               `json:"ttl"`
	Hosts  map[string]string `json:"hosts"`
}

type AppConfig struct {
  Addr          string `env:"ADDR" envDefault:"0.0.0.0:5353"`
	LogFormat     string `env:"LOG_FORMAT"`
  ConfigJSON    string `env:"CONFIG_JSON"`
  ConfigFile    string `env:"CONFIG_FILE"`
	Records       RecordsConfig         `json:"records"`
	DomainRecords []DomainRecordsConfig `json:"domain_records"`
}

func NewAppConfig() *AppConfig {
	return &AppConfig{
		Records: RecordsConfig{
			A:   []record.ARecord{},
			SOA: []record.SOARecord{},
		},
		DomainRecords: []DomainRecordsConfig{},
	}
}

func NewAppConfigFromEnv() *AppConfig {
	c := NewAppConfig()
  if err := env.Parse(c); err != nil {
		fmt.Printf("env parse error: %+v\n", err)
		return nil
	}

	if c.ConfigJSON != "" {
		c.LoadFromJSON([]byte(c.ConfigJSON))
	}

	if c.ConfigFile != "" {
		c.LoadFromFile(c.ConfigFile)
	}

	return c
}

func (c *AppConfig) LoadFromJSON(b []byte) {
	loadedConfig := NewAppConfig()
	err := json.Unmarshal(b, &loadedConfig)
	if err != nil {
		log.Error().Err(err).Msg("error loading from json")

		return
	}

	c.Records.A = append(c.Records.A, loadedConfig.Records.A...)
	c.Records.SOA = append(c.Records.SOA, loadedConfig.Records.SOA...)
	c.DomainRecords = append(c.DomainRecords, loadedConfig.DomainRecords...)
}

func (c *AppConfig) LoadFromFile(f string) {
	b, err := os.ReadFile(f)
	if err != nil {
		log.Error().Err(err).Msg("error loading from file")

		return
	}

	c.LoadFromJSON(b)
}

func (c *AppConfig) GetARecords() []record.ARecord {
	records := c.Records.A

	for _, r := range c.DomainRecords {
		for hostName, ip := range r.Hosts {
			records = append(records, record.ARecord{
				Name: hostName + "." + r.Domain,
				IP:   ip,
				TTL:  r.TTL,
			})
		}
	}

	return records
}

func (c *AppConfig) GetSOARecords() []record.SOARecord {
	return c.Records.SOA
}

func (c *AppConfig) ConfigureLogger(verbose bool) {
	if verbose {
		zerolog.SetGlobalLevel(zerolog.DebugLevel)
	} else {
		zerolog.SetGlobalLevel(zerolog.InfoLevel)
	}

	zerolog.TimeFieldFormat = time.RFC3339
	zerolog.ErrorStackMarshaler = pkgerrors.MarshalStack
	zerolog.CallerMarshalFunc = func(pc uintptr, file string, line int) string {
		dir := filepath.Dir(file)
		parent := filepath.Base(dir)
		return parent + "/" + filepath.Base(file) + ":" + strconv.Itoa(line)
	}

	log.Logger = log.
		With().
		Timestamp().
		Stack().
		Caller().
		Logger()

	if c.LogFormat != "json" {
		log.Logger = log.Logger.Output(zerolog.ConsoleWriter{Out: os.Stderr, TimeFormat: zerolog.TimeFieldFormat})
	}
}
