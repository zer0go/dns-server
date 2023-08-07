package record

type NS struct {
	Name       string `json:"name"`
	NameServer string `json:"nameserver"`
	TTL        int    `json:"ttl"`
}

func (r *NS) GetName() string {
	return r.Name + "."
}

func (r *NS) GetNameServer() string {
	return r.NameServer + "."
}

func (r *NS) GetTTL() uint32 {
	return uint32(r.TTL)
}
