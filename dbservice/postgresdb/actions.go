package postgresdb

import (
	"database/sql"
	"errors"
	"fmt"
	"time"

	"github.com/aquasecurity/postee/dbservice/dbparam"
)

func (p *Postgres) MayBeStoreMessage(message []byte, messageKey string, expired *time.Time) (wasStored bool, err error) {
	currentValue := ""
	sqlQuery := fmt.Sprintf("SELECT messageValue FROM %s WHERE (tenantName=$1 AND messageKey=$2)", dbparam.DbBucketName)
	if err = p.DB.Get(&currentValue, sqlQuery, p.TenantName, messageKey); err != nil {
		if !errors.Is(err, sql.ErrNoRows) {
			return false, err
		}
	}

	if currentValue != "" {
		return false, nil
	} else {
		if err = insertInTableName(p.DB, p.TenantName, messageKey, message, expired); err != nil {
			return false, err
		}
		return true, nil
	}
}
