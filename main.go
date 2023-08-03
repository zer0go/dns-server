package main

import (
	"flag"
	"fmt"
	"os"

	"github.com/miekg/dns"
	"github.com/rs/zerolog/log"
	"github.com/zer0go/dns-server/internal/config"
	"github.com/zer0go/dns-server/internal/handler"
	"github.com/zer0go/dns-server/internal/parser"
)

var (
	Version     = "development"
	DefaultAddr = "0.0.0.0:5353"
)

func main() {
	var verbose bool
	flag.BoolVar(&verbose, "v", false, "verbose output")
    flag.Parse()

	c := config.NewAppConfig()
	c.InitialLoadFromEnv()
	config.ConfigureLogger(verbose)

	log.Debug().
		Interface("config", c).
		Msg("configuration loaded")

	dnsHandler := handler.DNSHandler{
		Parser: parser.QuestionParser{
			RecordsA:   c.GetARecords(),
			RecordsSOA: c.GetSOARecords(),
		},
	}
	dns.HandleFunc(".", dnsHandler.Handle)

	addr := os.Getenv("ADDR")
	if addr == "" {
		addr = DefaultAddr
	}
	server := &dns.Server{
		Addr: addr,
		Net:  "udp",
	}
	fmt.Printf(
		"DNS Server %s | starting at udp://%s\n",
		Version,
		addr,
	)

	err := server.ListenAndServe()
	defer server.Shutdown()
	if err != nil {
		log.Error().Msgf("Failed to start server: %s\n ", err.Error())
	}
}
