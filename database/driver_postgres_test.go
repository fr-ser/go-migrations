package database

import (
	"database/sql"
	"io/ioutil"
	"os"
	"path/filepath"
	"testing"

	"github.com/lithammer/dedent"
)

var dbConn, _ = sql.Open("pgx", "postgresql://admin:admin_pass@localhost:35432/my_db")

func TestApplyBootstrap(t *testing.T) {
	cleanup, migrationPath := setupFolder(t)
	defer cleanup()

	bootstrapSQL := []byte(
		dedent.Dedent(`
			-- Some SQL comment
			CREATE SCHEMA db_main;
			CREATE TABLE db_main.migrations_changelog (
				  id VARCHAR(14) NOT NULL PRIMARY KEY
				, name TEXT NOT NULL
				, applied_at timestamptz NOT NULL
			);
			CREATE TABLE db_main.table_2 (
				id VARCHAR(14) NOT NULL PRIMARY KEY
		  	);
		`),
	)
	ioutil.WriteFile(filepath.Join(migrationPath, "bootstrap.sql"), bootstrapSQL, 0777)

	db, err := LoadDb(migrationPath, "development")
	if err != nil {
		t.Fatalf("Returned error loading database: %v", err)
	}
	err = db.Bootstrap()
	if err != nil {
		t.Fatalf("Error during bootstrap: %v", err)
	}

	_, err = dbConn.Exec("SELECT id, name, applied_at FROM db_main.migrations_changelog")
	if err != nil {
		t.Errorf("Error checking bootstrap: %v", err)
	}
	_, err = dbConn.Exec("SELECT id FROM db_main.table_2")
	if err != nil {
		t.Errorf("Error checking bootstrap: %v", err)
	}
}

func setupFolder(t *testing.T) (func(), string) {
	dir, err := ioutil.TempDir("", "go_mig")
	if err != nil {
		t.Fatalf("Returned error setting up the tmp directory: %v", err)
	}
	cleanup := func() { os.RemoveAll(dir) }

	err = os.Mkdir(filepath.Join(dir, "_environments"), 0777)
	if err != nil {
		t.Fatalf("Returned error creating env/config folder: %v", err)
		cleanup()
	}

	defaultConfig := dedent.Dedent(`
		db_type: postgres
		prepare: True
		host: localhost
		port: 35432
		db_name: my_db
		user: admin
		password: admin_pass
	`)
	err = ioutil.WriteFile(
		filepath.Join(dir, "_environments", "development.yaml"),
		[]byte(defaultConfig),
		0777,
	)
	if err != nil {
		t.Fatalf("Returned error creating default env: %v", err)
		cleanup()
	}

	return cleanup, dir
}
