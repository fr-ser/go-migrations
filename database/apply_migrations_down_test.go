package database

import (
	"database/sql"
	"fmt"
	"go-migrations/internal/direction"
	"testing"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/kylelemons/godebug/pretty"
)

func TestApplyDownMigration(t *testing.T) {
	expectedMigration := FileMigration{DownSQL: "SELECT 1", ID: "1"}

	var migrateDownCall FileMigration
	mockableMigrateDown = func(a *sql.DB, b FileMigration) error {
		migrateDownCall = b
		return nil
	}

	var removeFromChangelogCall FileMigration
	var calledChangelogTable string
	mockableRemoveFromChangelog = func(a *sql.DB, b FileMigration, c string) error {
		removeFromChangelogCall = b
		calledChangelogTable = c
		return nil
	}

	db, _, _ := sqlmock.New(sqlmock.QueryMatcherOption(sqlmock.QueryMatcherEqual))
	err := ApplyMigration(db, expectedMigration, "schema.sth", direction.Down)
	if err != nil {
		t.Errorf("Expected no error for applying up migrations, but got: %s", err)
	}

	if diff := pretty.Compare(expectedMigration, migrateDownCall); diff != "" {
		t.Errorf("Did not pass right FileMigrations to migrateDown:\n%s", diff)
	}

	if diff := pretty.Compare(expectedMigration, removeFromChangelogCall); diff != "" {
		t.Errorf("Did not pass right FileMigrations to insertToChangelog:\n%s", diff)
	}
	if calledChangelogTable != "schema.sth" {
		t.Errorf(
			"Expected passed changelogTable to be %s, but got %s",
			"schema.sth", calledChangelogTable,
		)
	}
}

func TestApplyDownMigrationDownMigrationError(t *testing.T) {
	var migrateDownCalled bool
	mockableMigrateDown = func(a *sql.DB, b FileMigration) error {
		migrateDownCalled = true
		return fmt.Errorf("test error")
	}

	var removeFromChangelogCalled bool
	mockableRemoveFromChangelog = func(a *sql.DB, b FileMigration, c string) error {
		removeFromChangelogCalled = true
		return nil
	}

	db, _, _ := sqlmock.New(sqlmock.QueryMatcherOption(sqlmock.QueryMatcherEqual))
	err := ApplyMigration(db, FileMigration{}, "sth", direction.Down)
	if err == nil {
		t.Errorf("Expected error for applying up migrations, but got nothing")
	}

	if !migrateDownCalled {
		t.Errorf("Expected migrateDown to be called once")
	}
	if removeFromChangelogCalled {
		t.Errorf("Expected insertToChangelog not to be called")
	}
}

func TestApplyDownMigrationChangelogError(t *testing.T) {
	var migrateDownCalled bool
	mockableMigrateDown = func(a *sql.DB, b FileMigration) error {
		migrateDownCalled = true
		return nil
	}

	var removeFromChangelogCalled bool
	mockableRemoveFromChangelog = func(a *sql.DB, b FileMigration, c string) error {
		removeFromChangelogCalled = true
		return fmt.Errorf("test error")
	}

	db, _, _ := sqlmock.New(sqlmock.QueryMatcherOption(sqlmock.QueryMatcherEqual))
	err := ApplyMigration(db, FileMigration{}, "sth", direction.Down)
	if err == nil {
		t.Errorf("Expected error for applying up migrations, but got nothing")
	}

	if !migrateDownCalled {
		t.Errorf("Expected migrateDown to be called once")
	}
	if !removeFromChangelogCalled {
		t.Errorf("Expected insertToChangelog  to be called once")
	}
}

func TestApplyDownSQL(t *testing.T) {
	db, mock, err := sqlmock.New(sqlmock.QueryMatcherOption(sqlmock.QueryMatcherEqual))
	migration := FileMigration{DownSQL: "SELECT 1"}

	mock.ExpectBegin()
	mock.ExpectExec("SELECT 1").WillReturnResult(sqlmock.NewResult(1, 1))
	mock.ExpectCommit()

	err = ApplyDownSQL(db, migration)
	if err != nil {
		t.Errorf("Expected no error, but got: %s", err)
	}

	if err := mock.ExpectationsWereMet(); err != nil {
		t.Errorf("there were unfulfilled expectations: %s", err)
	}
}

func TestApplyDownSQLMigrationError(t *testing.T) {
	db, mock, err := sqlmock.New(sqlmock.QueryMatcherOption(sqlmock.QueryMatcherEqual))
	migration := FileMigration{DownSQL: "SELECT 1"}

	mock.ExpectBegin()
	mock.ExpectExec("SELECT 1").WillReturnError(fmt.Errorf("Some error"))
	mock.ExpectRollback()

	err = ApplyDownSQL(db, migration)
	if err == nil {
		t.Errorf("Expected error, but got nothing")
	}

	if err := mock.ExpectationsWereMet(); err != nil {
		t.Errorf("there were unfulfilled expectations: %s", err)
	}
}

func TestRemoveFromChangelog(t *testing.T) {
	db, mock, err := sqlmock.New(sqlmock.QueryMatcherOption(sqlmock.QueryMatcherEqual))
	migration := FileMigration{ID: "1"}

	mock.ExpectExec(`DELETE FROM sth WHERE id = '1'`).WillReturnResult(sqlmock.NewResult(1, 1))

	err = RemoveFromChangelog(db, migration, "sth")
	if err != nil {
		t.Errorf("Expected no error, but got: %s", err)
	}

	if err := mock.ExpectationsWereMet(); err != nil {
		t.Errorf("there were unfulfilled expectations: %s", err)
	}
}

