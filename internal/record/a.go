package record

type ARecord struct {
	Name string `json:"name"`
	IP   string `json:"ip"`
}

func (r *ARecord) GetName() string {
	return r.Name + "."
}
