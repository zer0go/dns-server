package parser

import (
	"fmt"
	"strings"
	"time"

	"github.com/miekg/dns"
	"github.com/zer0go/dns-server/internal/record"
	"github.com/rs/zerolog/log"
)

type QuestionParser struct {
	RecordsA   []record.ARecord
	RecordsSOA []record.SOARecord
}

func (p *QuestionParser) Parse(questions []dns.Question) (dns.RR, error) {
	for _, q := range questions {
		name := strings.TrimSuffix(q.Name, ".")
		switch q.Qtype {
		case dns.TypeA:
			log.Info().Msgf("Query for [A] %s", q.Name)
			for _, r := range p.RecordsA {
				if r.Name == name {
					log.Info().Msgf("Founded record: [A] %s -> %s", r.Name, r.IP)

					return dns.NewRR(fmt.Sprintf("%s A %s", r.GetName(), r.IP))
				}
			}
		case dns.TypeSOA:
			log.Info().Msgf("Query for [SOA] %s", q.Name)
			for _, r := range p.RecordsSOA {
				if r.Name == name {
					mailBox := r.GetMailBox()
					log.Info().Msgf("Founded record: [SOA] %s (Ns: %s, Mbox: %s)", r.Name, r.NameServer, mailBox)

					result := new(dns.SOA)
					result.Hdr = dns.RR_Header{Name: r.GetName(), Rrtype: dns.TypeSOA, Class: dns.ClassINET, Ttl: 0}
					result.Ns = r.GetNameServer()
					result.Mbox = mailBox
					result.Serial = uint32(time.Now().Unix())
					result.Refresh = uint32(60)
					result.Retry = uint32(60)
					result.Expire = uint32(60)
					result.Minttl = uint32(60)

					return result, nil
				}
			}
		}
	}

	return nil, nil
}
