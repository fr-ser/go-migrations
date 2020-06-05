package start

import (
	"os/exec"
	"testing"

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

func fakeRunWithOutput(cmd *exec.Cmd) (stdout, stderr string, err error) {
	return "", "", nil
}

func assertDbCalls(t *testing.T, db internal.FakeDbWithSpy) {
	if !db.WaitForStartCalled {
		t.Error("WaitForStart not called")
	}
	if !db.BootstrapCalled {
		t.Error("Bootstrap not called")
	}
	if !db.ApplyAllUpMigrationsCalled {
		t.Error("ApplyAllUpMigrations not called")
	}
	if !db.EnsureMigrationsChangelogCalled {
		t.Error("EnsureMigrationsChangelog not called")
	}
}

func TestStartDbDefaults(t *testing.T) {
	runWithOutput = fakeRunWithOutput
	dbLoadDb = fakeLoadWithSpy

	args := []string{"sth.exe", "start"}
	if err := app.Run(args); err != nil {
		t.Errorf("Error running command - %s", err)
	}

	expected := []string{"./migrations/zlab", "development"}
	if !internal.StrSliceEqual(dbLoadArgs, expected) {
		t.Errorf("Expected to load db with '%v', but got %s", expected, dbLoadArgs)
	}

	assertDbCalls(t, fakeDb)

}

func TestStartDbWithFlags(t *testing.T) {
	runWithOutput = fakeRunWithOutput
	dbLoadDb = fakeLoadWithSpy

	args := []string{"sth.exe", "start", "-p", "/my/path", "-e", "my-env"}
	if err := app.Run(args); err != nil {
		t.Errorf("Error running command - %s", err)
	}

	expected := []string{"/my/path", "my-env"}
	if !internal.StrSliceEqual(dbLoadArgs, expected) {
		t.Errorf("Expected to load db with '%v', but got %s", expected, dbLoadArgs)
	}

	assertDbCalls(t, fakeDb)
}
