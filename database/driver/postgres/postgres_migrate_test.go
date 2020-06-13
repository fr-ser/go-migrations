package postgres

import (
	"database/sql"
	"fmt"
	"testing"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/kylelemons/godebug/pretty"

	"go-migrations/database"
	"go-migrations/internal/direction"
)

type migrateCallArgs struct {
	migration      database.FileMigration
	changelogTable string
	direction      direction.MigrateDirection
}

func TestApplyAllUpMigrations(t *testing.T) {
	defer resetMockVariables()
	db, mock, err := sqlmock.New(sqlmock.QueryMatcherOption(sqlmock.QueryMatcherEqual))
	mockableSQLOpen = func(a, b string) (*sql.DB, error) { return db, err }
	mock.ExpectClose()

	mockableGetFileMigrations = func(a string) ([]database.FileMigration, error) {
		return []database.FileMigration{{ID: "1"}, {ID: "2"}}, nil
	}

	receivedMigrateArgs := []migrateCallArgs{}
	mockableApplyMigration = func(
		db *sql.DB, f database.FileMigration, c string, d direction.MigrateDirection,
	) error {
		receivedMigrateArgs = append(
			receivedMigrateArgs,
			migrateCallArgs{migration: f, changelogTable: c, direction: d},
		)
		return nil
	}

	expectedArgs := []migrateCallArgs{
		{migration: database.FileMigration{ID: "1"}, changelogTable: changelogTable},
		{migration: database.FileMigration{ID: "2"}, changelogTable: changelogTable},
	}

	pg := Postgres{}
	err = pg.ApplyAllUpMigrations()
	if err != nil {
		t.Errorf("Expected no error, but got: %s", err)
	}

	if diff := pretty.Compare(expectedArgs, receivedMigrateArgs); diff != "" {
		t.Errorf("Did not pass right FileMigrations to migrateUp:\n%s", diff)
	}

	if err := mock.ExpectationsWereMet(); err != nil {
		t.Errorf("there were unfulfilled expectations: %s", err)
	}
}

func TestApplyMigrationsWithCount(t *testing.T) {
	defer resetMockVariables()
	for _, dir := range direction.Directions {

		db, mock, err := sqlmock.New(sqlmock.QueryMatcherOption(sqlmock.QueryMatcherEqual))
		mockableSQLOpen = func(a, b string) (*sql.DB, error) { return db, err }
		mock.ExpectClose()

		receivedMigrateArgs := []migrateCallArgs{}
		mockableApplyMigration = func(
			db *sql.DB, f database.FileMigration, c string, d direction.MigrateDirection,
		) error {
			receivedMigrateArgs = append(
				receivedMigrateArgs,
				migrateCallArgs{migration: f, changelogTable: c, direction: d},
			)
			return nil
		}

		type filterByCountArgs struct {
			count             uint
			all               bool
			direction         direction.MigrateDirection
			fileMigrations    []database.FileMigration
			appliedMigrations []database.AppliedMigration
		}

		var receivedFilterByCountArgs filterByCountArgs
		mockableFilterMigrationsByCount = func(c uint, a bool, d direction.MigrateDirection,
			f []database.FileMigration, app []database.AppliedMigration) (
			[]database.FileMigration, error,
		) {
			receivedFilterByCountArgs = filterByCountArgs{
				count: c, all: a, fileMigrations: f, appliedMigrations: app, direction: d,
			}
			return []database.FileMigration{{ID: "2"}, {ID: "3"}}, nil
		}

		fileMigrations := []database.FileMigration{{ID: "1"}, {ID: "2"}, {ID: "3"}}
		appliedMigrations := []database.AppliedMigration{{ID: "1"}}

		expectedMigrateArgs := []migrateCallArgs{
			{
				migration:      database.FileMigration{ID: "2"},
				changelogTable: changelogTable,
				direction:      dir.Direction,
			},
			{
				migration:      database.FileMigration{ID: "3"},
				changelogTable: changelogTable,
				direction:      dir.Direction,
			},
		}
		expectedFilterByCountArgs := filterByCountArgs{
			count: 2, all: false, direction: dir.Direction,
			fileMigrations: fileMigrations, appliedMigrations: appliedMigrations,
		}

		pg := Postgres{}
		pg.fileMigrations = fileMigrations
		pg.appliedMigrations = appliedMigrations
		err = pg.ApplyMigrationsWithCount(
			expectedFilterByCountArgs.count, expectedFilterByCountArgs.all, dir.Direction,
		)
		if err != nil {
			t.Errorf("Expected no error, but got: %s", err)
		}

		if diff := pretty.Compare(expectedFilterByCountArgs, receivedFilterByCountArgs); diff != "" {
			t.Errorf("Did not pass right arguments to FilterByCount:\n%s", diff)
		}

		if diff := pretty.Compare(expectedMigrateArgs, receivedMigrateArgs); diff != "" {
			t.Errorf("Did not pass right FileMigrations to migrateUp:\n%s", diff)
		}

		if err := mock.ExpectationsWereMet(); err != nil {
			t.Errorf("there were unfulfilled expectations: %s", err)
		}
	}
}

