package postgres

import (
	"database/sql"
	"go-migrations/database/common"
)

func resetMockVariables() {
	sqlOpen = sql.Open
	commonWaitForStart = common.WaitForStart
	commonBootstrap = common.ApplyBootstrapMigration
	commonGetMigrations = common.GetMigrations
}
