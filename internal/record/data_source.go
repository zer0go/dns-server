package record

import (
	"encoding/json"
	"os"
	"time"

	"github.com/rs/zerolog/log"
)

type DataSourceEntity struct {
	Records       RecordsResult         `json:"records"`
	DomainRecords []DomainRecordsResult `json:"domain_records"`
}

func NewDataSourceEntity() *DataSourceEntity {
	return &DataSourceEntity{
		Records: *NewRecordsResult(),
		DomainRecords: []DomainRecordsResult{},
	}
}

func (d *DataSourceEntity) GetMergedARecords() []A {
	records := d.Records.A

	for _, record := range d.DomainRecords {
		for hostName, ip := range record.Hosts {
			records = append(records, A{
				Name: hostName + "." + record.Domain,
				IP:   ip,
				TTL:  record.TTL,
			})
		}
	}

	return records
}

type DataSource struct {
	fileName       string
	records        *RecordsResult
	lastLoadedTime time.Time
}

func NewDataSource(fileName string) *DataSource {
	return &DataSource{
		fileName: fileName,
		records:  NewRecordsResult(),
	}
}

func (ds *DataSource) GetRecords() RecordsResult {
	return *ds.records
}

func (ds *DataSource) GetLastModified() time.Time {
	fileInfo, err := os.Stat(ds.fileName)

	if err != nil {
		log.Error().Err(err).Msg("error reading stat from file")

		return time.Time{}
	}

	return fileInfo.ModTime()
}

func (ds *DataSource) GetLastLoaded() time.Time {
	return ds.lastLoadedTime
}

func (ds *DataSource) Load() {
	b, err := os.ReadFile(ds.fileName)
	if err != nil {
		log.Error().Err(err).Msg("error loading from file")

		return
	}

	ds.loadFromJSON(b)
	ds.lastLoadedTime = time.Now()
}

func (ds *DataSource) loadFromJSON(b []byte) {
	dataSourceEntity := NewDataSourceEntity()
	err := json.Unmarshal(b, &dataSourceEntity)
	if err != nil {
		log.Error().Err(err).Msg("error loading from json")

		return
	}

	ds.records.InitializeWithDataSourceEntity(*dataSourceEntity)
}
