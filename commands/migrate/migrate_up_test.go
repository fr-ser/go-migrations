package migrate

import (
	"testing"

	"go-migrations/database"
	"go-migrations/internal"
)

var dbLoadArgsUp []string
var fakeDbUp internal.FakeDbWithSpy
var fakeDbLoadedUp bool

func fakeLoadWithSpyUp(migrationsPath, environment string) (database.Database, error) {
	dbLoadArgsUp = []string{migrationsPath, environment}
	fakeDbUp = internal.FakeDbWithSpy{}
	return &fakeDbUp, nil
}
func TestMigrateUpDefaults(t *testing.T) {
	mockableLoadDB = fakeLoadWithSpyUp

	args := []string{"sth.exe", "migrate", "up"}
	if err := app.Run(args); err != nil {
		t.Errorf("Error running command - %s", err)
	}

	expected := []string{"./migrations/zlab", "development"}
	if !internal.StrSliceEqual(dbLoadArgsUp, expected) {
		t.Errorf("Expected to load db with '%v', but got %s", expected, dbLoadArgsUp)
	}

	fakeDbUp.AssertWaitForStartCalled(t, true)
	fakeDbUp.AssertEnsureMigrationsChangelogCalled(t, true)
	fakeDbUp.AssertEnsureConsistentMigrationsCalled(t, true)
	fakeDbUp.AssertApplyUpMigrationsWithCountCalledWith(t, 1, false)
	fakeDbUp.AssertApplySpecificUpMigrationCalled(t, false)
}

func TestMigrateUpWithCount(t *testing.T) {
	mockableLoadDB = fakeLoadWithSpyUp

	args := []string{"sth.exe", "migrate", "up", "--count", "2"}
	if err := app.Run(args); err != nil {
		t.Errorf("Error running command - %s", err)
	}

	expected := []string{"./migrations/zlab", "development"}
	if !internal.StrSliceEqual(dbLoadArgsUp, expected) {
		t.Errorf("Expected to load db with '%v', but got %s", expected, dbLoadArgsUp)
	}

	fakeDbUp.AssertWaitForStartCalled(t, true)
	fakeDbUp.AssertEnsureMigrationsChangelogCalled(t, true)
	fakeDbUp.AssertEnsureConsistentMigrationsCalled(t, true)
	fakeDbUp.AssertApplyUpMigrationsWithCountCalledWith(t, 2, false)
	fakeDbUp.AssertApplySpecificUpMigrationCalled(t, false)

}

func TestMigrateUpWithAll(t *testing.T) {
	mockableLoadDB = fakeLoadWithSpyUp

	args := []string{"sth.exe", "migrate", "up", "--all"}
	if err := app.Run(args); err != nil {
		t.Errorf("Error running command - %s", err)
	}

	expected := []string{"./migrations/zlab", "development"}
	if !internal.StrSliceEqual(dbLoadArgsUp, expected) {
		t.Errorf("Expected to load db with '%v', but got %s", expected, dbLoadArgsUp)
	}

	fakeDbUp.AssertWaitForStartCalled(t, true)
	fakeDbUp.AssertEnsureMigrationsChangelogCalled(t, true)
	fakeDbUp.AssertEnsureConsistentMigrationsCalled(t, true)
	fakeDbUp.AssertApplyUpMigrationsWithCountCalledWith(t, 0, true)
	fakeDbUp.AssertApplySpecificUpMigrationCalled(t, false)

}

func TestMigrateUpWithOnly(t *testing.T) {
	mockableLoadDB = fakeLoadWithSpyUp

	args := []string{"sth.exe", "migrate", "up", "--only", "sth"}
	if err := app.Run(args); err != nil {
		t.Errorf("Error running command - %s", err)
	}

	expected := []string{"./migrations/zlab", "development"}
	if !internal.StrSliceEqual(dbLoadArgsUp, expected) {
		t.Errorf("Expected to load db with '%v', but got %s", expected, dbLoadArgsUp)
	}

	fakeDbUp.AssertWaitForStartCalled(t, true)
	fakeDbUp.AssertEnsureMigrationsChangelogCalled(t, true)
	fakeDbUp.AssertApplySpecificUpMigrationCalledWith(t, "sth")
	fakeDbUp.AssertApplyUpMigrationsWithCountCalled(t, false)
	fakeDbUp.AssertEnsureConsistentMigrationsCalled(t, false)
}

func TestMigrateUpErrorWithMultipleParams(t *testing.T) {
	fakeDbLoadedUp = false
	mockableLoadDB = fakeLoadWithSpyUp

	var invalidArgs = [][]string{
		{"sth.exe", "migrate", "up", "--only", "sth", "--all"},
		{"sth.exe", "migrate", "up", "--only", "sth", "--count", "1"},
		{"sth.exe", "migrate", "up", "--all", "--count", "1"},
		{"sth.exe", "migrate", "up", "--count", "-2"},
		{"sth.exe", "migrate", "up", "--count", "0"},
		{"sth.exe", "migrate", "up", "--count", "two"},
		{"sth.exe", "migrate", "up", "--all", "--count", "1", "--only", "sth"},
	}

	for _, args := range invalidArgs {
		if err := app.Run(args); err == nil {
			t.Errorf("Got no error for wrong parameters: %v", args)
		}
		if fakeDbLoadedUp {
			t.Errorf("Loaded the database even though it shouldn't")
		}

	}

}
