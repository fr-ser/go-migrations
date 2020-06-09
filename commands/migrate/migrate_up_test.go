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

	fakeDb.WaitForStartCalled(t)
	fakeDb.EnsureMigrationsChangelogCalled(t)
	fakeDb.ApplyUpMigrationsWithCountCalledWith(t, 1, false)
	fakeDb.ApplySpecificUpMigrationNotCalled(t)
}

func TestMigrateUpWithCountAndAll(t *testing.T) {
	dbLoadDb = fakeLoadWithSpy

	args := []string{"sth.exe", "migrate", "up", "--count", "2", "--all"}
	if err := app.Run(args); err != nil {
		t.Errorf("Error running command - %s", err)
	}

	expected := []string{"./migrations/zlab", "development"}
	if !internal.StrSliceEqual(dbLoadArgs, expected) {
		t.Errorf("Expected to load db with '%v', but got %s", expected, dbLoadArgs)
	}

	fakeDb.WaitForStartCalled(t)
	fakeDb.EnsureMigrationsChangelogCalled(t)
	fakeDb.ApplyUpMigrationsWithCountCalledWith(t, 2, true)
	fakeDb.ApplySpecificUpMigrationNotCalled(t)

}

func TestMigrateUpWithOnlyOverCountAndAll(t *testing.T) {
	dbLoadDb = fakeLoadWithSpy

	args := []string{"sth.exe", "migrate", "up", "--count", "2", "--all", "--only", "sth"}
	if err := app.Run(args); err != nil {
		t.Errorf("Error running command - %s", err)
	}

	expected := []string{"./migrations/zlab", "development"}
	if !internal.StrSliceEqual(dbLoadArgs, expected) {
		t.Errorf("Expected to load db with '%v', but got %s", expected, dbLoadArgs)
	}

	fakeDb.WaitForStartCalled(t)
	fakeDb.EnsureMigrationsChangelogCalled(t)
	fakeDb.ApplyUpMigrationsWithCountNotCalled(t)
	fakeDb.ApplySpecificUpMigrationCalledWith(t, "sth")
}
