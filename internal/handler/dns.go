package handler

import (
	"github.com/miekg/dns"
	"github.com/zer0go/dns-server/internal/parser"
	"github.com/rs/zerolog/log"
)

type DNSHandler struct {
	Parser parser.QuestionParser
}

func (h *DNSHandler) Handle(w dns.ResponseWriter, r *dns.Msg) {
	m := new(dns.Msg)
	m.SetReply(r)
	m.Compress = false

	switch r.Opcode {
	case dns.OpcodeQuery:
		rr, err := h.Parser.Parse(m.Question)
		if err != nil {
			log.Error().Msgf("Parse error: %s\n", err.Error())
			break
		}
		if rr == nil {
			break
		}

		m.Answer = append(m.Answer, rr)
	}

	w.WriteMsg(m)
}
