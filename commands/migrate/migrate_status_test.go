package migrate

import (
	"testing"

	"go-migrations/database"
	"go-migrations/internal"
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

	args := []string{"sth.exe", "migrate", "status", "-p", "./sth/else", "-e", "my-env"}
	if err := app.Run(args); err != nil {
		t.Errorf("Error running command - %s", err)
	}

	expected := []string{"./sth/else", "my-env"}
	if !internal.StrSliceEqual(dbLoadArgsStatus, expected) {
		t.Errorf("Expected to load db with '%v', but got %s", expected, dbLoadArgsStatus)
	}

	fakeDbStatus.AssertWaitForStartCalled(t, true)
	fakeDbStatus.AssertEnsureMigrationsChangelogCalled(t, true)
	fakeDbStatus.AssertPrintMigrationStatusCalled(t, true)
}