func TestFilterDownMigrationsByCount(t *testing.T) {
	testCases := []struct {
		count       uint
		all         bool
		fileIDs     []string
		appliedIDs  []string
		expectedIDs []string
	}{
		{
			count: 1, all: false, fileIDs: []string{"1"}, appliedIDs: []string{"1"},
			expectedIDs: []string{"1"},
		},
		{
			count: 1, all: false, fileIDs: []string{"1", "2", "3"}, appliedIDs: []string{"1", "2"},
			expectedIDs: []string{"2"},
		},
		{
			count: 2, all: false, fileIDs: []string{"1", "2", "3"}, appliedIDs: []string{"1", "2"},
			expectedIDs: []string{"2", "1"},
		},
		{
			count: 99, all: false,
			fileIDs: []string{"1", "2", "3", "4"}, appliedIDs: []string{"1", "2", "3"},
			expectedIDs: []string{"3", "2", "1"},
		},
		{
			count: 0, all: true, fileIDs: []string{"1", "2"}, appliedIDs: []string{"1", "2"},
			expectedIDs: []string{"2", "1"},
		},
	}

	for _, testCase := range testCases {
		fileMigrations := []FileMigration{}
		for _, id := range testCase.fileIDs {
			fileMigrations = append(fileMigrations, FileMigration{ID: id})
		}
		appliedMigrations := []AppliedMigration{}
		for _, id := range testCase.appliedIDs {
			appliedMigrations = append(appliedMigrations, AppliedMigration{ID: id})
		}
		expectedMigrations := []FileMigration{}
		for _, id := range testCase.expectedIDs {
			expectedMigrations = append(expectedMigrations, FileMigration{ID: id})
		}

		migrations, err := FilterMigrationsByCount(
			testCase.count, testCase.all, direction.Down, fileMigrations, appliedMigrations,
		)
		if err != nil {
			t.Errorf("Expected no error, but got: %s", err)
		}

		if diff := pretty.Compare(expectedMigrations, migrations); diff != "" {
			t.Errorf("Did not pass right FileMigrations to migrateUp:\n%s", diff)
		}
	}
}

func TestFilterDownMigrationsByCountNothingLeft(t *testing.T) {
	_, err := FilterMigrationsByCount(
		3, false, direction.Down, []FileMigration{{ID: "1"}}, []AppliedMigration{},
	)
	if err == nil {
		t.Errorf("Expected error, but got  nothing")
	}
}

func TestFilterDownMigrationsByText(t *testing.T) {
	testCases := []struct {
		filter            string
		expectedFilename  string
		fileMigrations    []FileMigration
		appliedMigrations []AppliedMigration
	}{
		{
			filter: "2017", expectedFilename: "20171101000003_buz.sql",
			fileMigrations: []FileMigration{
				{ID: "20161101000001", Filename: "20161101000001_foo.sql"},
				{ID: "20171101000003", Filename: "20171101000003_buz.sql"},
			},
			appliedMigrations: []AppliedMigration{
				{ID: "20161101000001", Name: "foo"},
				{ID: "20171101000003", Name: "buz"},
			},
		},
		{
			filter: "01_start_a.sql", expectedFilename: "20161101000001_start_a.sql",
			fileMigrations: []FileMigration{
				{ID: "20161101000001", Filename: "20161101000001_start_a.sql"},
				{ID: "20171101000003", Filename: "20171101000003_start_b.sql"},
			},
			appliedMigrations: []AppliedMigration{

				{ID: "20161101000001", Name: "start_a"},
				{ID: "20171101000003", Name: "start_b"},
			},
		},
		{
			filter: "2016", expectedFilename: "20161101000001_a.sql",
			fileMigrations: []FileMigration{
				{ID: "20161101000001", Filename: "20161101000001_a.sql"},
			},
			appliedMigrations: []AppliedMigration{
				{ID: "20161101000001"},
				{ID: "20161101000002"},
			},
		},
	}

	for _, testCase := range testCases {
		migration, err := FilterMigrationsByText(
			testCase.filter, direction.Down,
			testCase.fileMigrations, testCase.appliedMigrations,
		)
		if err != nil {
			t.Errorf("Expected no error, but got: %s", err)
		}

		if migration.Filename != testCase.expectedFilename {
			t.Errorf(
				"Wrong migration: Expected file '%s', but got '%s'",
				testCase.expectedFilename, migration.Filename,
			)
		}
	}
}

func TestFilterDownMigrationsByTextError(t *testing.T) {
	testCases := []struct {
		filter            string
		fileMigrations    []FileMigration
		appliedMigrations []AppliedMigration
	}{

		{
			filter: "2017",
			fileMigrations: []FileMigration{
				{ID: "20161101000001", Filename: "20161101000001_foo.sql"},
			},
			appliedMigrations: []AppliedMigration{{ID: "20161101000001"}},
		},
		{
			filter: "2017",
			fileMigrations: []FileMigration{
				{ID: "20171101000001", Filename: "20171101000001_foo.sql"},
				{ID: "20171101000003", Filename: "20171101000003_buz.sql"},
			},
			appliedMigrations: []AppliedMigration{
				{ID: "20171101000001", Name: "foo"},
				{ID: "20171101000003", Name: "buz"},
			},
		},

		{
			filter:            "2017",
			fileMigrations:    []FileMigration{},
			appliedMigrations: []AppliedMigration{{ID: "20171101000001"}},
		},
	}

	for _, testCase := range testCases {
		_, err := FilterMigrationsByText(
			testCase.filter, direction.Down,
			testCase.fileMigrations, testCase.appliedMigrations,
		)
		if err == nil {
			t.Errorf("Expected error, but got none")
		}
	}
}
