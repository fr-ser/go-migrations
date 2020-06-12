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
	mockableApplyDownMigration = database.ApplyDownMigration
	mockableFilterUpMigrationsByText = database.FilterUpMigrationsByText
	mockableFilterDownMigrationsByText = database.FilterDownMigrationsByText
	mockableFilterUpMigrationsByCount = database.FilterUpMigrationsByCount
	mockableFilterDownMigrationsByCount = database.FilterDownMigrationsByCount
}

func TestMain(m *testing.M) {
	log.SetOutput(ioutil.Discard)
	os.Exit(m.Run())
}
