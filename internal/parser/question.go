package parser

import (
	"net"
	"strconv"
	"strings"
	"time"

	"github.com/miekg/dns"
	"github.com/rs/zerolog/log"
	"github.com/zer0go/dns-server/internal/record"
)

type QuestionParser struct {
	RecordsA   []record.ARecord
	RecordsSOA []record.SOARecord
}

func (p *QuestionParser) Parse(questions []dns.Question) (rrs []dns.RR, err error) {
	for _, q := range questions {
		simplifiedName := simplifyName(q.Name)
		log.Info().
			Str("query_name", q.Name).
			Str("simplified_name", simplifiedName).
			Str("type", dns.Type(q.Qtype).String()).
			Msgf("Incoming question: '%s'", q.String())

		switch q.Qtype {
		case dns.TypeA:
			for _, r := range p.RecordsA {
				if r.Name == simplifiedName {
					name := r.GetName()
					log.Info().
						Str("name", name).
						Str("type", dns.Type(q.Qtype).String()).
						Str("ip", r.IP).
						Msgf("Founded record: %s", r.Name)

					result := new(dns.A)
					result.Hdr = dns.RR_Header{Name: name, Rrtype: dns.TypeA, Class: dns.ClassINET, Ttl: r.GetTTL()}
					result.A = net.ParseIP(r.IP)

					rrs = append(rrs, result)
				}
			}
		case dns.TypeSOA:
			for _, r := range p.RecordsSOA {
				if r.Name == simplifiedName {
					name := r.GetName()
					nameServer := r.GetNameServer()
					mailBox := r.GetMailBox()
					ttl := r.GetTTL()
					serial, _ := strconv.Atoi(time.Now().Format("20060102") + "00")
					log.Info().
						Str("name", name).
						Str("type", dns.Type(q.Qtype).String()).
						Str("name_server", nameServer).
						Str("email", r.Email).
						Str("mailBox", mailBox).
						Uint32("ttl", ttl).
						Int("serial", serial).
						Msgf("Founded record: %s", r.Name)

					result := new(dns.SOA)
					result.Hdr = dns.RR_Header{Name: name, Rrtype: dns.TypeSOA, Class: dns.ClassINET, Ttl: ttl}
					result.Ns = nameServer
					result.Mbox = mailBox
					result.Serial = uint32(serial)
					result.Refresh = uint32(3600)
					result.Retry = uint32(1800)
					result.Expire = uint32(2592000)
					result.Minttl = uint32(3600)

					rrs = append(rrs, result)
				}
			}
		}
	}
	
	return
}

func simplifyName(name string) string {
	return strings.TrimSuffix(strings.ToLower(name), ".")
}
