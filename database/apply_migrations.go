package database

import (
	"database/sql"
	"fmt"
	"strings"

	log "github.com/sirupsen/logrus"

	"github.com/lithammer/dedent"

	"go-migrations/internal/direction"
)

var (
	mockableMigrateUp           = ApplyUpSQL
	mockableMigrateDown         = ApplyDownSQL
	mockableInsertToChangelog   = InsertToChangelog
	mockableRemoveFromChangelog = RemoveFromChangelog
	mockableApplyVerify         = ApplyVerify
)

// FilterMigrationsByText filters the migrations by filename.
// If more then one migration remains an error is thrown
func FilterMigrationsByText(
	filter string, dir direction.MigrateDirection,
	fileMigrations []FileMigration, appliedMigrations []AppliedMigration,
) (FileMigration, error) {
	if dir == direction.Down {
		return filterDownMigrationsByText(filter, fileMigrations, appliedMigrations)
	}
	return filterUpMigrationsByText(filter, fileMigrations, appliedMigrations)
}

func filterUpMigrationsByText(filter string, fileMigrations []FileMigration,
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

func filterDownMigrationsByText(filter string, fileMigrations []FileMigration,
	appliedMigrations []AppliedMigration) (filteredMigration FileMigration, err error,
) {
	fileIDLookup := map[string]int{}
	for idx, mig := range fileMigrations {
		fileIDLookup[mig.ID] = idx
	}

	foundMigrations := []FileMigration{}
	for _, mig := range appliedMigrations {
		lookupText := fmt.Sprintf("%s_%s.sql", mig.ID, mig.Name)
		if strings.Contains(lookupText, filter) {
			idx, exists := fileIDLookup[mig.ID]
			if exists {
				foundMigrations = append(foundMigrations, fileMigrations[idx])
			}
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

// FilterMigrationsByCount filters the migrations for the next n migrations.
// This relies on a "consistent changelog" and returns an error if no migrations are left
// For down migrations the migrations are sorted reversed (descending by ID)
func FilterMigrationsByCount(
	count uint, all bool, dir direction.MigrateDirection,
	fileMigrations []FileMigration, appliedMigrations []AppliedMigration,
) ([]FileMigration, error) {
	if dir == direction.Down {
		return filterDownMigrationsByCount(count, all, fileMigrations, appliedMigrations)
	}
	return filterUpMigrationsByCount(count, all, fileMigrations, appliedMigrations)
}

func filterUpMigrationsByCount(count uint, all bool, fileMigrations []FileMigration,
	appliedMigrations []AppliedMigration) (migrations []FileMigration, err error,
) {
	appliedCount := len(appliedMigrations)

	if appliedCount == len(fileMigrations) {
		return migrations, fmt.Errorf("No migrations left to apply")
	}
	if count > 0 && appliedCount+int(count) > len(fileMigrations) {
		log.Warningf(
			dedent.Dedent(`
				The received count (%d) is bigger than the remaining migrations.
				All migrations will be applied.
			`), count,
		)
		all = true
	}

	if all {
		migrations = fileMigrations[appliedCount:]
	} else {
		migrations = fileMigrations[appliedCount : appliedCount+int(count)]
	}

	return migrations, nil
}

func filterDownMigrationsByCount(count uint, all bool, fileMigrations []FileMigration,
	appliedMigrations []AppliedMigration) (migrations []FileMigration, err error,
) {
	appliedCount := len(appliedMigrations)
	if appliedCount == 0 {
		return migrations, fmt.Errorf("No migrations left to remove")
	}

	if appliedCount < int(count) {
		log.Warningf(
			dedent.Dedent(`
				The received count (%d) is bigger than the applied migrations.
				All migrations will be removed.
			`), count,
		)
		all = true
	}

	var lastIdx int
	if all {
		lastIdx = 0
	} else {
		lastIdx = appliedCount - int(count)
	}

	for idx := appliedCount - 1; idx >= lastIdx; idx-- {
		migrations = append(migrations, fileMigrations[idx])
	}

	return migrations, nil
}

// ApplyMigration applies a migration in a transaction and updates the changelog
// For up migrations a verify script is executed and rolled back in a separate transaction.
func ApplyMigration(
	db *sql.DB, migration FileMigration, changelogTable string, dir direction.MigrateDirection,
) error {
	if dir == direction.Down {
		return applyDownMigration(db, migration, changelogTable)
	}
	return applyUpMigration(db, migration, changelogTable)
}

func applyDownMigration(db *sql.DB, migration FileMigration, changelogTable string) error {
	if err := mockableMigrateDown(db, migration); err != nil {
		return err
	}

	if err := mockableRemoveFromChangelog(db, migration, changelogTable); err != nil {
		return err
	}

	return nil
}

func applyUpMigration(db *sql.DB, migration FileMigration, changelogTable string) error {
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
		return fmt.Errorf(
			"Error during commit of up migration of %s: %s",
			migration.Filename, err,
		)
	}
	return nil
}

// ApplyDownSQL is an internal helper to apply the down migration in a transaction
// it does not perform anything else (like changelog update)
func ApplyDownSQL(db *sql.DB, migration FileMigration) error {
	tx, err := db.Begin()
	if err != nil {
		return fmt.Errorf("Error opening transaction: %v", err)
	}

	_, err = tx.Exec(migration.DownSQL)
	if err != nil {
		rollbackError := tx.Rollback()
		if rollbackError != nil {
			return fmt.Errorf(
				"Error during down migration of %s: %s \n and rollback error: %s",
				migration.Filename,
				err,
				rollbackError,
			)
		}
		return fmt.Errorf("Error during down migration of %s: %s", migration.Filename, err)
	}

	err = tx.Commit()
	if err != nil {
		return fmt.Errorf(
			"Error during commit of down migration of %s: %s",
			migration.Filename, err,
		)
	}
	return nil
}

// InsertToChangelog is an internal helper to insert the migration into the changelog
func InsertToChangelog(db *sql.DB, migration FileMigration, changelogTable string) error {
	_, err := db.Exec(fmt.Sprintf(
		`INSERT INTO %s (id, name, applied_at) VALUES ('%s', '%s', now())`,
		changelogTable,
		migration.ID,
		migration.Description,
	))
	if err != nil {
		return fmt.Errorf(
			"Could not add the migration %s from the changelog: %v",
			migration.Filename, err,
		)
	}
	return nil
}

// RemoveFromChangelog is an internal helper to remove the migration from the changelog
func RemoveFromChangelog(db *sql.DB, migration FileMigration, changelogTable string) error {
	_, err := db.Exec(fmt.Sprintf(`DELETE FROM %s WHERE id = '%s'`, changelogTable, migration.ID))
	if err != nil {
		return fmt.Errorf(
			"Could not remove the migration %s from the changelog: %v",
			migration.Filename, err,
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
