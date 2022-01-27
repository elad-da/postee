package postgresdb

import (
	"testing"
	"time"

	"github.com/jmoiron/sqlx"
	sqlxmock "github.com/zhashkevych/go-sqlxmock"
)

func TestStoreMessage(t *testing.T) {
	currentValueStoreMessage := []byte{}
	savedinsertInTableName := insertInTableName
	insertInTableName = func(db *sqlx.DB, tenantName, messageKey string, messageValue []byte, date *time.Time) error {
		currentValueStoreMessage = messageValue
		return nil
	}

	db, mock, err := sqlxmock.Newx()
	if err != nil {
		t.Errorf("failed to open sqlmock database: %v", err)
	}
	defer db.Close()
	rows := sqlxmock.NewRows([]string{"messagevalue"}).AddRow("")
	mock.ExpectQuery("SELECT").WillReturnRows(rows)

	pgDB := &Postgres{DB: db, TenantName: "tenantName"}
	defer func() {
		insertInTableName = savedinsertInTableName
	}()

	// Handling of first scan
	isNew, err := pgDB.MayBeStoreMessage([]byte(AlpineImageResult), AlpineImageKey, nil)
	if err != nil {
		t.Errorf("Error: %s\n", err)
	}
	if !isNew {
		t.Errorf("A first scan was found!\n")
	}

	rows = sqlxmock.NewRows([]string{"messagevalue"}).AddRow(currentValueStoreMessage)
	mock.ExpectQuery("SELECT").WillReturnRows(rows)

	// Handling of second scan with the same data
	isNew, err = pgDB.MayBeStoreMessage([]byte(AlpineImageResult), AlpineImageKey, nil)
	if err != nil {
		t.Errorf("Error: %s\n", err)
	}
	if isNew {
		t.Errorf("A old scan wasn't found!\n")
	}
	currentValueStoreMessage = []byte{}

}
