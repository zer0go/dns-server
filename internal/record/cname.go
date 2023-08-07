package record

type CNAME struct {
	Name   string `json:"name"`
	Target string `json:"target"`
	TTL    int    `json:"ttl"`
}

func (r *CNAME) GetName() string {
	return r.Name + "."
}

func (r *CNAME) GetTarget() string {
	return r.Target + "."
}

func (r *CNAME) GetTTL() uint32 {
	return uint32(r.TTL)
}
