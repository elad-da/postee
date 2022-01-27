package postgresdb

import (
	"database/sql"
	"testing"

	"github.com/jmoiron/sqlx"
	sqlxmock "github.com/zhashkevych/go-sqlxmock"
)

func TestApiKey(t *testing.T) {
	savedInsertInTableSharedConfig := insertInTableSharedConfig
	insertInTableSharedConfig = func(db *sqlx.DB, tenantName, apikeyname, value string) error { return nil }

	db, mock, err := sqlxmock.Newx()
	if err != nil {
		t.Errorf("failed to open sqlmock database: %v", err)
	}
	rows := sqlxmock.NewRows([]string{"value"}).AddRow("key")
	mock.ExpectQuery("SELECT").WillReturnRows(rows)
	pgDB := &Postgres{DB: db}

	defer func() {
		insertInTableSharedConfig = savedInsertInTableSharedConfig
	}()

	if err := pgDB.EnsureApiKey(); err != nil {
		t.Errorf("Unexpected EnsureApiKey error: %v", err)
	}

	key, err := pgDB.GetApiKey()
	if err != nil {
		t.Fatal("error while getting value of API key")
	}
	if key == "" {
		t.Fatal("empty key received")
	}
}

func TestApiKeyWithoutInit(t *testing.T) {

	db, mock, err := sqlxmock.Newx()
	if err != nil {
		t.Errorf("failed to open sqlmock database: %v", err)
	}
	mock.ExpectQuery("SELECT").WillReturnError(sql.ErrNoRows)
	pgDB := &Postgres{DB: db}

	key, err := pgDB.GetApiKey()
	if err == nil {
		t.Fatal("Error is expected")
	}
	if key != "" {
		t.Fatal("Empty key is expected")
	}
}

func TestApiKeyRenewal(t *testing.T) {
	receivedKey := ""
	savedInsertInTableSharedConfig := insertInTableSharedConfig
	insertInTableSharedConfig = func(db *sqlx.DB, tenantName, apikeyname, value string) error {
		receivedKey = value
		return nil
	}

	db, mock, err := sqlxmock.Newx()
	if err != nil {
		t.Errorf("failed to open sqlmock database: %v", err)
	}
	defer db.Close()
	pgDB := &Postgres{DB: db}

	defer func() {
		insertInTableSharedConfig = savedInsertInTableSharedConfig
	}()

	var keys [2]string
	for i := 0; i < 2; i++ {
		if err := pgDB.EnsureApiKey(); err != nil {
			t.Errorf("Unexpected EnsureApiKey error: %v", err)
		}

		rows := sqlxmock.NewRows([]string{"value"}).AddRow(receivedKey)
		mock.ExpectQuery("SELECT").WillReturnRows(rows)
		key, err := pgDB.GetApiKey()
		if err != nil {
			t.Fatal("error while getting value of API key")
		}
		if key == "" {
			t.Fatal("empty key received")
		}
		keys[i] = key
	}
	if keys[0] == keys[1] {
		t.Errorf("Key is not updated. (before: %s and after update: %s)", keys[0], keys[1])
	}
}
