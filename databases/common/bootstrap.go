package common

import (
	"database/sql"
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"
)

// ApplyBootstrapMigration applies the bootstrap.sql, which it finds by itself based
// on the migrations path
func ApplyBootstrapMigration(db *sql.DB, migrationsPath string) error {
	bootstrapFile, err := os.Open(filepath.Join(migrationsPath, "bootstrap.sql"))
	if err != nil {
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
