// +build !unit

package postgres_test

import (
	"database/sql"
	"io/ioutil"
	"os"
	"path/filepath"
	"testing"

	"github.com/lithammer/dedent"

	"go-migrations/database/driver"
)

var dbConn, _ = sql.Open("pgx", "postgresql://admin:admin_pass@localhost:35432/my_db")

func TestApplyBootstrap(t *testing.T) {
	cleanup, migrationPath := setupFolder(t)
	defer cleanup()

	bootstrapSQL := []byte(
		dedent.Dedent(`
			-- Some SQL comment
			CREATE SCHEMA boot_baz;
			CREATE TABLE boot_baz.foo (
				  id VARCHAR(14) NOT NULL PRIMARY KEY
				, name TEXT NOT NULL
				, applied_at timestamptz NOT NULL
			);
			CREATE TABLE boot_baz.bar (
				id VARCHAR(14) NOT NULL PRIMARY KEY
		  	);
		`),
	)
	ioutil.WriteFile(filepath.Join(migrationPath, "bootstrap.sql"), bootstrapSQL, 0777)

	db, err := driver.LoadDb(migrationPath, "development")
	if err != nil {
		t.Fatalf("Returned error loading database: %v", err)
	}
	err = db.Bootstrap()
	if err != nil {
		t.Fatalf("Error during bootstrap: %v", err)
	}

	_, err = dbConn.Exec("SELECT id, name, applied_at FROM boot_baz.foo")
	if err != nil {
		t.Errorf("Error checking bootstrap: %v", err)
	}
	_, err = dbConn.Exec("SELECT id FROM boot_baz.bar")
	if err != nil {
		t.Errorf("Error checking bootstrap: %v", err)
	}
}

func TestApplyAllUpMigrations(t *testing.T) {
	cleanupFileMigrations, migrationPath := setupFolder(t)
	defer cleanupFileMigrations()
	defer cleanupChangelog()

	firstMigration := []byte(dedent.Dedent(`
		CREATE TABLE public.all_fiz (fuz TEXT PRIMARY KEY)
		-- //@UNDO
		DROP TABLE public.all_fiz
	`))
	ioutil.WriteFile(
		filepath.Join(migrationPath, "common", "20171101000001_foo.sql"),
		firstMigration,
		0777,
	)
	ioutil.WriteFile(filepath.Join(
		migrationPath, "common", "verify", "20171101000001_foo.sql"),
		[]byte("SELECT 1"),
		0777,
	)

	secondMigration := []byte(dedent.Dedent(`
		CREATE TABLE public.all_biz (buz TEXT PRIMARY KEY)
		-- //@UNDO
		DROP TABLE public.all_biz
	`))
	os.Mkdir(filepath.Join(migrationPath, "analytics"), 0777)
	ioutil.WriteFile(
		filepath.Join(migrationPath, "analytics", "20171101000002_bar.sql"),
		secondMigration,
		0777,
	)
	os.Mkdir(filepath.Join(migrationPath, "analytics", "verify"), 0777)
	ioutil.WriteFile(filepath.Join(
		migrationPath, "analytics", "verify", "20171101000002_bar.sql"),
		[]byte("SELECT 2"),
		0777,
	)

	db, err := driver.LoadDb(migrationPath, "development")
	if err != nil {
		t.Fatalf("Returned error loading database: %v", err)
	}

	if _, err := db.EnsureMigrationsChangelog(); err != nil {
		t.Fatalf("Error during changelog creation: %v", err)
	}

	if err := db.ApplyAllUpMigrations(); err != nil {
		t.Fatalf("Error during up migration: %v", err)
	}

	_, err = dbConn.Exec("SELECT fuz FROM public.all_fiz")
	if err != nil {
		t.Errorf("Error checking firstMigration: %v", err)
	}
	_, err = dbConn.Exec("SELECT buz FROM public.all_biz")
	if err != nil {
		t.Errorf("Error checking secondMigration: %v", err)
	}
}

