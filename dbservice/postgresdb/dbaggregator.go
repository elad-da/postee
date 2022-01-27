package postgresdb

import (
	"database/sql"
	"encoding/json"
	"errors"
	"fmt"

	"github.com/aquasecurity/postee/dbservice/dbparam"
)

func (p *Postgres) AggregateScans(output string,
	currentScan map[string]string,
	scansPerTicket int,
	ignoreTheQuantity bool) ([]map[string]string, error) {

	aggregatedScans := make([]map[string]string, 0, scansPerTicket)
	if len(currentScan) > 0 {
		aggregatedScans = append(aggregatedScans, currentScan)
	}
	currentValue := []byte{}
	sqlQuery := fmt.Sprintf("SELECT %s FROM %s WHERE (tenantName=$1 AND %s=$2)", "saving", dbparam.DbBucketAggregator, "output")
	err := p.DB.Get(&currentValue, sqlQuery, p.TenantName, output)
	if err != nil {
		if !errors.Is(err, sql.ErrNoRows) {
			return nil, err
		}
	}

	if len(currentValue) > 0 {
		var savedScans []map[string]string
		err = json.Unmarshal([]byte(currentValue), &savedScans)
		if err != nil {
			return nil, err
		}
		aggregatedScans = append(aggregatedScans, savedScans...)
	}

	if ignoreTheQuantity || len(aggregatedScans) < scansPerTicket {
		saving, err := json.Marshal(aggregatedScans)
		if err != nil {
			return nil, err
		}
		if err = insertInTableAggregator(p.DB, p.TenantName, output, saving); err != nil {

			return nil, err
		}
		return nil, nil
	}
	if err = insertInTableAggregator(p.DB, p.TenantName, output, nil); err != nil {
		return nil, err
	}
	return aggregatedScans, nil
}
