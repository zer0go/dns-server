package config

import (
	"encoding/json"
	"os"
	"path/filepath"
	"strconv"
	"time"

	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
	"github.com/rs/zerolog/pkgerrors"
	"github.com/zer0go/dns-server/internal/record"
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

	configJSON := os.Getenv("CONFIG_JSON")
	if configJSON != "" {
		c.LoadFromJSON([]byte(configJSON))
	}

	configJSONFile := os.Getenv("CONFIG_JSON_FILE")
	if configJSONFile != "" {
		c.LoadFromFile(configJSONFile)
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

func ConfigureLogger(verbose bool) {
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

	if os.Getenv("LOG_FORMAT") != "json" {
		log.Logger = log.Logger.Output(zerolog.ConsoleWriter{Out: os.Stderr, TimeFormat: zerolog.TimeFieldFormat})
	}
}
