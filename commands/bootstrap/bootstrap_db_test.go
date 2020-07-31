package bootstrap

import (
	"io/ioutil"
	"log"
	"os"
	"testing"

	"github.com/urfave/cli/v2"

	"go-migrations/database"
	"go-migrations/internal"
)

var dbLoadArgs []string
var fakeDb internal.FakeDbWithSpy

func fakeLoadWithSpy(migrationsPath, environment string) (database.Database, error) {
	dbLoadArgs = []string{migrationsPath, environment}
	fakeDb = internal.FakeDbWithSpy{}
	return &fakeDb, nil
}

var app = cli.NewApp()

func TestMain(m *testing.M) {
	app.Commands = []*cli.Command{
		BootstrapCommand,
	}
	log.SetOutput(ioutil.Discard)
	os.Exit(m.Run())
}

func TestBootstrapDbDefaults(t *testing.T) {
	mockableLoadDB = fakeLoadWithSpy

	args := []string{"sth.exe", "bootstrap"}
	if err := app.Run(args); err != nil {
		t.Errorf("Error running command - %s", err)
	}

	expected := []string{"./migrations/zlab", "development"}
	if !internal.StrSliceEqual(dbLoadArgs, expected) {
		t.Errorf("Expected to load db with '%v', but got %s", expected, dbLoadArgs)
	}

	fakeDb.AssertWaitForStartCalled(t, true)
	fakeDb.AssertEnsureMigrationsChangelogCalled(t, true)
	fakeDb.AssertBootstrapCalled(t, true)
	fakeDb.AssertApplyAllUpMigrationsCalled(t, true)

}

func TestBootstrapDbWithFlags(t *testing.T) {
	mockableLoadDB = fakeLoadWithSpy

	args := []string{"sth.exe", "bootstrap", "-p", "/my/path", "-e", "my-env"}
	if err := app.Run(args); err != nil {
		t.Errorf("Error running command - %s", err)
	}

	expected := []string{"/my/path", "my-env"}
	if !internal.StrSliceEqual(dbLoadArgs, expected) {
		t.Errorf("Expected to load db with '%v', but got %s", expected, dbLoadArgs)
	}

	fakeDb.AssertWaitForStartCalled(t, true)
	fakeDb.AssertEnsureMigrationsChangelogCalled(t, true)
	fakeDb.AssertBootstrapCalled(t, true)
	fakeDb.AssertApplyAllUpMigrationsCalled(t, true)
}
