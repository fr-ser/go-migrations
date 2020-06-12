package database

import (
	"database/sql"
	"fmt"
	"strings"

	log "github.com/sirupsen/logrus"

	"github.com/lithammer/dedent"
)

var (
	mockableMigrateUp         = ApplyUpSQL
	mockableInsertToChangelog = InsertToChangelog
	mockableApplyVerify       = ApplyVerify
)

// FilterUpMigrationsByText filters the migrations by filename. If more then one unapplied migration
// remain an error is thrown
func FilterUpMigrationsByText(filter string, fileMigrations []FileMigration,
	appliedMigrations []AppliedMigration) (filteredMigration FileMigration, err error,
) {
	appliedIDLookup := map[string]bool{}
	for _, mig := range appliedMigrations {
		appliedIDLookup[mig.ID] = true
	}

	foundMigrations := []FileMigration{}
	for _, mig := range fileMigrations {
		if strings.Contains(mig.Filename, filter) && !appliedIDLookup[mig.ID] {
			foundMigrations = append(foundMigrations, mig)
		}
	}

	if len(foundMigrations) == 0 {
		return filteredMigration, fmt.Errorf("Found no migration matching the filter: %s", filter)
	} else if len(foundMigrations) > 1 {
		matchedNames := ""
		for _, mig := range foundMigrations {
			matchedNames = fmt.Sprintf("%s\n%s/%s", matchedNames, mig.Application, mig.Filename)
		}
		return filteredMigration, fmt.Errorf(
			"Found multiple matches for the filter: %s %s", filter, matchedNames,
		)
	}

	filteredMigration = foundMigrations[0]

	return filteredMigration, nil
}

// FilterUpMigrationsByCount filters the migrations for the next n unapplied migrations.
// This relies on a "consistent changelog" and return an error if no migrations are left to apply
func FilterUpMigrationsByCount(count uint, all bool, fileMigrations []FileMigration,
	appliedMigrations []AppliedMigration) (migrations []FileMigration, err error,
) {
	appliedCount := len(appliedMigrations)

	if count > 0 && appliedCount+int(count) > len(fileMigrations) {
		log.Warningf(
			dedent.Dedent(`
				The received count (%d) is bigger than the remaining migrations.
				All migrations will be applied.
			`), count,
		)
		all = true
	}
	if appliedCount == len(fileMigrations) {
		return migrations, fmt.Errorf("No migrations left to apply")
	}

	if all {
		migrations = fileMigrations[appliedCount:]
	} else {
		migrations = fileMigrations[appliedCount : appliedCount+int(count)]
	}

	return migrations, nil
}

// ApplyUpMigration applies the up migration in a transaction
// After the migration a verify script is executed and rolled back in a separate transaction.
// If the verify script fails the downmigration is executed (also in a transaction)
func ApplyUpMigration(db *sql.DB, migration FileMigration, changelogTable string) error {
	if err := mockableMigrateUp(db, migration); err != nil {
		return err
	}

	if err := mockableInsertToChangelog(db, migration, changelogTable); err != nil {
		return err
	}

	if err := mockableApplyVerify(db, migration); err != nil {
		return err
	}

	return nil
}

// ApplyUpSQL is an internal helper to apply the up migration in a transaction
// it does not perform anything else (like verify execution)
func ApplyUpSQL(db *sql.DB, migration FileMigration) error {
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
		`INSERT INTO %s (id, name, applied_at) VALUES ('%s', '%s', now())`,
		changelogTable,
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
