package migrate

import (
	"testing"

	"go-migrations/database"
	"go-migrations/internal"

	"github.com/kylelemons/godebug/pretty"
)

var dbLoadArgsStatus []string
var fakeDbStatus internal.FakeDbWithSpy
var fakeDbLoadedStatus bool

func fakeLoadWithSpyStatus(migrationsPath, environment string) (database.Database, error) {
	dbLoadArgsStatus = []string{migrationsPath, environment}
	fakeDbStatus = internal.FakeDbWithSpy{}
	return &fakeDbStatus, nil
}

func TestMigrateStatus(t *testing.T) {
	mockableLoadDB = fakeLoadWithSpyStatus

	expectedRows := []database.MigrateStatusRow{{ID: "abc"}}
	expectedStatus := "All good"

	mockableGetMigrationStatus = func(
		fileMigrations []database.FileMigration, appliedMigrations []database.AppliedMigration,
	) (rows []database.MigrateStatusRow, statusNote string, err error) {
		return expectedRows, expectedStatus, nil
	}
	var gotRows []database.MigrateStatusRow
	var gotStatus string
	mockablePrintStatusTable = func(rows []database.MigrateStatusRow, statusNote string) {
		gotRows = rows
		gotStatus = statusNote
	}

	args := []string{"sth.exe", "migrate", "status", "-p", "./sth/else", "-e", "my-env"}
	if err := app.Run(args); err != nil {
		t.Errorf("Error running command - %s", err)
	}

	fakeDbStatus.AssertGetFileMigrationsCalled(t, true)
	fakeDbStatus.AssertGetAppliedMigrationsCalled(t, true)

	if diff := pretty.Compare(expectedRows, gotRows); diff != "" {
		t.Errorf("Did not pass right rows for print:\n%s", diff)
	}
	if expectedStatus != gotStatus {
		t.Errorf("Expected status of '%s' but got '%s'", expectedStatus, gotStatus)
	}

}
