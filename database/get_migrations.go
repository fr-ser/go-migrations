package database

import (
	"fmt"
	"io/ioutil"
	"path/filepath"
	"sort"
)

// GetMigrations gets all migration files within the database/migration folder's subfolders
// it returns a list of FileMigrations sorted (ascending) by the ID.
// The migrations are filtered by application if the passed list is not empty
func GetMigrations(migrationFolder string, applicationFilter []string) (
	migrations []FileMigration, err error,
) {
	skippedFolders := map[string]bool{"_environments": true}
	fileMigrations := map[string]FileMigration{}
	shouldFilter := applicationFilter != nil && len(applicationFilter) > 0
	appsToInclude := map[string]bool{"_common": true}

	if shouldFilter {
		for _, app := range applicationFilter {
			appsToInclude[app] = true
		}
	}

	apps, err := ioutil.ReadDir(migrationFolder)
	if err != nil {
		return nil, fmt.Errorf(
			"Could not read content of migrationFolder %s - Err: %v", migrationFolder, err,
		)
	}
	for _, app := range apps {
		if skippedFolders[app.Name()] || !app.IsDir() {
			continue
		} else if shouldFilter && !appsToInclude[app.Name()] {
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
