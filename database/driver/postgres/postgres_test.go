package postgres

import (
	"database/sql"
	"go-migrations/database"
)

func resetMockVariables() {
	sqlOpen = sql.Open
	commonWaitForStart = database.WaitForStart
	commonBootstrap = database.ApplyBootstrapMigration
	commonGetMigrations = database.GetMigrations
	commonEnsureConsistentMigrations = database.EnsureConsistentMigrations
}
