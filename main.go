package main

import (
	"fmt"
	"log"
	"os"
	"strconv"

	"github.com/miekg/dns"
	"github.com/zer0go/dns-server/internal/config"
	"github.com/zer0go/dns-server/internal/handler"
	"github.com/zer0go/dns-server/internal/question"
)

var Version = "development"
var DefaultAddr = "0.0.0.0:15353"

func main() {
	config := config.NewAppConfig()
	recordsJSON := os.Getenv("RECORDS_JSON")
	if recordsJSON != "" {
		config.LoadFromJSON([]byte(recordsJSON))
	}
	recordsJSONFile := os.Getenv("RECORDS_JSON_FILE")
	if recordsJSONFile != "" {
		config.LoadFromFile(recordsJSONFile)
	}

	baseDomain := os.Getenv("BASE_DOMAIN")
	dnsHandler := handler.DNSHandler{
		Parser: question.Parser{
			BaseDomain: baseDomain,
			Records: config.Records,
		},
	}
	dns.HandleFunc(baseDomain + ".", dnsHandler.Handle)

	addr := os.Getenv("ADDR")
	if addr == "" {
		addr = DefaultAddr
	}
	server := &dns.Server{
		Addr: addr,
		Net:  "udp",
	}
	fmt.Printf("DNS Server %s | starting at udp://%s | base domain '%s' \n", Version, addr, baseDomain)
	err := server.ListenAndServe()
	defer server.Shutdown()
	if err != nil {
		log.Fatalf("Failed to start server: %s\n ", err.Error())
	}
}
