package record

type A struct {
	Name string `json:"name"`
	IP   string `json:"ip"`
	TTL  int    `json:"ttl"`
}

func (r *A) GetName() string {
	return r.Name + "."
}

func (r *A) GetTTL() uint32 {
	return uint32(r.TTL)
}
