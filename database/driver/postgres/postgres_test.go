package postgres

import (
	"database/sql"
	"go-migrations/database"
)

func resetMockVariables() {
	sqlOpen = sql.Open
	commonWaitForStart = database.WaitForStart
	commonBootstrap = database.ApplyBootstrapMigration
	commonEnsureConsistentMigrations = database.EnsureConsistentMigrations
	commonGetFileMigrations = database.GetFileMigrations
	commonGetAppliedMigrations = database.GetAppliedMigrations
}
