package dbservice

import (
	"os"
	"reflect"
	"testing"
)

func TestConfigurateBoltDbPathUsedEnv(t *testing.T) {
	tests := []struct {
		name         string
		dbPath       string
		expectedPath string
	}{
		{"happy configuration BoltDB with dbPath", "database/webhooks.db", "database/webhooks.db"},
		{"happy configuration BoltDB with empty dbPath", "", "/server/database/webhooks.db"},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {

			testInterval := 2
			if err := ConfigureStorage(test.dbPath, "", ""); err != nil {
				t.Errorf("Unexpected error: %v", err)
			}
			if testInterval != 2 {
				t.Error("test interval error, expected: 2, got: ", testInterval)
			}
			if test.expectedPath != reflect.Indirect(reflect.ValueOf(Db)).FieldByName("DbPath").Interface() {
				t.Errorf("paths do not match, expected: %s, got: %s", test.expectedPath, reflect.Indirect(reflect.ValueOf(Db)).FieldByName("DbPath").Interface())
			}
		})
		defer os.RemoveAll("database/")
	}
}

// func TestConfiguratePostgresDbUrlAndTenantName(t *testing.T) {
// 	tests := []struct {
// 		name          string
// 		url           string
// 		tenantName    string
// 		expectedError error
// 	}{
// 		{"happy configuration postgres with url", "postgresql://user:secret@localhost", "test-tenantName", nil},
// 		{"bad tenantName", "postgresql://user:secret@localhost", "", errConfigPsqlEmptyTenantName},
// 		{"bad url", "badUrl", "test-tenantName", errors.New("badUrl error")},
// 	}

// 	for _, test := range tests {
// 		t.Run(test.name, func(t *testing.T) {
// 			initPostgresDbSaved := postgresdb.PsqlConnect
// 			postgresdb.PsqlConnect = func(connectUrl string) (*sqlx.DB, error) {
// 				if connectUrl == "badUrl" {
// 					return nil, test.expectedError
// 				}
// 				return &sqlx.DB{}, nil
// 			}
// 			defer func() {
// 				postgresdb.PsqlConnect = initPostgresDbSaved
// 			}()

// 			err := ConfigureStorage("", test.url, test.tenantName)
// 			if err != nil {
// 				if !errors.Is(err, test.expectedError) {
// 					t.Errorf("Unexpected error, expected: %v, got: %v", test.expectedError, err)
// 				}
// 			}
// 		})
// 	}
// }
