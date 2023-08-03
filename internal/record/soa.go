package record

import "strings"

type SOARecord struct {
	Name       string `json:"name"`
	NameServer string `json:"nameserver"`
	Email      string `json:"email"`
	TTL        int    `json:"ttl"`
}

func (r *SOARecord) GetName() string {
	return r.Name + "."
}

func (r *SOARecord) GetNameServer() string {
	return r.NameServer + "."
}

func (r *SOARecord) GetMailBox() string {
	return strings.ReplaceAll(r.Email, "@", ".") + "."
}

func (r *SOARecord) GetTTL() uint32 {
	return uint32(r.TTL)
}
