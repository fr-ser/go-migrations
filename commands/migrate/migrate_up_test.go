package migrate

import (
	"io/ioutil"
	"log"
	"os"
	"testing"

	"github.com/urfave/cli/v2"

	"go-migrations/database"
	"go-migrations/internal"
)

var app = cli.NewApp()

var dbLoadArgs []string
var fakeDb internal.FakeDbWithSpy
var fakeDbLoaded bool

func TestMain(m *testing.M) {
	app.Commands = []*cli.Command{
		MigrateCommands,
	}
	log.SetOutput(ioutil.Discard)
	os.Exit(m.Run())
}

func fakeLoadWithSpy(migrationsPath, environment string) (database.Database, error) {
	dbLoadArgs = []string{migrationsPath, environment}
	fakeDb = internal.FakeDbWithSpy{}
	return &fakeDb, nil
}

func TestMigrateUpDefaults(t *testing.T) {
	dbLoadDb = fakeLoadWithSpy

	args := []string{"sth.exe", "migrate", "up"}
	if err := app.Run(args); err != nil {
		t.Errorf("Error running command - %s", err)
	}

	expected := []string{"./migrations/zlab", "development"}
	if !internal.StrSliceEqual(dbLoadArgs, expected) {
		t.Errorf("Expected to load db with '%v', but got %s", expected, dbLoadArgs)
	}

	fakeDb.AssertWaitForStartCalled(t, true)
	fakeDb.AssertEnsureMigrationsChangelogCalled(t, true)
	fakeDb.AssertEnsureConsistentMigrationsCalled(t, true)
	fakeDb.AssertApplyUpMigrationsWithCountCalledWith(t, 1, false)
	fakeDb.AssertApplySpecificUpMigrationCalled(t, false)
}

func TestMigrateUpWithCount(t *testing.T) {
	dbLoadDb = fakeLoadWithSpy

	args := []string{"sth.exe", "migrate", "up", "--count", "2"}
	if err := app.Run(args); err != nil {
		t.Errorf("Error running command - %s", err)
	}

	expected := []string{"./migrations/zlab", "development"}
	if !internal.StrSliceEqual(dbLoadArgs, expected) {
		t.Errorf("Expected to load db with '%v', but got %s", expected, dbLoadArgs)
	}

	fakeDb.AssertWaitForStartCalled(t, true)
	fakeDb.AssertEnsureMigrationsChangelogCalled(t, true)
	fakeDb.AssertEnsureConsistentMigrationsCalled(t, true)
	fakeDb.AssertApplyUpMigrationsWithCountCalledWith(t, 2, false)
	fakeDb.AssertApplySpecificUpMigrationCalled(t, false)

}

func TestMigrateUpWithAll(t *testing.T) {
	dbLoadDb = fakeLoadWithSpy

	args := []string{"sth.exe", "migrate", "up", "--all"}
	if err := app.Run(args); err != nil {
		t.Errorf("Error running command - %s", err)
	}

	expected := []string{"./migrations/zlab", "development"}
	if !internal.StrSliceEqual(dbLoadArgs, expected) {
		t.Errorf("Expected to load db with '%v', but got %s", expected, dbLoadArgs)
	}

	fakeDb.AssertWaitForStartCalled(t, true)
	fakeDb.AssertEnsureMigrationsChangelogCalled(t, true)
	fakeDb.AssertEnsureConsistentMigrationsCalled(t, true)
	fakeDb.AssertApplyUpMigrationsWithCountCalledWith(t, 0, true)
	fakeDb.AssertApplySpecificUpMigrationCalled(t, false)

}

func TestMigrateUpWithOnly(t *testing.T) {
	dbLoadDb = fakeLoadWithSpy

	args := []string{"sth.exe", "migrate", "up", "--only", "sth"}
	if err := app.Run(args); err != nil {
		t.Errorf("Error running command - %s", err)
	}

	expected := []string{"./migrations/zlab", "development"}
	if !internal.StrSliceEqual(dbLoadArgs, expected) {
		t.Errorf("Expected to load db with '%v', but got %s", expected, dbLoadArgs)
	}

	fakeDb.AssertWaitForStartCalled(t, true)
	fakeDb.AssertEnsureMigrationsChangelogCalled(t, true)
	fakeDb.AssertApplySpecificUpMigrationCalledWith(t, "sth")
	fakeDb.AssertApplyUpMigrationsWithCountCalled(t, false)
	fakeDb.AssertEnsureConsistentMigrationsCalled(t, false)
}

func TestMigrateUpErrorWithMultipleParams(t *testing.T) {
	fakeDbLoaded = false
	dbLoadDb = fakeLoadWithSpy

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
		if fakeDbLoaded {
			t.Errorf("Loaded the database even though it shouldn't")
		}

	}

}
