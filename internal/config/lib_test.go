package config

import (
	"os"
	"testing"

	"github.com/google/uuid"
	"gotest.tools/v3/assert"
)

func TestLoadConfig(t *testing.T) {
	conf, err := LoadConfig()
	assert.NilError(t, err)

	assert.Equal(t, conf.DB.DBName, "postgres")

	dbName := uuid.NewString()
	os.Setenv("RACING_DRAFT_DB_DBNAME", dbName)

	conf, err = LoadConfig()
	assert.NilError(t, err)
	assert.Equal(t, conf.DB.DBName, dbName)
}
