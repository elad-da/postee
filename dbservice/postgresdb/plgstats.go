package postgresdb

import (
	"database/sql"
	"errors"
	"fmt"

	"github.com/aquasecurity/postee/dbservice/dbparam"
	_ "github.com/lib/pq"
)

func (p *Postgres) RegisterPlgnInvctn(name string) error {
	amount := 0
	sqlQuery := fmt.Sprintf("SELECT %s FROM %s WHERE (tenantName=$1 AND %s=$2)", "amount", dbparam.DbBucketOutputStats, "outputName")
	err := p.DB.Get(&amount, sqlQuery, p.TenantName, name)
	if err != nil && !errors.Is(err, sql.ErrNoRows) {
		return err
	}
	amount += 1
	err = insertOutputStats(p.DB, p.TenantName, name, amount)
	if err != nil {
		return err
	}
	return nil
}
