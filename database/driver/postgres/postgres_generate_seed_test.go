package postgres

import (
	"fmt"
	"io/ioutil"
	"os"
	"regexp"
	"strings"
	"testing"

	"go-migrations/database"
)

func TestGenerateSeedSQL(t *testing.T) {
	defer resetMockVariables()

	tmpFile, _ := ioutil.TempFile(os.TempDir(), "g-mig-test-")
	defer os.Remove(tmpFile.Name())

	migrations := []database.FileMigration{
		{ID: "1", UpSQL: "SELECT 1"},
		{ID: "2", UpSQL: "SELECT 2"},
	}
	mockableGetFileMigrations = func(a string) ([]database.FileMigration, error) {
		return migrations, nil
	}
	mockableGetBootstrapSQL = func(p string) (string, error) {
		return "SELECT 'bootstrap';", nil
	}

	pg := Postgres{}
	err := pg.GenerateSeedSQL(tmpFile)
	if err != nil {
		t.Errorf("Expected no error, but got: %s", err)
	}

	content, _ := ioutil.ReadFile(tmpFile.Name())
	generatedSeed := string(content)
	expectedSeed := fmt.Sprintf(
		`
			%s
			SELECT 'bootstrap';
			%s;
			INSERT INTO %s (id, name, applied_at) VALUES ('%s', '%s', now());
			%s;
			INSERT INTO %s (id, name, applied_at) VALUES ('%s', '%s', now());
		`,
		createChangelogSQL,
		migrations[0].UpSQL,
		changelogTable, migrations[0].ID, migrations[0].Description,
		migrations[1].UpSQL,
		changelogTable, migrations[1].ID, migrations[1].Description,
	)

	parsedExpected := string(
		regexp.MustCompile(`\s{2,}`).ReplaceAll(
			[]byte(strings.ReplaceAll(expectedSeed, "\n", " ")),
			[]byte(" "),
		),
	)
	parsedReceived := string(
		regexp.MustCompile(`\s{2,}`).ReplaceAll(
			[]byte(strings.ReplaceAll(generatedSeed, "\n", " ")),
			[]byte(" "),
		),
	)

	if parsedReceived != parsedExpected {
		t.Errorf(
			"%s. Sorta expected:\n%s\nReceived:\n%s\n",
			"Did not get expected seed output (after regexing)",
			expectedSeed,
			generatedSeed,
		)
	}
}

func TestGenerateSeedSQLNoBootstrap(t *testing.T) {
	defer resetMockVariables()

	tmpFile, _ := ioutil.TempFile(os.TempDir(), "g-mig-test-")
	defer os.Remove(tmpFile.Name())

	migrations := []database.FileMigration{
		{ID: "1", UpSQL: "SELECT 1"},
	}
	mockableGetFileMigrations = func(a string) ([]database.FileMigration, error) {
		return migrations, nil
	}

	pg := Postgres{}
	err := pg.GenerateSeedSQL(tmpFile)
	if err != nil {
		t.Errorf("Expected no error, but got: %s", err)
	}

	content, _ := ioutil.ReadFile(tmpFile.Name())
	generatedSeed := string(content)
	expectedSeed := fmt.Sprintf(
		`
			%s
			%s;
			INSERT INTO %s (id, name, applied_at) VALUES ('%s', '%s', now());
		`,
		createChangelogSQL,
		migrations[0].UpSQL,
		changelogTable, migrations[0].ID, migrations[0].Description,
	)

	parsedExpected := string(
		regexp.MustCompile(`\s{2,}`).ReplaceAll(
			[]byte(strings.ReplaceAll(expectedSeed, "\n", " ")),
			[]byte(" "),
		),
	)
	parsedReceived := string(
		regexp.MustCompile(`\s{2,}`).ReplaceAll(
			[]byte(strings.ReplaceAll(generatedSeed, "\n", " ")),
			[]byte(" "),
		),
	)

	if parsedReceived != parsedExpected {
		t.Errorf(
			"%s. Sorta expected:\n%s\nReceived:\n%s\n",
			"Did not get expected seed output (after regexing)",
			expectedSeed,
			generatedSeed,
		)
	}
}
