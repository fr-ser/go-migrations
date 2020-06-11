package postgres

import (
	"database/sql"
	"go-migrations/database"
)

func resetMockVariables() {
	sqlOpen = sql.Open
	commonWaitForStart = database.WaitForStart
	commonBootstrap = database.ApplyBootstrapMigration
	commonGetFileMigrations = database.GetFileMigrations
	commonEnsureConsistentMigrations = database.EnsureConsistentMigrations
}
