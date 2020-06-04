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
			CREATE SCHEMA baz;
			CREATE TABLE baz.foo (
				  id VARCHAR(14) NOT NULL PRIMARY KEY
				, name TEXT NOT NULL
				, applied_at timestamptz NOT NULL
			);
			CREATE TABLE baz.bar (
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

	_, err = dbConn.Exec("SELECT id, name, applied_at FROM baz.foo")
	if err != nil {
		t.Errorf("Error checking bootstrap: %v", err)
	}
	_, err = dbConn.Exec("SELECT id FROM baz.bar")
	if err != nil {
		t.Errorf("Error checking bootstrap: %v", err)
	}
}

func TestApplyAllUpMigrations(t *testing.T) {
	cleanup, migrationPath := setupFolder(t)
	defer cleanup()

	bootstrapSQL := []byte(
		dedent.Dedent(`
			CREATE TABLE public.migrations_changelog (
				id VARCHAR(14) NOT NULL PRIMARY KEY
				, name TEXT NOT NULL
				, applied_at timestamptz NOT NULL
			);
		`),
	)

	ioutil.WriteFile(filepath.Join(migrationPath, "bootstrap.sql"), bootstrapSQL, 0777)

	firstMigration := []byte(dedent.Dedent(`
		CREATE TABLE public.fiz (fuz TEXT PRIMARY KEY)
		-- //@UNDO
		DROP TABLE public.fiz
	`))
	ioutil.WriteFile(
		filepath.Join(migrationPath, "_common", "20171101000001_foo.sql"),
		firstMigration,
		0777,
	)
	ioutil.WriteFile(filepath.Join(
		migrationPath, "_common", "verify", "20171101000001_foo.sql"),
		[]byte("SELECT 1"),
		0777,
	)

	secondMigration := []byte(dedent.Dedent(`
		CREATE TABLE public.biz (buz TEXT PRIMARY KEY)
		-- //@UNDO
		DROP TABLE public.biz
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

	if err := db.Bootstrap(); err != nil {
		t.Fatalf("Error during bootstrap: %v", err)
	}

	if err := db.ApplyAllUpMigrations(); err != nil {
		t.Fatalf("Error during up migration: %v", err)
	}

	_, err = dbConn.Exec("SELECT fuz FROM public.fiz")
	if err != nil {
		t.Errorf("Error checking firstMigration: %v", err)
	}
	_, err = dbConn.Exec("SELECT buz FROM public.biz")
	if err != nil {
		t.Errorf("Error checking secondMigration: %v", err)
	}
}

func setupFolder(t *testing.T) (func(), string) {
	dir, err := ioutil.TempDir("", "go_mig")
	if err != nil {
		t.Fatalf("Returned error setting up the tmp directory: %v", err)
	}
	cleanup := func() { os.RemoveAll(dir) }

	os.Mkdir(filepath.Join(dir, "_environments"), 0777)
	os.Mkdir(filepath.Join(dir, "_common"), 0777)
	os.Mkdir(filepath.Join(dir, "_common", "verify"), 0777)
	os.Mkdir(filepath.Join(dir, "_common", "prepare"), 0777)

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
