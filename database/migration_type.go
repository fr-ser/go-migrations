package database

import (
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"
	"regexp"
	"strings"
	"time"
)

// FileMigration is a struct around a local migration with all attached SQL files
type FileMigration struct {
	UpSQL       string
	DownSQL     string
	VerifySQL   string
	ID          string
	Description string
	Filename    string
	Application string
}

// LoadFromFile loads all properties based on the filepath of the migration itself
func (mig *FileMigration) LoadFromFile(migrationPath string) error {
	mig.Application = filepath.Base(filepath.Dir(migrationPath))
	mig.Filename = filepath.Base(migrationPath)

	validName := regexp.MustCompile(`^\d{14}_[\w]+\.sql$`).Match([]byte(mig.Filename))
	if validName == false {
		return fmt.Errorf("The migration file at '%s' was invalid", migrationPath)

	}

	idMatch := regexp.MustCompile(`(^\d+)_`).FindStringSubmatch(mig.Filename)
	if len(idMatch) < 2 {
		return fmt.Errorf("Could not find id in filename '%s'", mig.Filename)
	}
	mig.ID = idMatch[1]

	runes := []rune(mig.Filename)
	mig.Description = string(runes[15 : len(runes)-4])

	if err := mig.loadMigration(migrationPath); err != nil {
		return err
	}
	if err := mig.loadVerify(migrationPath); err != nil {
		return err
	}

	return nil
}

func (mig *FileMigration) loadMigration(migrationPath string) error {

	migrationFile, err := os.Open(migrationPath)
	if err != nil {
		return fmt.Errorf("Couldn't open migration file: %v", err)
	}
	defer migrationFile.Close()
	migration, err := ioutil.ReadAll(migrationFile)
	if err != nil {
		return fmt.Errorf("Couldn't read migration file: %v", err)
	}

	if string(migration) == "" {
		return fmt.Errorf("The migration at '%s' was empty", mig.Filename)
	}

	UpDownMigration := strings.Split(string(migration), "\n-- //@UNDO\n")
	if len(UpDownMigration) != 2 {
		return fmt.Errorf("Could not find up and down migration in '%s'", mig.Filename)
	}

	mig.UpSQL = strings.Trim(strings.Trim(UpDownMigration[0], "\n"), " ")
	if mig.UpSQL == "" {
		return fmt.Errorf("The up migration at '%s' was empty", mig.Filename)
	}

	mig.DownSQL = strings.Trim(strings.Trim(UpDownMigration[1], "\n"), " ")
	if mig.DownSQL == "" {
		return fmt.Errorf("The down migration at '%s' was empty", mig.Filename)
	}

	return nil
}

func (mig *FileMigration) loadVerify(migrationPath string) error {
	verifyPath := filepath.Join(
		filepath.Dir(migrationPath), "verify", filepath.Base(migrationPath),
	)
	verifyFile, err := os.Open(verifyPath)
	if err != nil {
		return fmt.Errorf("Couldn't open verify file: %v", err)
	}
	defer verifyFile.Close()
	verify, err := ioutil.ReadAll(verifyFile)
	if err != nil {
		return fmt.Errorf("Couldn't read verify file: %v", err)
	}
	mig.VerifySQL = strings.Trim(string(verify), "\n")
	if mig.VerifySQL == "" {
		return fmt.Errorf("Verify file for %s was empty", migrationPath)
	}

	return nil
}

// AppliedMigration is a struct around a migration in the database / changelog table
type AppliedMigration struct {
	ID        string
	Name      string
	AppliedAt time.Time
}
