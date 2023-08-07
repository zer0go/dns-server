package question

import (
	"net"
	"strconv"
	"strings"
	"time"

	"github.com/miekg/dns"
	"github.com/rs/zerolog/log"
	"github.com/zer0go/dns-server/internal/record"
)

type Resolver struct {
	dao record.DAO
}

func NewResolver(dao *record.DAO) *Resolver {
	return &Resolver{
		dao: *dao,
	}
}

func (r *Resolver) Resolve(questions []dns.Question) (rrs []dns.RR, err error) {
	for _, q := range questions {
		simplifiedName := simplifyName(q.Name)
		log.Info().
			Str("query_name", q.Name).
			Str("simplified_name", simplifiedName).
			Str("type", dns.Type(q.Qtype).String()).
			Msgf("Incoming question: '%s'", q.String())

		switch q.Qtype {
		case dns.TypeA:
			for _, r := range r.dao.GetARecordsByName(simplifiedName) {
				name := r.GetName()
				ttl := r.GetTTL()
				log.Info().
					Str("name", name).
					Str("type", dns.Type(q.Qtype).String()).
					Str("ip", r.IP).
					Uint32("ttl", ttl).
					Msgf("Founded record: %s", r.Name)

				result := new(dns.A)
				result.Hdr = dns.RR_Header{Name: name, Rrtype: q.Qtype, Class: dns.ClassINET, Ttl: ttl}
				result.A = net.ParseIP(r.IP)

				rrs = append(rrs, result)
			}
		case dns.TypeSOA:
			for _, r := range r.dao.GetSOARecordsByName(simplifiedName) {
				name := r.GetName()
				nameServer := r.GetNameServer()
				mailBox := r.GetMailBox()
				serial, _ := strconv.Atoi(time.Now().Format("20060102") + "00")
				ttl := r.GetTTL()
				log.Info().
					Str("name", name).
					Str("type", dns.Type(q.Qtype).String()).
					Str("name_server", nameServer).
					Str("email", r.Email).
					Str("mail_box", mailBox).
					Int("serial", serial).
					Uint32("ttl", ttl).
					Msgf("Founded record: %s", r.Name)

				result := new(dns.SOA)
				result.Hdr = dns.RR_Header{Name: name, Rrtype: q.Qtype, Class: dns.ClassINET, Ttl: ttl}
				result.Ns = nameServer
				result.Mbox = mailBox
				result.Serial = uint32(serial)
				result.Refresh = uint32(3600)
				result.Retry = uint32(1800)
				result.Expire = uint32(2592000)
				result.Minttl = uint32(3600)

				rrs = append(rrs, result)
			}
		case dns.TypeCNAME:
			for _, r := range r.dao.GetCNAMERecordsByName(simplifiedName) {
				name := q.Name
				target := r.GetTarget()
				ttl := r.GetTTL()
				log.Info().
					Str("name", name).
					Str("type", dns.Type(q.Qtype).String()).
					Str("target", target).
					Uint32("ttl", ttl).
					Msgf("Founded record: %s", r.Name)

				result := new(dns.CNAME)
				result.Hdr = dns.RR_Header{Name: name, Rrtype: q.Qtype, Class: dns.ClassINET, Ttl: ttl}
				result.Target = target
				
				rrs = append(rrs, result)
			}
		case dns.TypeNS:
			for _, r := range r.dao.GetNSRecordsByName(simplifiedName) {
				name := r.GetName()
				nameServer := r.GetNameServer()
				ttl := r.GetTTL()
				log.Info().
					Str("name", name).
					Str("type", dns.Type(q.Qtype).String()).
					Str("name_server", nameServer).
					Uint32("ttl", ttl).
					Msgf("Founded record: %s", r.Name)

				result := new(dns.NS)
				result.Hdr = dns.RR_Header{Name: name, Rrtype: q.Qtype, Class: dns.ClassINET, Ttl: ttl}
				result.Ns = nameServer
				
				rrs = append(rrs, result)
			}
		}
	}

	return
}

func simplifyName(name string) string {
	return strings.TrimSuffix(strings.ToLower(name), ".")
}
