package postgresdb

import (
	"testing"
	"time"

	"github.com/aquasecurity/postee/dbservice/dbparam"
	"github.com/jmoiron/sqlx"
	sqlxmock "github.com/zhashkevych/go-sqlxmock"
)

func TestExpiredDates(t *testing.T) {
	tests := []struct {
		name        string
		deleteError bool
		wasDeleted  bool
	}{
		{"happy delete rows", false, true},
		{"bad delete rows", true, false},
	}

	savedDeleteRow := deleteRowsByTenantNameAndTime
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			deleted := false

			db, _, err := sqlxmock.Newx()
			if err != nil {
				t.Errorf("failed to open sqlmock database: %v", err)
			}
			pgDB := &Postgres{DB: db}

			deleteRowsByTenantNameAndTime = func(db *sqlx.DB, tenantName string, t time.Time) error {
				if !test.deleteError {
					deleted = true
				}
				return nil
			}

			pgDB.CheckExpiredData()
			if deleted != test.wasDeleted {
				t.Errorf("error deleted rows")
			}

			deleteRowsByTenantNameAndTime = savedDeleteRow
		})
	}
}

func TestSizeLimit(t *testing.T) {
	tests := []struct {
		name       string
		sizeLimit  int
		size       int
		wasDeleted bool
	}{
		{"No size limit", 0, 10, false},
		{"Size less then limit", 5, 10, true},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			deleted := false

			db, mock, err := sqlxmock.Newx()
			if err != nil {
				t.Errorf("failed to open sqlmock database: %v", err)
			}
			rows := sqlxmock.NewRows([]string{"size"}).AddRow(test.size)
			mock.ExpectQuery("SELECT").WillReturnRows(rows)
			pgDB := &Postgres{DB: db}

			savedDeleteRowsByTenantName := deleteRowsByTenantName
			deleteRowsByTenantName = func(db *sqlx.DB, table, tenantName string) error {
				deleted = true
				return nil
			}

			dbparam.DbSizeLimit = test.sizeLimit
			pgDB.CheckSizeLimit()
			if deleted != test.wasDeleted {
				t.Errorf("error deleted rows")
			}

			deleteRowsByTenantName = savedDeleteRowsByTenantName
		})
	}
}
