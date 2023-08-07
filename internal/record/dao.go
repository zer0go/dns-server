package record

import "strings"

type DAO struct {
	dataSource DataSource
}

func NewDAO(dataSource *DataSource) *DAO {
	return &DAO{
		dataSource: *dataSource,
	}
}

func (r *DAO) GetARecordsByName(name string) []A {
	var records []A
	for _, r := range r.dataSource.records.A {
		if r.Name == name {
			records = append(records, r)
		}
	}

	return records
}

func (r *DAO) GetSOARecordsByName(name string) []SOA {
	var records []SOA
	for _, r := range r.dataSource.records.SOA {
		if r.Name == name {
			records = append(records, r)
		}
	}

	return records
}

func (r *DAO) GetCNAMERecordsByName(name string) []CNAME {
	var records []CNAME
	for _, r := range r.dataSource.records.CNAME {
		if r.Name == name || strings.HasSuffix(getWithoutWildCardCharacter(name), getWithoutWildCardCharacter(r.Name)) {
			records = append(records, r)
		}
	}

	return records
}

func (r *DAO) GetNSRecordsByName(name string) []NS {
	var records []NS
	for _, r := range r.dataSource.records.NS {
		if r.Name == name {
			records = append(records, r)
		}
	}

	return records
}

func getWithoutWildCardCharacter(name string) string {
	return strings.Replace(name, "*", "", 1)
}
