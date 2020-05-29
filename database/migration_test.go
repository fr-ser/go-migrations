package database

import (
	"io/ioutil"
	"os"
	"path/filepath"
	"testing"

	"github.com/kylelemons/godebug/pretty"
	"github.com/lithammer/dedent"
)

func TestLoadMigration(t *testing.T) {
	cleanup, migrationPath := setupFolder(t)
	defer cleanup()

	filename := "20171101000001_foo.sql"

	migrationSQL := []byte(dedent.Dedent(`
		CREATE SCHEMA foo;
		-- //@UNDO
		DROP SCHEMA foo;
	`))
	ioutil.WriteFile(
		filepath.Join(migrationPath, "_common", filename),
		migrationSQL, 0777,
	)

	verifySQL := []byte("SELECT 1")
	ioutil.WriteFile(
		filepath.Join(migrationPath, "_common", "verify", filename),
		verifySQL, 0777,
	)

	prepareSQL := []byte("SELECT 2")
	ioutil.WriteFile(
		filepath.Join(migrationPath, "_common", "prepare", filename),
		prepareSQL, 0777,
	)

	migration := FileMigration{}
	err := migration.LoadFromFile(filepath.Join(migrationPath, "_common", filename))
	if err != nil {
		t.Errorf("Returned error loading migration: %v", err)
	}

	expectedMigration := FileMigration{
		UpSQL:       "CREATE SCHEMA foo;",
		DownSQL:     "DROP SCHEMA foo;",
		VerifySQL:   "SELECT 1",
		PrepareSQL:  "SELECT 2",
		ID:          "20171101000001",
		Filename:    filename,
		Application: "_common",
	}

	if diff := pretty.Compare(migration, expectedMigration); diff != "" {
		t.Errorf("The migration was not as expected:\n%s", diff)
	}
}

func TestInvalidFilenames(t *testing.T) {
	var filenames = []string{
		"foo.sql",
		"20171101000001_foo",
		"20171101000001_foo.pdf",
		"2017_11_01_00_00_01_foo.sql",
		"2017-11-01_00:00:01_foo.sql",
	}
	for _, filename := range filenames {
		t.Run(filename, func(t *testing.T) {
			appPath, cleanup := setupMigrationFor(t, filename)
			defer cleanup()

			failingMigration := FileMigration{}
			err := failingMigration.LoadFromFile(filepath.Join(appPath, filename))
			if err == nil {
				t.Errorf("Did not get an error for the invalid fileName '%s'", filename)
			}
		})
	}
}
func TestRequireVerifySqlFile(t *testing.T) {
	filename := "20171101000001_foo.sql"
	appPath, cleanup := setupMigrationFor(t, filename)
	defer cleanup()

	os.Remove(filepath.Join(appPath, "verify", filename))

	failingMigration := FileMigration{}
	err := failingMigration.LoadFromFile(filepath.Join(appPath, filename))
	if err == nil {
		t.Error("Did not get an error for a missing verify script")
	}
}

func TestRequireVerifySqlContent(t *testing.T) {
	filename := "20171101000001_foo.sql"
	appPath, cleanup := setupMigrationFor(t, filename)
	defer cleanup()

	file, _ := os.OpenFile(filepath.Join(appPath, "verify", filename), os.O_RDWR, 0666)
	file.Truncate(0)

	failingMigration := FileMigration{}
	err := failingMigration.LoadFromFile(filepath.Join(appPath, filename))
	if err == nil {
		t.Error("Did not get an error for an empty verify script")
	}
}

func TestOptionalPrepareSQL(t *testing.T) {
	filename := "20171101000001_foo.sql"
	appPath, cleanup := setupMigrationFor(t, filename)
	defer cleanup()

	os.Remove(filepath.Join(appPath, "prepare", filename))

	migration := FileMigration{}
	err := migration.LoadFromFile(filepath.Join(appPath, filename))
	if err != nil {
		t.Errorf("Did get an error for a missing prepare script: %v", err)
	}
	if migration.PrepareSQL != "" {
		t.Errorf("Expected prepareSQL of '', but got '%s'", migration.PrepareSQL)
	}
}

func TestRequireMigrationContent(t *testing.T) {
	var invalidContents = []struct{ name, migration string }{
		{"empty migration", ""},
		{"missing up", "\n-- //@UNDO\nDROP SCHEMA template;"},
		{"missing down", "CREATE SCHEMA template;\n-- //@UNDO\n"},
	}
	for _, content := range invalidContents {
		t.Run(content.name, func(t *testing.T) {
			filename := "20171101000001_foo.sql"
			appPath, cleanup := setupMigrationFor(t, filename)
			defer cleanup()

			file, _ := os.OpenFile(filepath.Join(appPath, filename), os.O_RDWR, 0666)
			file.Truncate(0)
			ioutil.WriteFile(filepath.Join(appPath, filename), []byte(content.migration), 0777)

			failingMigration := FileMigration{}
			err := failingMigration.LoadFromFile(filepath.Join(appPath, filename))
			if err == nil {
				t.Errorf("Did not get an error for an invalid migration: %s", content.name)
			}
		})
	}
}

func setupFolder(t *testing.T) (func(), string) {
	dir, err := ioutil.TempDir("", "go_mig")
	if err != nil {
		t.Fatalf("Returned error setting up the tmp directory: %v", err)
	}
	cleanup := func() { os.RemoveAll(dir) }

	os.Mkdir(filepath.Join(dir, "_common"), 0777)
	os.Mkdir(filepath.Join(dir, "_common", "verify"), 0777)
	os.Mkdir(filepath.Join(dir, "_common", "prepare"), 0777)

	return cleanup, dir
}

func setupMigrationFor(t *testing.T, migrationName string) (
	applicationPath string, cleanup func(),
) {
	dir, err := ioutil.TempDir("", "go_mig")
	if err != nil {
		t.Fatalf("Returned error setting up the tmp directory: %v", err)
	}
	cleanup = func() { os.RemoveAll(dir) }

	os.Mkdir(filepath.Join(dir, "_common"), 0777)
	os.Mkdir(filepath.Join(dir, "_common", "verify"), 0777)
	os.Mkdir(filepath.Join(dir, "_common", "prepare"), 0777)

	basePath := filepath.Join(dir, "_common")

	migrationSQL := []byte("CREATE SCHEMA template;\n-- //@UNDO\nDROP SCHEMA template;")
	ioutil.WriteFile(filepath.Join(basePath, migrationName), migrationSQL, 0777)

	verifySQL := []byte("SELECT 1")
	ioutil.WriteFile(filepath.Join(basePath, "verify", migrationName), verifySQL, 0777)

	prepareSQL := []byte("SELECT 2")
	ioutil.WriteFile(filepath.Join(basePath, "prepare", migrationName), prepareSQL, 0777)

	return basePath, cleanup
}
