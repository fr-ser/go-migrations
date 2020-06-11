package postgres

import (
	"database/sql"
	"io/ioutil"
	"os"
	"testing"

	log "github.com/sirupsen/logrus"

	"go-migrations/database"
)

func resetMockVariables() {
	mockableSQLOpen = sql.Open
	mockableWaitForStart = database.WaitForStart
	mockableBootstrap = database.ApplyBootstrapMigration
	mockableEnsureConsistentMigrations = database.EnsureConsistentMigrations
	mockableGetFileMigrations = database.GetFileMigrations
	mockableGetAppliedMigrations = database.GetAppliedMigrations
	mockableApplyUpMigration = database.ApplyUpMigration
	mockableFilterUpMigrationsByText = database.FilterUpMigrationsByText
	mockableFilterUpMigrationsByCount = database.FilterUpMigrationsByCount
}

func TestMain(m *testing.M) {
	log.SetOutput(ioutil.Discard)
	os.Exit(m.Run())
}
