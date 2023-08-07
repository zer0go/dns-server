package handler

import (
	"github.com/miekg/dns"
	"github.com/rs/zerolog/log"
	"github.com/zer0go/dns-server/internal/question"
)

type DNSHandler struct {
	Resolver *question.Resolver
}

func (h *DNSHandler) Handle(w dns.ResponseWriter, r *dns.Msg) {
	m := new(dns.Msg)
	m.SetReply(r)
	m.Compress = false

	switch r.Opcode {
	case dns.OpcodeQuery:
		rr, err := h.Resolver.Resolve(m.Question)
		if err != nil {
			log.Error().Msgf("Parse error: %s\n", err.Error())
			break
		}
		if len(rr) == 0 {
			break
		}

		m.Answer = append(m.Answer, rr...)
	}

	w.WriteMsg(m)
}
