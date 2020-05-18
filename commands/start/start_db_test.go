package start

import (
	"os/exec"
	"testing"

	"go-migrations/databases"
	"go-migrations/databases/config"
	"go-migrations/internal"
)

type fakeDbWithSpy struct {
	// WaitForStartCalled tracks whether the method was called
	WaitForStartCalled bool
	// WaitForStartCalled tracks whether the method was called
	BootstrapCalled bool
	// WaitForStartCalled tracks whether the method was called
	ApplyUpMigrationsCalled bool
}

func (db *fakeDbWithSpy) WaitForStart() error {
	db.WaitForStartCalled = true
	return nil
}
func (db *fakeDbWithSpy) Bootstrap() error {
	db.BootstrapCalled = true
	return nil
}
func (db *fakeDbWithSpy) ApplyUpMigrations() error {
	db.ApplyUpMigrationsCalled = true
	return nil
}

func (db *fakeDbWithSpy) Init(_ config.Config) error {
	return nil
}

var dbLoadArgs []string
var fakeDb fakeDbWithSpy

func fakeLoadWithSpy(migrationsPath, environment string) (databases.Database, error) {
	dbLoadArgs = []string{migrationsPath, environment}
	fakeDb = fakeDbWithSpy{}
	return &fakeDb, nil
}

func fakeRunWithOutput(cmd *exec.Cmd) (stdout, stderr string, err error) {
	return "", "", nil
}

func assertDbCalls(t *testing.T, db fakeDbWithSpy) {
	if !db.WaitForStartCalled {
		t.Error("WaitForStart not called")
	}
	if !db.BootstrapCalled {
		t.Error("Bootstrap not called")
	}
	if !db.ApplyUpMigrationsCalled {
		t.Error("ApplyUpMigrations not called")
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
