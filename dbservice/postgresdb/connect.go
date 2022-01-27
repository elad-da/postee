package postgresdb

import (
	"errors"

	"github.com/jmoiron/sqlx"
)

type Postgres struct {
	DB         *sqlx.DB
	TenantName string
}

func New(connectURL string, tenantName string) (*Postgres, error) {
	db, err := PsqlConnect(connectURL)
	if err != nil {
		return nil, err
	}
	err = InitAllTables(db)
	if err != nil {
		return nil, err
	}

	return &Postgres{DB: db, TenantName: tenantName}, nil
}

//TODO: retry support
var PsqlConnect = func(connectUrl string) (*sqlx.DB, error) {
	db, err := sqlx.Connect("postgres", connectUrl)
	if err != nil {
		return nil, err
	}
	return db, nil
}

var testConnect = func(connectUrl string) (*sqlx.DB, error) {
	db, err := PsqlConnect(connectUrl)
	if err != nil {
		return nil, errors.New("Error postgresDb test connect: " + err.Error())
	}
	return db, nil
}
