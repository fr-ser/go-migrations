package migrate

import (
	"testing"

	"go-migrations/database"
	"go-migrations/internal"
	"go-migrations/internal/direction"
)

var dbLoadArgsDown []string
var fakeDbDown internal.FakeDbWithSpy
var fakeDbLoadedDown bool

func fakeLoadWithSpyDown(migrationsPath, environment string) (database.Database, error) {
	dbLoadArgsDown = []string{migrationsPath, environment}
	fakeDbDown = internal.FakeDbWithSpy{}
	return &fakeDbDown, nil
}
func TestMigrateDownDefaults(t *testing.T) {
	mockableLoadDB = fakeLoadWithSpyDown

	args := []string{"sth.exe", "migrate", "down"}
	if err := app.Run(args); err != nil {
		t.Errorf("Error running command - %s", err)
	}

	expected := []string{"./migrations/zlab", "development"}
	if !internal.StrSliceEqual(dbLoadArgsDown, expected) {
		t.Errorf("Expected to load db with '%v', but got %s", expected, dbLoadArgsDown)
	}

	fakeDbDown.AssertWaitForStartCalled(t, true)
	fakeDbDown.AssertEnsureMigrationsChangelogCalled(t, true)
	fakeDbDown.AssertEnsureConsistentMigrationsCalled(t, true)
	fakeDbDown.AssertApplyMigrationsWithCountCalledWith(t, 1, false, direction.Down)
	fakeDbDown.AssertApplySpecificMigrationCalled(t, false)
}

func TestMigrateDownWithCount(t *testing.T) {
	mockableLoadDB = fakeLoadWithSpyDown

	args := []string{"sth.exe", "migrate", "down", "--count", "2"}
	if err := app.Run(args); err != nil {
		t.Errorf("Error running command - %s", err)
	}

	expected := []string{"./migrations/zlab", "development"}
	if !internal.StrSliceEqual(dbLoadArgsDown, expected) {
		t.Errorf("Expected to load db with '%v', but got %s", expected, dbLoadArgsDown)
	}

	fakeDbDown.AssertWaitForStartCalled(t, true)
	fakeDbDown.AssertEnsureMigrationsChangelogCalled(t, true)
	fakeDbDown.AssertEnsureConsistentMigrationsCalled(t, true)
	fakeDbDown.AssertApplyMigrationsWithCountCalledWith(t, 2, false, direction.Down)
	fakeDbDown.AssertApplySpecificMigrationCalled(t, false)

}

func TestMigrateDownWithAll(t *testing.T) {
	mockableLoadDB = fakeLoadWithSpyDown

	args := []string{"sth.exe", "migrate", "down", "--all"}
	if err := app.Run(args); err != nil {
		t.Errorf("Error running command - %s", err)
	}

	expected := []string{"./migrations/zlab", "development"}
	if !internal.StrSliceEqual(dbLoadArgsDown, expected) {
		t.Errorf("Expected to load db with '%v', but got %s", expected, dbLoadArgsDown)
	}

	fakeDbDown.AssertWaitForStartCalled(t, true)
	fakeDbDown.AssertEnsureMigrationsChangelogCalled(t, true)
	fakeDbDown.AssertEnsureConsistentMigrationsCalled(t, true)
	fakeDbDown.AssertApplyMigrationsWithCountCalledWith(t, 0, true, direction.Down)
	fakeDbDown.AssertApplySpecificMigrationCalled(t, false)

}

func TestMigrateDownWithOnly(t *testing.T) {
	mockableLoadDB = fakeLoadWithSpyDown

	args := []string{"sth.exe", "migrate", "down", "--only", "sth"}
	if err := app.Run(args); err != nil {
		t.Errorf("Error running command - %s", err)
	}

	expected := []string{"./migrations/zlab", "development"}
	if !internal.StrSliceEqual(dbLoadArgsDown, expected) {
		t.Errorf("Expected to load db with '%v', but got %s", expected, dbLoadArgsDown)
	}

	fakeDbDown.AssertWaitForStartCalled(t, true)
	fakeDbDown.AssertEnsureMigrationsChangelogCalled(t, true)
	fakeDbDown.AssertApplySpecificMigrationCalledWith(t, "sth", direction.Down)
	fakeDbDown.AssertApplyMigrationsWithCountCalled(t, false)
	fakeDbDown.AssertEnsureConsistentMigrationsCalled(t, false)
}

func TestMigrateDownErrorWithMultipleParams(t *testing.T) {
	fakeDbLoadedDown = false
	mockableLoadDB = fakeLoadWithSpyDown

	var invalidArgs = [][]string{
		{"sth.exe", "migrate", "down", "--only", "sth", "--all"},
		{"sth.exe", "migrate", "down", "--only", "sth", "--count", "1"},
		{"sth.exe", "migrate", "down", "--all", "--count", "1"},
		{"sth.exe", "migrate", "down", "--count", "-2"},
		{"sth.exe", "migrate", "down", "--count", "0"},
		{"sth.exe", "migrate", "down", "--count", "two"},
		{"sth.exe", "migrate", "down", "--all", "--count", "1", "--only", "sth"},
	}

	for _, args := range invalidArgs {
		if err := app.Run(args); err == nil {
			t.Errorf("Got no error for wrong parameters: %v", args)
		}
		if fakeDbLoadedDown {
			t.Errorf("Loaded the database even though it shouldn't")
		}

	}

}
