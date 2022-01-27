package dbservice

import (
	"fmt"
	"time"

	"github.com/aquasecurity/postee/dbservice/boltdb"
	"github.com/aquasecurity/postee/dbservice/postgresdb"
)

var (
	Db DbProvider

	errConfigPsqlEmptyTenantName = fmt.Errorf("error configuring postgres: 'tenantName' is empty")
)

type DbProvider interface {
	MayBeStoreMessage(message []byte, messageKey string, expired *time.Time) (wasStored bool, err error)
	CheckSizeLimit()
	CheckExpiredData()
	AggregateScans(output string, currentScan map[string]string, scansPerTicket int, ignoreTheQuantity bool) ([]map[string]string, error)
	RegisterPlgnInvctn(name string) error
	EnsureApiKey() error
	GetApiKey() (string, error)
}

func ConfigureStorage(boltDBPath, postgresUrl, tenantName string) error {
	if postgresUrl != "" {
		return ConfigurePostgresDB(postgresUrl, tenantName)
	}

	return ConfigureBoltDB(boltDBPath)
}

func ConfigureBoltDB(pathToDb string) error {
	boltdb := boltdb.NewBoltDb()
	if pathToDb != "" {
		if err := boltdb.SetNewDbPath(pathToDb); err != nil {
			return err
		}
	}
	Db = boltdb
	return nil
}

func ConfigurePostgresDB(postgresUrl, tenantName string) error {
	if tenantName == "" {
		return errConfigPsqlEmptyTenantName
	}
	postgresDb, err := postgresdb.New(postgresUrl, tenantName)
	if err != nil {
		return err
	}
	Db = postgresDb
	return nil
}
