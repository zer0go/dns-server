package record

import "strings"

type SOA struct {
	Name       string `json:"name"`
	NameServer string `json:"nameserver"`
	Email      string `json:"email"`
	TTL        int    `json:"ttl"`
}

func (r *SOA) GetName() string {
	return r.Name + "."
}

func (r *SOA) GetNameServer() string {
	return r.NameServer + "."
}

func (r *SOA) GetMailBox() string {
	return strings.ReplaceAll(r.Email, "@", ".") + "."
}

func (r *SOA) GetTTL() uint32 {
	return uint32(r.TTL)
}
