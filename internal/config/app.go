package config

import (
	"fmt"
	"os"
	"path/filepath"
	"strconv"
	"time"

	"github.com/caarlos0/env/v9"
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
	"github.com/rs/zerolog/pkgerrors"
)

type AppConfig struct {
	Addr                  string `env:"ADDR" envDefault:"0.0.0.0:5353"`
	LogFormat             string `env:"LOG_FORMAT"`
	RecordsFile           string `env:"RECORDS_FILE"`
	ReloadRecordsInterval int    `env:"RELOAD_RECORDS_INTERVAL" envDefault:"0"`
}

func NewAppConfig() *AppConfig {
	return new(AppConfig)
}

func NewAppConfigFromEnv() *AppConfig {
	c := NewAppConfig()
	if err := env.Parse(c); err != nil {
		fmt.Printf("env parse error: %+v\n", err)
		return nil
	}

	return c
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
