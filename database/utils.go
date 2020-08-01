package database

import (
	"database/sql"
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"
	"time"
)

// WaitForStart tries to connect to the database
// parameters are the number of retries and the sleep interval in milliseconds between the retries
func WaitForStart(db *sql.DB, pollInterval time.Duration, retries int) error {
	var err error

	for retry := 0; retry < retries; retry++ {
		_, err = db.Exec("SELECT 1")
		if err == nil {
			return nil
		}
		time.Sleep(pollInterval)
	}

	return fmt.Errorf("Timed out connecting to database: %v", err)
}

// GetBootstrapSQL returns the SQL string of the bootstrap file
// or returns an empty string if the file does not exist
func GetBootstrapSQL(migrationsPath string) (sql string, err error) {
	bootstrapFile, err := os.Open(filepath.Join(migrationsPath, "bootstrap.sql"))
	if err != nil {
		if os.IsNotExist(err) {
			return "", nil
		}

		return "", fmt.Errorf("Couldn't open bootstrap file: %v", err)
	}
	defer bootstrapFile.Close()
	fileContent, err := ioutil.ReadAll(bootstrapFile)
	if err != nil {
		return "", fmt.Errorf("Couldn't read config file: %v", err)
	}

	return string(fileContent), nil
}

// ApplyBootstrapMigration applies the bootstrap.sql, which it finds by itself based
// on the migrations path
func ApplyBootstrapMigration(db *sql.DB, migrationsPath string) (err error) {
	fileContent, err := GetBootstrapSQL(migrationsPath)
	if err != nil {
		return err
	}
	if fileContent == "" {
		return nil
	}
	_, err = db.Exec(string(fileContent))
	if err != nil {
		return fmt.Errorf("Could not apply bootstrap.sql: %v", err)
	}

	return nil
}

// EnsureConsistentMigrations checks if all applied migrations (by ID) exist as local files
// and if no local migration has been "skipped" (newer migrations applied)
func EnsureConsistentMigrations(fileMigrations []FileMigration, appliedMigrations []AppliedMigration) error {
	for idx := 0; idx < len(appliedMigrations); idx++ {
		if len(fileMigrations) <= idx || fileMigrations[idx].ID != appliedMigrations[idx].ID {
			moreInfo := "For more information execute the migrate status command"
			if idx > 0 {
				return fmt.Errorf(
					"FileMigrations and AppliedMigrations are out of sync after %s\n%s",
					fileMigrations[idx-1].ID, moreInfo,
				)
			}
			return fmt.Errorf(
				"Local and applied migrations are out of sync already at the first migration.\n%s",
				moreInfo,
			)
		}
	}
	return nil
}
