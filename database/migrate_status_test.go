package database

import (
	"testing"
	"time"

	"github.com/kylelemons/godebug/pretty"
)

func TestGetMigrationStatus(t *testing.T) {
	fileMigrations := []FileMigration{
		{ID: "1", Application: "buz", Description: "foo_bar"},
		{ID: "2", Application: "fuz", Description: "baz_biz"},
	}
	appliedTime, _ := time.Parse(time.RFC3339, "2020-06-13T17:17:44.371Z")
	appliedMigrations := []AppliedMigration{{ID: "1", Name: "foo_bar", AppliedAt: appliedTime}}

	rows, statusNote, err := GetMigrationStatus(fileMigrations, appliedMigrations)
	if err != nil {
		t.Errorf("Expected no error, but got: %v", err)
	}
	if statusNote != "" {
		t.Errorf("Expected empty statusNote but got: %s", statusNote)
	}
	expectedRows := []MigrateStatusRow{
		{ID: "1", Name: "foo_bar", Application: "buz", Status: "applied at 2020-06-13 17:17:44 UTC"},
		{ID: "2", Name: "baz_biz", Application: "fuz", Status: "not applied"},
	}
	if diff := pretty.Compare(expectedRows, rows); diff != "" {
		t.Errorf(diff)
	}
}

func TestGetMigrationStatusNotFoundLocally(t *testing.T) {
	fileMigrations := []FileMigration{}
	appliedTime, _ := time.Parse(time.RFC3339, "2020-06-13T17:17:44.371Z")
	appliedMigrations := []AppliedMigration{{ID: "1", Name: "foo_bar", AppliedAt: appliedTime}}

	rows, statusNote, err := GetMigrationStatus(fileMigrations, appliedMigrations)
	if err != nil {
		t.Errorf("Expected no error, but got: %v", err)
	}
	expectedNote := "\nAn applied migration was not found locally"
	if statusNote != expectedNote {
		t.Errorf("Expected statusNote: '%s' \nReceived: %s", expectedNote, statusNote)
	}
	expectedRows := []MigrateStatusRow{
		{
			ID: "1", Name: "foo_bar",
			Status: "applied at 2020-06-13 17:17:44 UTC", Info: "Migration not found locally",
		},
	}
	if diff := pretty.Compare(expectedRows, rows); diff != "" {
		t.Errorf(diff)
	}
}

func TestGetMigrationStatusAppliedGap(t *testing.T) {
	fileMigrations := []FileMigration{
		{ID: "1", Application: "common", Description: "one"},
		{ID: "2", Application: "common", Description: "two"},
		{ID: "3", Application: "common", Description: "three"},
	}
	appliedTime, _ := time.Parse(time.RFC3339, "2020-06-13T17:17:44.371Z")
	appliedMigrations := []AppliedMigration{
		{ID: "1", Name: "one", AppliedAt: appliedTime},
		{ID: "3", Name: "three", AppliedAt: appliedTime},
	}

	rows, statusNote, err := GetMigrationStatus(fileMigrations, appliedMigrations)
	if err != nil {
		t.Errorf("Expected no error, but got: %v", err)
	}
	expectedNote := "\nThere was a gap in the changelog, making it inconsistent"
	if statusNote != expectedNote {
		t.Errorf("Expected statusNote: '%s' \nReceived: %s", expectedNote, statusNote)
	}
	expectedRows := []MigrateStatusRow{
		{
			ID: "1", Name: "one", Application: "common",
			Status: "applied at 2020-06-13 17:17:44 UTC",
		},
		{
			ID: "2", Name: "two", Application: "common",
			Status: "not applied", Info: "Gap in migrations - inconsistency",
		},
		{
			ID: "3", Name: "three", Application: "common",
			Status: "applied at 2020-06-13 17:17:44 UTC",
		},
	}
	if diff := pretty.Compare(expectedRows, rows); diff != "" {
		t.Errorf(diff)
	}
}

func TestGetMigrationStatusEmpty(t *testing.T) {
	appliedMigrations := []AppliedMigration{}
	fileMigrations := []FileMigration{}

	rows, statusNote, err := GetMigrationStatus(fileMigrations, appliedMigrations)
	if err != nil {
		t.Errorf("Expected no error, but got: %v", err)
	}
	if statusNote != "" {
		t.Errorf("Expected empty statusNote but got: %s", statusNote)
	}
	if len(rows) != 0 {
		t.Errorf("Expected no rows but got %v", rows)
	}
}
