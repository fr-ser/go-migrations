package database

import (
	"database/sql"
	"fmt"
	"io/ioutil"
	"path/filepath"
	"sort"
	"time"
)

// GetFileMigrations gets all migration files within the database/migration folder's subfolders
// it returns a list of FileMigrations sorted (ascending) by the ID.
func GetFileMigrations(migrationFolder string) (migrations []FileMigration, err error) {
	skippedFolders := map[string]bool{"_environments": true}
	fileMigrations := map[string]FileMigration{}

	apps, err := ioutil.ReadDir(migrationFolder)
	if err != nil {
		return nil, fmt.Errorf(
			"Could not read content of migrationFolder %s - Err: %v", migrationFolder, err,
		)
	}
	for _, app := range apps {
		if skippedFolders[app.Name()] || !app.IsDir() {
			continue
		}

		migFiles, err := ioutil.ReadDir(filepath.Join(migrationFolder, app.Name()))
		if err != nil {
			return nil, fmt.Errorf(
				"Could not read content of appFolder %s - Err: %v", app.Name(), err,
			)
		}

		for _, migFile := range migFiles {
			if migFile.IsDir() {
				continue
			}

			mig := FileMigration{}
			err = mig.LoadFromFile(filepath.Join(migrationFolder, app.Name(), migFile.Name()))
			if err != nil {
				return nil, err
			}

			if prevMig, alreadyExists := fileMigrations[mig.ID]; alreadyExists {
				return nil, fmt.Errorf(
					"The id %s is not unique. It exists for %s and %s",
					mig.ID, prevMig.Filename, mig.Filename,
				)
			}
			fileMigrations[mig.ID] = mig
		}

	}

	for _, fileMigration := range fileMigrations {
		migrations = append(migrations, fileMigration)
	}
	sort.Slice(migrations, func(i, j int) bool { return migrations[i].ID < migrations[j].ID })
	return migrations, nil
}

// GetAppliedMigrations gets all applied migrations from the changelog (sorted by ID)
func GetAppliedMigrations(db *sql.DB, changelogTable string) (
	migrations []AppliedMigration, err error,
) {
	rows, err := db.Query(fmt.Sprintf(
		`SELECT id, name, applied_at FROM %s ORDER BY id ASC`,
		changelogTable,
	))
	if err != nil {
		return nil, fmt.Errorf("Got error getting applied migrations: %v", err)
	}
	defer rows.Close()
	for rows.Next() {
		var id, name string
		var appliedAt time.Time
		if err := rows.Scan(&id, &name, &appliedAt); err != nil {
			return nil, fmt.Errorf("Error scanning row for applied migrations: %v", err)
		}
		migrations = append(migrations, AppliedMigration{ID: id, Name: name, AppliedAt: appliedAt})
	}
	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("Error after row iteration for getting applied migrations: %v", err)
	}

	return migrations, nil
}
