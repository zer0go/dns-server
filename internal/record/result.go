package record

type DomainRecordsResult struct {
	Domain string            `json:"domain"`
	TTL    int               `json:"ttl"`
	Hosts  map[string]string `json:"hosts"`
}

type RecordsResult struct {
	A     []A     `json:"A"`
	SOA   []SOA   `json:"SOA"`
	CNAME []CNAME `json:"CNAME"`
	NS    []NS    `json:"NS"`
}

func NewRecordsResult() *RecordsResult {
	return &RecordsResult{
		A:     []A{},
		SOA:   []SOA{},
		CNAME: []CNAME{},
		NS:    []NS{},
	}
}

func (r *RecordsResult) InitializeWithDataSourceEntity(entity DataSourceEntity) {
	r.A = entity.GetMergedARecords()
	r.SOA = entity.Records.SOA
	r.CNAME = entity.Records.CNAME
	r.NS = entity.Records.NS
}