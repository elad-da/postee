package postgresdb

import (
	"database/sql"
	"errors"
	"testing"

	"github.com/jmoiron/sqlx"
	sqlxmock "github.com/zhashkevych/go-sqlxmock"
)

func TestRegisterPlgnInvctn(t *testing.T) {
	receivedKey := 0
	savedInsertOutputStats := insertOutputStats
	insertOutputStats = func(db *sqlx.DB, tenantName, outputName string, amount int) error {
		receivedKey = amount
		return nil
	}

	db, mock, err := sqlxmock.Newx()
	if err != nil {
		t.Errorf("failed to open sqlmock database: %v", err)
	}
	rows := sqlxmock.NewRows([]string{"amount"}).AddRow(receivedKey)
	mock.ExpectQuery("SELECT").WillReturnRows(rows)
	pgDB := &Postgres{DB: db}

	defer func() {
		insertOutputStats = savedInsertOutputStats
	}()

	expectedCnt := 3
	keyToTest := "test"
	for i := 0; i < expectedCnt; i++ {
		if err := pgDB.RegisterPlgnInvctn(keyToTest); err != nil {
			t.Errorf("unexpected error: %v", err)
		}
	}
	if receivedKey != expectedCnt {
		t.Errorf("Persisted count doesn't match expected. Expected %d, got %d\n", receivedKey, expectedCnt)
	}
}

func TestRegisterPlgnInvctnErrors(t *testing.T) {
	var tests = []struct {
		name        string
		errIn       error
		expectedErr error
	}{
		{"No result rows error", sql.ErrNoRows, nil},
		{"Other errors", sql.ErrConnDone, sql.ErrConnDone},
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			savedInsertOutputStats := insertOutputStats
			insertOutputStats = func(db *sqlx.DB, tenantName, outputName string, amount int) error { return nil }

			db, mock, err := sqlxmock.Newx()
			if err != nil {
				t.Errorf("failed to open sqlmock database: %v", err)
			}
			mock.ExpectQuery("SELECT").WillReturnError(test.errIn)
			pgDB := &Postgres{DB: db}

			err = pgDB.RegisterPlgnInvctn("testName")
			if !errors.Is(err, test.expectedErr) {
				t.Errorf("Errors no contains: expected: %v, got: %v", test.expectedErr, err)
			}

			insertOutputStats = savedInsertOutputStats

		})
	}

	db, mock, err := sqlxmock.Newx()
	if err != nil {
		t.Errorf("failed to open sqlmock database: %v", err)
	}
	mock.ExpectQuery("SELECT").WillReturnError(sql.ErrConnDone)
	pgDB := &Postgres{DB: db}

	key, err := pgDB.GetApiKey()
	if err == nil {
		t.Fatal("Error is expected")
	}
	if key != "" {
		t.Fatal("Empty key is expected")
	}
}
