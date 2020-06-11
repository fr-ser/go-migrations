package database

import (
	"io/ioutil"
	"os"
	"path/filepath"
	"testing"

	"github.com/kylelemons/godebug/pretty"
)

func TestGetFileMigrations(t *testing.T) {
	basePath, err := ioutil.TempDir("", "go_mig")
	if err != nil {
		t.Fatalf("Returned error setting up the tmp directory: %v", err)
	}
	defer os.RemoveAll(basePath)

	expectedMigrations := []FileMigration{
		saveMigrationFor(basePath, "z_my_app", "20171101000001_foo.sql"),
		saveMigrationFor(basePath, "a_other_app", "20171101000002_bar.sql"),
	}

	gotMigrations, err := GetFileMigrations(basePath)
	if err != nil {
		t.Fatalf("Got an error loading migrations: %v", err)
	}

	diff := pretty.Compare(expectedMigrations, gotMigrations)
	if diff != "" {
		t.Error(diff)
	}
}

func TestIgnoreFilesWithoutApps(t *testing.T) {
	basePath, err := ioutil.TempDir("", "go_mig")
	if err != nil {
		t.Fatalf("Returned error setting up the tmp directory: %v", err)
	}
	defer os.RemoveAll(basePath)

	ioutil.WriteFile(filepath.Join(basePath, "foo.sql"), []byte("foo"), 0777)

	gotMigrations, err := GetFileMigrations(basePath)
	if err != nil {
		t.Fatalf("Got an error loading migrations: %v", err)
	}

	if len(gotMigrations) > 0 {
		t.Errorf("Expected no migrations, but got %d", len(gotMigrations))
	}
}
func TestIgnoreFoldersInApps(t *testing.T) {
	basePath, err := ioutil.TempDir("", "go_mig")
	if err != nil {
		t.Fatalf("Returned error setting up the tmp directory: %v", err)
	}
	defer os.RemoveAll(basePath)

	os.Mkdir(filepath.Join(basePath, "some_app"), 0777)
	os.Mkdir(filepath.Join(basePath, "some_app", "some_folder"), 0777)
	ioutil.WriteFile(
		filepath.Join(basePath, "some_app", "some_folder", "bar.sql"), []byte("bar"), 0777,
	)

	gotMigrations, err := GetFileMigrations(basePath)
	if err != nil {
		t.Fatalf("Got an error loading migrations: %v", err)
	}

	if len(gotMigrations) > 0 {
		t.Errorf("Expected no migrations, but got %d", len(gotMigrations))
	}
}

func TestIgnoreEnvironments(t *testing.T) {
	basePath, err := ioutil.TempDir("", "go_mig")
	if err != nil {
		t.Fatalf("Returned error setting up the tmp directory: %v", err)
	}
	defer os.RemoveAll(basePath)

	os.Mkdir(filepath.Join(basePath, "_environments"), 0777)
	ioutil.WriteFile(filepath.Join(basePath, "_environments", "some_env.yaml"), []byte("1"), 0777)

	gotMigrations, err := GetFileMigrations(basePath)
	if err != nil {
		t.Fatalf("Got an error loading migrations: %v", err)
	}

	if len(gotMigrations) > 0 {
		t.Errorf("Expected no migrations, but got %d", len(gotMigrations))
	}
}

func TestNoDuplicateIDs(t *testing.T) {
	basePath, err := ioutil.TempDir("", "go_mig")
	if err != nil {
		t.Fatalf("Returned error setting up the tmp directory: %v", err)
	}
	defer os.RemoveAll(basePath)

	saveMigrationFor(basePath, "my_app", "20171101000001_foo.sql")
	saveMigrationFor(basePath, "other_app", "20171101000001_bar.sql")

	_, err = GetFileMigrations(basePath)
	if err == nil {
		t.Fatal("Got no error with duplicate IDs")
	}
}

func saveMigrationFor(basePath, application, migrationName string) FileMigration {
	os.Mkdir(filepath.Join(basePath, application), 0777)
	os.Mkdir(filepath.Join(basePath, application, "verify"), 0777)
	os.Mkdir(filepath.Join(basePath, application, "prepare"), 0777)

	appPath := filepath.Join(basePath, application)

	migrationSQL := []byte("CREATE SCHEMA template;\n-- //@UNDO\nDROP SCHEMA template;")
	ioutil.WriteFile(filepath.Join(appPath, migrationName), migrationSQL, 0777)

	verifySQL := []byte("SELECT 1")
	ioutil.WriteFile(filepath.Join(appPath, "verify", migrationName), verifySQL, 0777)

	prepareSQL := []byte("SELECT 2")
	ioutil.WriteFile(filepath.Join(appPath, "prepare", migrationName), prepareSQL, 0777)

	return FileMigration{
		Filename:    migrationName,
		ID:          migrationName[0:14],
		Description: migrationName[15 : len(migrationName)-4],
		Application: application,
		UpSQL:       "CREATE SCHEMA template;",
		DownSQL:     "DROP SCHEMA template;",
		PrepareSQL:  "SELECT 2",
		VerifySQL:   "SELECT 1",
	}
}
