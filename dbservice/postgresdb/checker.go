package postgresdb

import (
	"fmt"
	"time"

	"github.com/aquasecurity/postee/dbservice/dbparam"
	"github.com/aquasecurity/postee/log"
)

func (p *Postgres) CheckSizeLimit() {
	if dbparam.DbSizeLimit == 0 {
		return
	}

	size := 0
	err := p.DB.Get(&size, fmt.Sprintf("SELECT pg_total_relation_size('%s');", dbparam.DbBucketName))
	if err != nil {
		log.Logger.Error("CheckSizeLimit: Can't get db size")
		return
	}
	if size > dbparam.DbSizeLimit {
		if err = deleteRowsByTenantName(p.DB, dbparam.DbBucketName, p.TenantName); err != nil {
			log.Logger.Errorf("CheckSizeLimit: Can't delete tenantName's: %s from table: %s", p.TenantName, dbparam.DbBucketName)
			return
		}
	}
}

func (p *Postgres) CheckExpiredData() {
	max := time.Now().UTC() //remove expired records
	if err := deleteRowsByTenantNameAndTime(p.DB, p.TenantName, max); err != nil {
		log.Logger.Errorf("CheckExpiredData: Can't delete dates from table:%s, err: %v", dbparam.DbBucketName, err)
	}
}
