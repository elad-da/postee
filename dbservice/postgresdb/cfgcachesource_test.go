package postgresdb

import (
	"database/sql"
	"errors"
	"log"
	"testing"

	"github.com/jmoiron/sqlx"
	sqlxmock "github.com/zhashkevych/go-sqlxmock"
)

func TestUpdateCfgCacheSource(t *testing.T) {
	cfgFile := `{"name": "tenant", "aqua-server": "https://myserver.aquasec.com"}`

	db, mock, err := sqlxmock.Newx()
	if err != nil {
		t.Errorf("failed to open sqlmock database: %v", err)
	}
	rows := sqlxmock.NewRows([]string{"cfgFile"}).AddRow(cfgFile)
	mock.ExpectQuery("SELECT").WillReturnRows(rows)
	pgDB := &Postgres{DB: db}

	savedInsertCfgCacheSource := insertCfgCacheSource
	insertCfgCacheSource = func(db *sqlx.DB, tenantName, cfgFile string) error { return nil }
	defer func() {
		insertCfgCacheSource = savedInsertCfgCacheSource
	}()

	if err := UpdateCfgCacheSource(pgDB, "cfgFile"); err != nil {
		t.Errorf("Unexpected error: %v", err)
	}

	cfg, err := GetCfgCacheSource(pgDB)
	if err != nil {
		t.Errorf("Unexpected error: %v", err)
	}
	if cfgFile != cfg {
		t.Errorf("CfgFiles not equals, expected: %s, got: %s", cfgFile, cfg)
	}
}

func TestGetCfgCacheSourceErrors(t *testing.T) {
	tests := []struct {
		name        string
		err         error
		expectedCfg string
		expectedErr string
	}{
		{"Norows error", sql.ErrNoRows, "{}", ""},
		{"select error", errors.New("select error"), "", "error getting cfg cache source: select error"},
	}

	savedInsertCfgCacheSource := insertCfgCacheSource
	insertCfgCacheSource = func(db *sqlx.DB, tenantName, cfgFile string) error { return nil }
	defer func() {
		insertCfgCacheSource = savedInsertCfgCacheSource
	}()

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			db, mock, err := sqlxmock.Newx()
			if err != nil {
				log.Println("failed to open sqlmock database:", err)
			}
			mock.ExpectQuery("SELECT").WillReturnError(test.err)
			pgDB := &Postgres{DB: db}

			cfg, err := GetCfgCacheSource(pgDB)
			if test.expectedErr != "" || err != nil {
				if err.Error() != test.expectedErr {
					t.Errorf("Unexpected err, expected: %v, got: %v", test.expectedErr, err)
				}
			}

			if cfg != test.expectedCfg {
				t.Errorf("Bad cfg, expected: %s, got: %s", test.expectedCfg, cfg)
			}
		})

	}

}
