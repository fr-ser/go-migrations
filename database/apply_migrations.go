package database

import (
	"database/sql"
	"fmt"

	"github.com/lithammer/dedent"
)

// ApplyUpMigration is an internal helper to apply the up migration in a transaction
// it does not perform anything else (like verify execution)
func ApplyUpMigration(db *sql.DB, migration FileMigration) error {
	upTx, err := db.Begin()
	if err != nil {
		return fmt.Errorf("Error opening transaction: %v", err)
	}

	_, err = upTx.Exec(migration.UpSQL)
	if err != nil {
		rollbackError := upTx.Rollback()
		if rollbackError != nil {
			return fmt.Errorf(
				"Error during up migration of %s: %s \n and rollback error: %s",
				migration.Filename,
				err,
				rollbackError,
			)
		}
		return fmt.Errorf("Error during up migration of %s: %s", migration.Filename, err)
	}

	err = upTx.Commit()
	if err != nil {
		return fmt.Errorf("Error during commit of %s: %s", migration.Filename, err)
	}
	return nil
}

// InsertToChangelog is an internal helper to insert the migration into the changelog
// The table name is passed to allow specifying database specific paths and schemas
func InsertToChangelog(db *sql.DB, migration FileMigration, changelogTable string) error {
	_, err := db.Exec(fmt.Sprintf(
		"INSERT INTO public.migrations_changelog(id, name, applied_at) VALUES ('%s', '%s', now())",
		migration.ID,
		migration.Description,
	))
	if err != nil {
		return fmt.Errorf(
			"Could not update migration changelog for %s: %v",
			migration.Filename,
			err,
		)
	}
	return nil
}

// ApplyVerify is an internal helper to apply the verify script in a transaction and roll it back
func ApplyVerify(db *sql.DB, migration FileMigration) error {
	verifyTx, err := db.Begin()
	if err != nil {
		return fmt.Errorf("Error opening transaction for verify: %v", err)
	}

	_, verifyErr := verifyTx.Exec(migration.VerifySQL)
	rollbackError := verifyTx.Rollback()
	if verifyErr != nil && rollbackError != nil {
		return fmt.Errorf(
			dedent.Dedent(`
				Got an error for verify. Please asses the necessity of a down migration.
				Migration: %s
				Verify Error: %s
				Rollback error: %s
			`),
			migration.Filename,
			verifyErr,
			rollbackError,
		)
	} else if verifyErr != nil {
		return fmt.Errorf("Error during verify for %s: %s", migration.Filename, verifyErr)
	} else if rollbackError != nil {
		return fmt.Errorf(
			"Error during rollback of verify for %s: %s",
			migration.Filename,
			rollbackError,
		)
	}
	return nil
}
