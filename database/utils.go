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

// ApplyBootstrapMigration applies the bootstrap.sql, which it finds by itself based
// on the migrations path
func ApplyBootstrapMigration(db *sql.DB, migrationsPath string) error {
	bootstrapFile, err := os.Open(filepath.Join(migrationsPath, "bootstrap.sql"))
	if err != nil {
		if os.IsNotExist(err) {
			return nil
		}

		return fmt.Errorf("Couldn't open bootstrap file: %v", err)
	}
	defer bootstrapFile.Close()
	fileContent, err := ioutil.ReadAll(bootstrapFile)
	if err != nil {
		return fmt.Errorf("Couldn't read config file: %v", err)
	}
	_, err = db.Exec(string(fileContent))
	return err
}
