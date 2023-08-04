package main

import (
	"flag"
	"fmt"
	//"time"

	"github.com/miekg/dns"
	"github.com/rs/zerolog/log"
	"github.com/zer0go/dns-server/internal/config"
	"github.com/zer0go/dns-server/internal/handler"
	"github.com/zer0go/dns-server/internal/parser"
	//"github.com/zer0go/dns-server/internal/record"
)

var (
	Version     = "development"
)

func main() {
	var verbose bool
	flag.BoolVar(&verbose, "v", false, "verbose output")
	flag.Parse()

	c := config.NewAppConfigFromEnv()
	c.ConfigureLogger(verbose)

	log.Debug().
		Interface("config", c).
		Msg("configuration loaded")

	dnsHandler := handler.DNSHandler{
		Parser: parser.QuestionParser{
			RecordsA:   c.GetARecords(),
			RecordsSOA: c.GetSOARecords(),
		},
	}
	
// 	go func() {
//     for {
//       time.Sleep(10 * time.Second)
      
//       dnsHandler = handler.DNSHandler{
//     		Parser: parser.QuestionParser{
//     			RecordsA:   []record.ARecord{},
//     			RecordsSOA: []record.SOARecord{},
//     		},
//     	}
      
//       log.Debug().
// 		    Interface("config", c).
// 		    Msg("configuration updated")
//     }
//   }()
  
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
