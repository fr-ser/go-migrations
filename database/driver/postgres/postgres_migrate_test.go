package postgres

import (
	"database/sql"
	"fmt"
	"testing"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/kylelemons/godebug/pretty"

	"go-migrations/database"
)

type migrateCallArgs struct {
	migration      database.FileMigration
	changelogTable string
}

func TestApplyAllUpMigrations(t *testing.T) {
	defer resetMockVariables()
	db, mock, err := sqlmock.New(sqlmock.QueryMatcherOption(sqlmock.QueryMatcherEqual))
	mockableSQLOpen = func(a, b string) (*sql.DB, error) { return db, err }
	mock.ExpectClose()

	mockableGetFileMigrations = func(a string) ([]database.FileMigration, error) {
		return []database.FileMigration{{ID: "1"}, {ID: "2"}}, nil
	}

	migrateUpCalls := []migrateCallArgs{}
	mockableApplyUpMigration = func(a *sql.DB, b database.FileMigration, c string) error {
		migrateUpCalls = append(migrateUpCalls, migrateCallArgs{migration: b, changelogTable: c})
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

	if diff := pretty.Compare(expectedArgs, migrateUpCalls); diff != "" {
		t.Errorf("Did not pass right FileMigrations to migrateUp:\n%s", diff)
	}

	if err := mock.ExpectationsWereMet(); err != nil {
		t.Errorf("there were unfulfilled expectations: %s", err)
	}
}

func TestApplyDownMigrationsWithCount(t *testing.T) {
	defer resetMockVariables()
	db, mock, err := sqlmock.New(sqlmock.QueryMatcherOption(sqlmock.QueryMatcherEqual))
	mockableSQLOpen = func(a, b string) (*sql.DB, error) { return db, err }
	mock.ExpectClose()

	receivedMigrateArgs := []migrateCallArgs{}
	mockableApplyDownMigration = func(a *sql.DB, b database.FileMigration, c string) error {
		receivedMigrateArgs = append(
			receivedMigrateArgs,
			migrateCallArgs{migration: b, changelogTable: c},
		)
		return nil
	}

	type filterByCountArgs struct {
		count             uint
		all               bool
		fileMigrations    []database.FileMigration
		appliedMigrations []database.AppliedMigration
	}

	var receivedFilterByCountArgs filterByCountArgs
	mockableFilterDownMigrationsByCount = func(a uint, b bool, c []database.FileMigration,
		d []database.AppliedMigration) ([]database.FileMigration, error) {
		receivedFilterByCountArgs = filterByCountArgs{
			count: a, all: b, fileMigrations: c, appliedMigrations: d,
		}
		return []database.FileMigration{{ID: "2"}, {ID: "3"}}, nil
	}

	fileMigrations := []database.FileMigration{{ID: "1"}, {ID: "2"}, {ID: "3"}}
	appliedMigrations := []database.AppliedMigration{{ID: "1"}}

	expectedMigrateArgs := []migrateCallArgs{
		{migration: database.FileMigration{ID: "2"}, changelogTable: changelogTable},
		{migration: database.FileMigration{ID: "3"}, changelogTable: changelogTable},
	}
	expectedFilterByCountArgs := filterByCountArgs{
		count: 2, all: false, fileMigrations: fileMigrations,
		appliedMigrations: appliedMigrations,
	}

	pg := Postgres{}
	pg.fileMigrations = fileMigrations
	pg.appliedMigrations = appliedMigrations
	err = pg.ApplyDownMigrationsWithCount(
		expectedFilterByCountArgs.count, expectedFilterByCountArgs.all,
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

func TestApplyUpMigrationsWithCount(t *testing.T) {
	defer resetMockVariables()
	db, mock, err := sqlmock.New(sqlmock.QueryMatcherOption(sqlmock.QueryMatcherEqual))
	mockableSQLOpen = func(a, b string) (*sql.DB, error) { return db, err }
	mock.ExpectClose()

	receivedMigrateArgs := []migrateCallArgs{}
	mockableApplyUpMigration = func(a *sql.DB, b database.FileMigration, c string) error {
		receivedMigrateArgs = append(
			receivedMigrateArgs,
			migrateCallArgs{migration: b, changelogTable: c},
		)
		return nil
	}

	type filterByCountArgs struct {
		count             uint
		all               bool
		fileMigrations    []database.FileMigration
		appliedMigrations []database.AppliedMigration
	}

	var receivedFilterByCountArgs filterByCountArgs
	mockableFilterUpMigrationsByCount = func(a uint, b bool, c []database.FileMigration,
		d []database.AppliedMigration) ([]database.FileMigration, error) {
		receivedFilterByCountArgs = filterByCountArgs{
			count: a, all: b, fileMigrations: c, appliedMigrations: d,
		}
		return []database.FileMigration{{ID: "2"}, {ID: "3"}}, nil
	}

	fileMigrations := []database.FileMigration{{ID: "1"}, {ID: "2"}, {ID: "3"}}
	appliedMigrations := []database.AppliedMigration{{ID: "1"}}

	expectedMigrateArgs := []migrateCallArgs{
		{migration: database.FileMigration{ID: "2"}, changelogTable: changelogTable},
		{migration: database.FileMigration{ID: "3"}, changelogTable: changelogTable},
	}
	expectedFilterByCountArgs := filterByCountArgs{
		count: 2, all: false, fileMigrations: fileMigrations,
		appliedMigrations: appliedMigrations,
	}

	pg := Postgres{}
	pg.fileMigrations = fileMigrations
	pg.appliedMigrations = appliedMigrations
	err = pg.ApplyUpMigrationsWithCount(
		expectedFilterByCountArgs.count, expectedFilterByCountArgs.all,
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

func TestApplyDownMigrationsWithCountError(t *testing.T) {
	defer resetMockVariables()
	db, mock, err := sqlmock.New(sqlmock.QueryMatcherOption(sqlmock.QueryMatcherEqual))
	mockableSQLOpen = func(a, b string) (*sql.DB, error) { return db, err }
	mock.ExpectClose()

	var migrateCalled bool
	mockableApplyUpMigration = func(a *sql.DB, b database.FileMigration, c string) error {
		migrateCalled = true
		return nil
	}
	mockableFilterDownMigrationsByCount = func(a uint, b bool, c []database.FileMigration,
		d []database.AppliedMigration) ([]database.FileMigration, error) {
		return []database.FileMigration{}, fmt.Errorf("test")
	}

	pg := Postgres{}
	pg.fileMigrations = []database.FileMigration{{ID: "1"}}
	pg.appliedMigrations = []database.AppliedMigration{{ID: "1"}}
	err = pg.ApplyDownMigrationsWithCount(3, false)
	if err == nil {
		t.Errorf("Expected error, but got nothing")
	}

	if migrateCalled {
		t.Errorf("Did not expect to call down migrations, but they were called")
	}

	if err := mock.ExpectationsWereMet(); err != nil {
		t.Errorf("there were unfulfilled expectations: %s", err)
	}
}

func TestApplyUpMigrationsWithCountError(t *testing.T) {
	defer resetMockVariables()
	db, mock, err := sqlmock.New(sqlmock.QueryMatcherOption(sqlmock.QueryMatcherEqual))
	mockableSQLOpen = func(a, b string) (*sql.DB, error) { return db, err }
	mock.ExpectClose()

	var migrateCalled bool
	mockableApplyUpMigration = func(a *sql.DB, b database.FileMigration, c string) error {
		migrateCalled = true
		return nil
	}
	mockableFilterUpMigrationsByCount = func(a uint, b bool, c []database.FileMigration,
		d []database.AppliedMigration) ([]database.FileMigration, error) {
		return []database.FileMigration{}, fmt.Errorf("test")
	}

	pg := Postgres{}
	pg.fileMigrations = []database.FileMigration{{ID: "1"}}
	pg.appliedMigrations = []database.AppliedMigration{{ID: "1"}}
	err = pg.ApplyUpMigrationsWithCount(3, false)
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

func TestApplySpecificUpMigration(t *testing.T) {
	defer resetMockVariables()
	db, mock, err := sqlmock.New(sqlmock.QueryMatcherOption(sqlmock.QueryMatcherEqual))
	mockableSQLOpen = func(a, b string) (*sql.DB, error) { return db, err }
	mock.ExpectClose()

	var migrateMigration database.FileMigration
	var migrateChangelog string
	mockableApplyUpMigration = func(a *sql.DB, b database.FileMigration, c string) error {
		migrateMigration = b
		migrateChangelog = c
		return nil
	}

	expectedMigration := database.FileMigration{ID: "expected"}
	var FilterMigrationsByTextFilter string
	mockableFilterUpMigrationsByText = func(a string, b []database.FileMigration,
		c []database.AppliedMigration) (database.FileMigration, error) {
		FilterMigrationsByTextFilter = a
		return expectedMigration, nil
	}

	pg := Postgres{}
	pg.fileMigrations = []database.FileMigration{}
	pg.appliedMigrations = []database.AppliedMigration{}
	err = pg.ApplySpecificUpMigration("sth")
	if err != nil {
		t.Errorf("Expected no error, but got %v", err)
	}

	if migrateChangelog != changelogTable {
		t.Errorf("Expected changelogtable '%s', but got %s", changelogTable, migrateChangelog)
	}
	if migrateMigration != expectedMigration {
		t.Errorf("Expected migration '%v', but got %v", expectedMigration, migrateMigration)
	}
	if FilterMigrationsByTextFilter != "sth" {
		t.Errorf(
			"Expected FilterUpMigration to be called with 'sht', but got %s",
			FilterMigrationsByTextFilter,
		)
	}

	if err := mock.ExpectationsWereMet(); err != nil {
		t.Errorf("there were unfulfilled expectations: %s", err)
	}
}

func TestApplySpecificDownMigration(t *testing.T) {
	defer resetMockVariables()
	db, mock, err := sqlmock.New(sqlmock.QueryMatcherOption(sqlmock.QueryMatcherEqual))
	mockableSQLOpen = func(a, b string) (*sql.DB, error) { return db, err }
	mock.ExpectClose()

	var migrateMigration database.FileMigration
	var migrateChangelog string
	mockableApplyDownMigration = func(a *sql.DB, b database.FileMigration, c string) error {
		migrateMigration = b
		migrateChangelog = c
		return nil
	}

	expectedMigration := database.FileMigration{ID: "expected"}
	var FilterMigrationsByTextFilter string
	mockableFilterDownMigrationsByText = func(a string, b []database.FileMigration,
		c []database.AppliedMigration) (database.FileMigration, error) {
		FilterMigrationsByTextFilter = a
		return expectedMigration, nil
	}

	pg := Postgres{}
	pg.fileMigrations = []database.FileMigration{}
	pg.appliedMigrations = []database.AppliedMigration{}
	err = pg.ApplySpecificDownMigration("sth")
	if err != nil {
		t.Errorf("Expected no error, but got %v", err)
	}

	if migrateChangelog != changelogTable {
		t.Errorf("Expected changelogtable '%s', but got %s", changelogTable, migrateChangelog)
	}
	if migrateMigration != expectedMigration {
		t.Errorf("Expected migration '%v', but got %v", expectedMigration, migrateMigration)
	}
	if FilterMigrationsByTextFilter != "sth" {
		t.Errorf(
			"Expected FilterDownMigration to be called with 'sht', but got %s",
			FilterMigrationsByTextFilter,
		)
	}

	if err := mock.ExpectationsWereMet(); err != nil {
		t.Errorf("there were unfulfilled expectations: %s", err)
	}
}

func TestApplySpecificUpMigrationError(t *testing.T) {
	defer resetMockVariables()
	db, mock, err := sqlmock.New(sqlmock.QueryMatcherOption(sqlmock.QueryMatcherEqual))
	mockableSQLOpen = func(a, b string) (*sql.DB, error) { return db, err }
	mock.ExpectClose()

	var migrateCalled bool
	mockableApplyUpMigration = func(a *sql.DB, b database.FileMigration, c string) error {
		migrateCalled = true
		return nil
	}

	mockableFilterUpMigrationsByText = func(a string, b []database.FileMigration,
		c []database.AppliedMigration) (database.FileMigration, error) {
		return database.FileMigration{}, fmt.Errorf("test")
	}

	pg := Postgres{}
	pg.fileMigrations = []database.FileMigration{}
	pg.appliedMigrations = []database.AppliedMigration{}
	err = pg.ApplySpecificUpMigration("sth")
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

func TestApplySpecificDownMigrationError(t *testing.T) {
	defer resetMockVariables()
	db, mock, err := sqlmock.New(sqlmock.QueryMatcherOption(sqlmock.QueryMatcherEqual))
	mockableSQLOpen = func(a, b string) (*sql.DB, error) { return db, err }
	mock.ExpectClose()

	var migrateCalled bool
	mockableApplyDownMigration = func(a *sql.DB, b database.FileMigration, c string) error {
		migrateCalled = true
		return nil
	}

	mockableFilterDownMigrationsByText = func(a string, b []database.FileMigration,
		c []database.AppliedMigration) (database.FileMigration, error) {
		return database.FileMigration{}, fmt.Errorf("test")
	}

	pg := Postgres{}
	pg.fileMigrations = []database.FileMigration{}
	pg.appliedMigrations = []database.AppliedMigration{}
	err = pg.ApplySpecificDownMigration("sth")
	if err == nil {
		t.Errorf("Expected error, but got none")
	}

	if migrateCalled {
		t.Errorf("Expected migrateDown not to be called, but it was")
	}

	if err := mock.ExpectationsWereMet(); err != nil {
		t.Errorf("there were unfulfilled expectations: %s", err)
	}
}
