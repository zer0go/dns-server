package record

type ARecord struct {
	Name string `json:"name"`
	IP   string `json:"ip"`
	TTL  int    `json:"ttl"`
}

func (r *ARecord) GetName() string {
	return r.Name + "."
}

func (r *ARecord) GetTTL() uint32 {
	return uint32(r.TTL)
}