func TestApplyMigrationsWithCountError(t *testing.T) {
	defer resetMockVariables()
	db, mock, err := sqlmock.New(sqlmock.QueryMatcherOption(sqlmock.QueryMatcherEqual))
	mockableSQLOpen = func(a, b string) (*sql.DB, error) { return db, err }
	mock.ExpectClose()

	var migrateCalled bool
	mockableApplyMigration = func(
		db *sql.DB, f database.FileMigration, c string, d direction.MigrateDirection,
	) error {
		migrateCalled = true
		return nil
	}
	mockableFilterMigrationsByCount = func(c uint, a bool, d direction.MigrateDirection,
		f []database.FileMigration, app []database.AppliedMigration) (
		[]database.FileMigration, error,
	) {
		return []database.FileMigration{}, fmt.Errorf("test")
	}

	pg := Postgres{}
	pg.fileMigrations = []database.FileMigration{{ID: "1"}}
	pg.appliedMigrations = []database.AppliedMigration{{ID: "1"}}
	err = pg.ApplyMigrationsWithCount(3, false, direction.Up)
	if err == nil {
		t.Errorf("Expected error, but got nothing")
	}

	if migrateCalled {
		t.Errorf("Did not expect to call up migrations, but they were called")
	}

	if err := mock.ExpectationsWereMet(); err != nil {
		t.Errorf("there were unfulfilled expectations: %s", err)
	}
}

func TestApplySpecificMigration(t *testing.T) {
	defer resetMockVariables()
	for _, dir := range direction.Directions {
		db, mock, err := sqlmock.New(sqlmock.QueryMatcherOption(sqlmock.QueryMatcherEqual))
		mockableSQLOpen = func(a, b string) (*sql.DB, error) { return db, err }
		mock.ExpectClose()

		var migrateMigration database.FileMigration
		var migrateChangelog string
		var migrateDirection direction.MigrateDirection
		mockableApplyMigration = func(db *sql.DB, f database.FileMigration, c string,
			d direction.MigrateDirection,
		) error {
			migrateMigration = f
			migrateChangelog = c
			migrateDirection = d
			return nil
		}

		expectedMigration := database.FileMigration{ID: "expected"}
		var filterMigrationsByTextFilter string
		var filterMigrationsByTextDirection direction.MigrateDirection
		mockableFilterMigrationsByText = func(fi string, d direction.MigrateDirection,
			f []database.FileMigration, a []database.AppliedMigration) (database.FileMigration, error) {
			filterMigrationsByTextFilter = fi
			filterMigrationsByTextDirection = d
			return expectedMigration, nil
		}

		pg := Postgres{}
		pg.fileMigrations = []database.FileMigration{}
		pg.appliedMigrations = []database.AppliedMigration{}
		err = pg.ApplySpecificMigration("sth", dir.Direction)
		if err != nil {
			t.Errorf("Expected no error, but got %v", err)
		}

		if migrateChangelog != changelogTable {
			t.Errorf("Expected changelogtable '%s', but got %s", changelogTable, migrateChangelog)
		}
		if migrateDirection != dir.Direction {
			t.Errorf("Expected %s migration, but got the other direction", dir.Name)
		}
		if migrateMigration != expectedMigration {
			t.Errorf("Expected migration '%v', but got %v", expectedMigration, migrateMigration)
		}
		if filterMigrationsByTextFilter != "sth" {
			t.Errorf(
				"Expected FilterUpMigration to be called with 'sht', but got %s",
				filterMigrationsByTextFilter,
			)
		}
		if filterMigrationsByTextDirection != dir.Direction {
			t.Errorf(
				"Expected FilterMigration for %s to be called, but got the other direction",
				dir.Name,
			)
		}

		if err := mock.ExpectationsWereMet(); err != nil {
			t.Errorf("there were unfulfilled expectations: %s", err)
		}
	}
}

func TestApplySpecificMigrationError(t *testing.T) {
	defer resetMockVariables()
	db, mock, err := sqlmock.New(sqlmock.QueryMatcherOption(sqlmock.QueryMatcherEqual))
	mockableSQLOpen = func(a, b string) (*sql.DB, error) { return db, err }
	mock.ExpectClose()

	var migrateCalled bool
	mockableApplyMigration = func(
		db *sql.DB, f database.FileMigration, c string, d direction.MigrateDirection,
	) error {
		migrateCalled = true
		return nil
	}

	mockableFilterMigrationsByText = func(fi string, d direction.MigrateDirection,
		f []database.FileMigration, a []database.AppliedMigration) (database.FileMigration, error) {
		return database.FileMigration{}, fmt.Errorf("test")
	}

	pg := Postgres{}
	pg.fileMigrations = []database.FileMigration{}
	pg.appliedMigrations = []database.AppliedMigration{}
	err = pg.ApplySpecificMigration("sth", direction.Up)
	if err == nil {
		t.Errorf("Expected error, but got none")
	}

	if migrateCalled {
		t.Errorf("Expected migrateUp not to be called, but it was")
	}

	if err := mock.ExpectationsWereMet(); err != nil {
		t.Errorf("there were unfulfilled expectations: %s", err)
	}
}
