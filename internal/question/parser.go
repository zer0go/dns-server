package question

import (
	"fmt"
	"log"

	"github.com/miekg/dns"
)

type Parser struct {
	BaseDomain string
	Records    map[string]string
}

func (p *Parser) Parse(questions []dns.Question) (dns.RR, error) {
	baseDomain := "."
	if p.BaseDomain != "" {
		baseDomain += p.BaseDomain + "."
	}
	for _, q := range questions {
		switch q.Qtype {
		case dns.TypeA:
			trimmedName := q.Name[:len(q.Name)-len(baseDomain)]
			log.Printf("Query for %s -> %s\n", q.Name, trimmedName)

			ip := p.Records[trimmedName]
			if ip == "" {
				break
			}

			log.Printf("Founded record: %s -> %s\n", trimmedName, ip)

			return dns.NewRR(fmt.Sprintf("%s A %s", q.Name, ip))
		}
	}

	return nil, nil
}