func TestApplyUpMigrationsWithCount(t *testing.T) {
	cleanupFileMigrations, migrationPath := setupFolder(t)
	defer cleanupFileMigrations()
	defer cleanupChangelog()

	firstMigration := []byte(dedent.Dedent(`
		CREATE TABLE public.count_foo (fuz TEXT PRIMARY KEY)
		-- //@UNDO
		DROP TABLE public.count_foo
	`))
	ioutil.WriteFile(
		filepath.Join(migrationPath, "common", "20171101000001_foo.sql"),
		firstMigration,
		0777,
	)
	ioutil.WriteFile(filepath.Join(
		migrationPath, "common", "verify", "20171101000001_foo.sql"),
		[]byte("SELECT 1"),
		0777,
	)

	secondMigration := []byte(dedent.Dedent(`
		INSERT INTO public.count_foo (fuz) VALUES ('one');
		-- //@UNDO
		DELETE FROM public.count_foo WHERE fuz = 'one'
	`))
	os.Mkdir(filepath.Join(migrationPath, "analytics"), 0777)
	ioutil.WriteFile(
		filepath.Join(migrationPath, "analytics", "20171101000002_bar.sql"),
		secondMigration,
		0777,
	)
	os.Mkdir(filepath.Join(migrationPath, "analytics", "verify"), 0777)
	ioutil.WriteFile(filepath.Join(
		migrationPath, "analytics", "verify", "20171101000002_bar.sql"),
		[]byte("SELECT 2"),
		0777,
	)

	thirdMigration := []byte(dedent.Dedent(`
		INSERT INTO public.count_foo (fuz) VALUES ('two');
		-- //@UNDO
		DELETE FROM public.count_foo WHERE fuz = 'two'
	`))
	os.Mkdir(filepath.Join(migrationPath, "analytics"), 0777)
	ioutil.WriteFile(
		filepath.Join(migrationPath, "analytics", "20171101000003_buz.sql"),
		thirdMigration,
		0777,
	)
	os.Mkdir(filepath.Join(migrationPath, "analytics", "verify"), 0777)
	ioutil.WriteFile(filepath.Join(
		migrationPath, "analytics", "verify", "20171101000003_buz.sql"),
		[]byte("SELECT 2"),
		0777,
	)

	db, err := driver.LoadDb(migrationPath, "development")
	if err != nil {
		t.Fatalf("Returned error loading database: %v", err)
	}

	if _, err := db.EnsureMigrationsChangelog(); err != nil {
		t.Fatalf("Error during changelog creation: %v", err)
	}

	if err := db.ApplyUpMigrationsWithCount(3, false); err != nil {
		t.Fatalf("Error during up migration: %v", err)
	}

	var rowCount int
	verifyCount := dbConn.QueryRow("SELECT COUNT(*) FROM public.count_foo")
	scanErr := verifyCount.Scan(&rowCount)
	if rowCount != 2 {
		t.Fatalf(
			dedent.Dedent(`
				Expected rowCount of %d, but got %d. Incorrect up migration.
				ScanErr: %v
			`),
			2, rowCount, scanErr,
		)
	}
}

func cleanupChangelog() {
	dbConn.Exec(`TRUNCATE public.migrations_changelog`)
}

func setupFolder(t *testing.T) (func(), string) {
	dir, err := ioutil.TempDir("", "go_mig")
	if err != nil {
		t.Fatalf("Returned error setting up the tmp directory: %v", err)
	}
	cleanup := func() { os.RemoveAll(dir) }

	os.Mkdir(filepath.Join(dir, "_environments"), 0777)
	os.Mkdir(filepath.Join(dir, "common"), 0777)
	os.Mkdir(filepath.Join(dir, "common", "verify"), 0777)
	os.Mkdir(filepath.Join(dir, "common", "prepare"), 0777)

	defaultConfig := dedent.Dedent(`
		db_type: postgres
		prepare: True
		host: localhost
		port: 35432
		db_name: my_db
		user: admin
		password: admin_pass
	`)
	ioutil.WriteFile(
		filepath.Join(dir, "_environments", "development.yaml"),
		[]byte(defaultConfig),
		0777,
	)

	return cleanup, dir
}
