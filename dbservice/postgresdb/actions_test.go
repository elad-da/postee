package postgresdb

import (
	"testing"
	"time"

	"github.com/jmoiron/sqlx"
	sqlxmock "github.com/zhashkevych/go-sqlxmock"
)

func TestStoreMessage(t *testing.T) {
	currentValueStoreMessage := []byte{}

	db, mock, err := sqlxmock.Newx()
	if err != nil {
		t.Errorf("failed to open sqlmock database: %v", err)
	}
	rows := sqlxmock.NewRows([]string{"messagevalue"}).AddRow(currentValueStoreMessage)
	mock.ExpectQuery("SELECT (.+) FROM WebhookBucket").WithArgs("tenantName", AlpineImageKey).WillReturnRows(rows)
	pgDB := &Postgres{DB: db, TenantName: "tenantName"}

	savedinsertInTableName := insertInTableName
	insertInTableName = func(db *sqlx.DB, tenantName, messageKey string, messageValue []byte, date *time.Time) error {
		currentValueStoreMessage = messageValue
		return nil
	}

	defer func() {
		insertInTableName = savedinsertInTableName
	}()

	var tests = []struct {
		input *string
		t     *time.Time
	}{
		{&AlpineImageResult, nil},
		{&AlpineImageResult, &time.Time{}},
	}

	for _, test := range tests {

		// Handling of first scan
		isNew, err := pgDB.MayBeStoreMessage([]byte(*test.input), AlpineImageKey, test.t)
		if err != nil {
			t.Errorf("Error: %s\n", err)
		}
		if !isNew {
			t.Errorf("A first scan was found!\n")
		}

		// Handling of second scan with the same data
		isNew, err = pgDB.MayBeStoreMessage([]byte(*test.input), AlpineImageKey, nil)
		if err != nil {
			t.Errorf("Error: %s\n", err)
		}
		if isNew {
			t.Errorf("A old scan wasn't found!\n")
		}
		currentValueStoreMessage = []byte{}
	}
}
