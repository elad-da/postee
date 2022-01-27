package postgresdb

import (
	"database/sql"
	"errors"
	"fmt"

	"github.com/aquasecurity/postee/dbservice/dbparam"
)

const emptyCfg = `{}`

var UpdateCfgCacheSource = func(p *Postgres, cfgfile string) error {
	if err := insertCfgCacheSource(p.DB, p.TenantName, cfgfile); err != nil {
		return err
	}
	return nil
}

var GetCfgCacheSource = func(p *Postgres) (string, error) {
	cfgFile := ""
	sqlQuery := fmt.Sprintf("SELECT configfile FROM %s WHERE tenantName=$1", dbparam.DbTableCfgCacheSource)
	if err := p.DB.Get(&cfgFile, sqlQuery, p.TenantName); err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return emptyCfg, nil
		}
		return "", fmt.Errorf("error getting cfg cache source: %w", err)
	}
	return cfgFile, nil
}
