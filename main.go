package main

import (
	"fmt"
	"log"
	"os"

	"github.com/miekg/dns"
	"github.com/zer0go/dns-server/internal/config"
	"github.com/zer0go/dns-server/internal/handler"
	"github.com/zer0go/dns-server/internal/parser"
)

var (
	Version     = "development"
	DefaultAddr = "0.0.0.0:5353"
)

func main() {
	config := config.NewAppConfig()
	config.InitialLoadFromEnv()

	reloadConfig := os.Getenv("RELOAD_CONFIG") == "1"
	baseDomain := os.Getenv("BASE_DOMAIN")
	dnsHandler := handler.DNSHandler{
		Parser: parser.QuestionParser{
			BaseDomain: baseDomain,
			Records:    config.Records,
			SOADomain:  os.Getenv("SOA_DOMAIN"),
			SOAMailBox: os.Getenv("SOA_MAILBOX"),
		},
	}
	dns.HandleFunc(".", func(w dns.ResponseWriter, r *dns.Msg) {
		if reloadConfig {
			config.InitialLoadFromEnv()
		}
		dnsHandler.Handle(w, r)
	})

	addr := os.Getenv("ADDR")
	if addr == "" {
		addr = DefaultAddr
	}
	server := &dns.Server{
		Addr: addr,
		Net:  "udp",
	}
	fmt.Printf(
		"DNS Server %s | starting at udp://%s | base domain '%s' | reload config '%t' \n",
		Version,
		addr,
		baseDomain,
		reloadConfig,
	)

	err := server.ListenAndServe()
	defer server.Shutdown()
	if err != nil {
		log.Fatalf("Failed to start server: %s\n ", err.Error())
	}
}
