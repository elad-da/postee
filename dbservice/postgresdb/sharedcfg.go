package postgresdb

import (
	"fmt"

	"github.com/aquasecurity/postee/dbservice/dbparam"
	_ "github.com/lib/pq"
)

var apiKeyName = "POSTEE_API_KEY"

func (p *Postgres) EnsureApiKey() error {
	apiKey, err := dbparam.GenerateApiKey(32)
	if err != nil {
		return err
	}

	if err = insertInTableSharedConfig(p.DB, p.TenantName, apiKeyName, apiKey); err != nil {
		return err
	}
	return nil
}

func (p *Postgres) GetApiKey() (string, error) {
	value := ""
	sqlQuery := fmt.Sprintf("SELECT %s FROM %s WHERE (tenantName=$1 AND %s=$2)", "value", dbparam.DbBucketSharedConfig, "apikeyname")
	err := p.DB.Get(&value, sqlQuery, p.TenantName, apiKeyName)
	if err != nil {
		return "", err
	}
	return value, nil
}
