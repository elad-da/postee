package postgresdb

import (
	"crypto/rand"
	"encoding/hex"
	"fmt"
	"io"

	_ "github.com/lib/pq"
)

var apiKeyName = "POSTEE_API_KEY"

func (postgresDb *PostgresDb) EnsureApiKey() error {
	db, err := psqlConnect(postgresDb.psqlInfo)
	if err != nil {
		return err
	}
	defer db.Close()
	err = initTable(db, dbTableSharedConfig)
	if err != nil {
		return err
	}

	apiKey, err := generateApiKey(32)
	if err != nil {
		return err
	}

	if err = insert(db, dbTableSharedConfig, postgresDb.id, "apikeyname", apiKeyName, "value", apiKey); err != nil {
		return err
	}

	return nil
}

func (postgresDb *PostgresDb) GetApiKey() (string, error) {
	db, err := psqlConnect(postgresDb.psqlInfo)
	if err != nil {
		return "", err
	}
	defer db.Close()
	value := ""
	err = db.Get(&value, fmt.Sprintf("SELECT %s FROM %s WHERE (%s=$1 AND %s=$2)", "value", dbTableSharedConfig, "id", "apikeyname"), postgresDb.id, apiKeyName)
	if err != nil {
		return "", err
	}
	return value, nil

}

func generateApiKey(length int) (string, error) {
	k := make([]byte, length)
	if _, err := io.ReadFull(rand.Reader, k); err != nil {
		return "", err
	}
	return hex.EncodeToString(k), nil
}