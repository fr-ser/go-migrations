package database

import (
	"testing"
	"time"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/kylelemons/godebug/pretty"
	"github.com/lithammer/dedent"
)

func TestGetAppliedMigrations(t *testing.T) {
	db, mock, err := sqlmock.New(sqlmock.QueryMatcherOption(sqlmock.QueryMatcherEqual))
	mockRows := sqlmock.NewRows([]string{"id", "name", "applied_at"})

	time1, _ := time.Parse(time.RFC3339, "2014-11-12T11:45:26.371Z")
	time2, _ := time.Parse(time.RFC3339, "2015-12-11T10:46:23.378Z")
	expectedMigrations := []AppliedMigration{
		{ID: "20171101000001", Name: "foo", AppliedAt: time1},
		{ID: "20171101000002", Name: "bar", AppliedAt: time2},
	}

	for _, mig := range expectedMigrations {
		mockRows.AddRow(mig.ID, mig.Name, mig.AppliedAt)
	}

	mock.ExpectQuery(dedent.Dedent(`
		SELECT id, name, applied_at
		FROM schema.changelog
		ORDER BY id ASC
	`)).WillReturnRows(mockRows)

	gotMigrations, err := GetAppliedMigrations(db, "schema.changelog")
	if err != nil {
		t.Fatalf("Got an error loading migrations: %v", err)
	}

	diff := pretty.Compare(expectedMigrations, gotMigrations)
	if diff != "" {
		t.Error(diff)
	}
}
