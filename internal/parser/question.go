package parser

import (
	"fmt"
	"log"
	"time"

	"github.com/miekg/dns"
)

type QuestionParser struct {
	BaseDomain string
	Records    map[string]string
	SOADomain  string
	SOAMailBox string
}

func (p *QuestionParser) Parse(questions []dns.Question) (dns.RR, error) {
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
		case dns.TypeSOA:
			record := new(dns.SOA)
			record.Hdr = dns.RR_Header{Name: p.SOADomain + ".", Rrtype: dns.TypeSOA, Class: dns.ClassINET, Ttl: 0}
			record.Ns = q.Name
			record.Mbox = p.SOAMailBox + "."
			record.Serial = uint32(time.Now().Unix())
			record.Refresh = uint32(60)
			record.Retry = uint32(60)
			record.Expire = uint32(60)
			record.Minttl = uint32(60)

			return record, nil
		}

	}

	return nil, nil
}
