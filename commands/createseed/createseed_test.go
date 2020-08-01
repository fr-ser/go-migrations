package createseed

import (
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"testing"
	"time"

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
		CreateSeedCommand,
	}
	log.SetOutput(ioutil.Discard)
	os.Exit(m.Run())
}

func TestCreateSeedFlags(t *testing.T) {
	mockableLoadDB = fakeLoadWithSpy

	uniqueFilename := fmt.Sprintf("%s/go-mig-test-%d", os.TempDir(), time.Now().UnixNano())

	args := []string{
		"sth.exe", "create-seed",
		"-p", "/my/path", "-e", "my-env", "-t", uniqueFilename,
	}
	if err := app.Run(args); err != nil {
		t.Errorf("Error running command - %s", err)
	}

	expected := []string{"/my/path", "my-env"}
	if !internal.StrSliceEqual(dbLoadArgs, expected) {
		t.Errorf("Expected to load db with '%v', but got %s", expected, dbLoadArgs)
	}

	fakeDb.AssertGenerateSeedSQLCalled(t, true)
}

func TestErrorForExistingTargets(t *testing.T) {
	fakeDb = internal.FakeDbWithSpy{}
	mockableLoadDB = fakeLoadWithSpy

	target, _ := ioutil.TempFile(os.TempDir(), "g-mig-test-")
	defer os.Remove(target.Name())

	args := []string{"sth.exe", "create-seed", "-t", target.Name()}
	err := app.Run(args)

	expectedMsg := fmt.Sprintf("The file %s already exists", target.Name())
	if err == nil || err.Error() != expectedMsg {
		t.Errorf("Expected a file already exists error but got - %s", err)
	}

	fakeDb.AssertGenerateSeedSQLCalled(t, false)
}
