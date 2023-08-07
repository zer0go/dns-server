package main

import (
	"flag"
	"fmt"
	"time"

	"github.com/miekg/dns"
	"github.com/rs/zerolog/log"
	"github.com/zer0go/dns-server/internal/config"
	"github.com/zer0go/dns-server/internal/handler"
	"github.com/zer0go/dns-server/internal/question"
	"github.com/zer0go/dns-server/internal/record"
)

var (
	Version = "development"
)

func reloadRecords(ds *record.DataSource, i int) {
	if i == 0 {
		return
	}

	for {
		time.Sleep(time.Duration(i) * time.Second)

		lastModified := ds.GetLastModified()
		lastLoaded := ds.GetLastLoaded()

		if lastModified.Before(lastLoaded) {
			continue
		}

		ds.Load()

		log.Debug().
			Interface("last_loaded", lastLoaded).
			Interface("last_modified", lastModified).
			Interface("records", ds.GetRecords()).
			Msg("records updated")
	}
}

func main() {
	var verbose bool
	flag.BoolVar(&verbose, "v", false, "verbose output")
	flag.Parse()

	c := config.NewAppConfigFromEnv()
	c.ConfigureLogger(verbose)

	log.Debug().
		Interface("config", c).
		Msg("configuration loaded")

	dataSource := record.NewDataSource(c.RecordsFile)
	dataSource.Load()

	log.Debug().
		Interface("records", dataSource.GetRecords()).
		Msg("records loaded")

	go reloadRecords(dataSource, c.ReloadRecordsInterval)

	dnsHandler := handler.DNSHandler{
		Resolver: question.NewResolver(
			record.NewDAO(dataSource),
		),
	}
	dns.HandleFunc(".", dnsHandler.Handle)

	server := &dns.Server{
		Addr: c.Addr,
		Net:  "udp",
	}
	fmt.Printf(
		"DNS Server %s | starting at udp://%s\n",
		Version,
		c.Addr,
	)

	err := server.ListenAndServe()
	defer server.Shutdown()
	if err != nil {
		log.Error().Err(err).Msgf("failed to start server")
	}
}
