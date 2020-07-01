package database

import (
	"fmt"
	"os"
	"sort"

	"github.com/jedib0t/go-pretty/v6/table"
	"github.com/jedib0t/go-pretty/v6/text"
)

func colorizeRow(row table.Row) text.Colors {
	if row[4] != "" {
		return text.Colors{text.Reset, text.FgHiYellow}
	}
	return nil
}

// PrintStatusTable prints the migration status in a table format
// optimized for human readability
func PrintStatusTable(rows []MigrateStatusRow, statusNote string) {
	t := table.NewWriter()
	t.SetRowPainter(colorizeRow)
	t.SetOutputMirror(os.Stdout)
	t.AppendHeader(table.Row{"ID", "Name", "App", "Status", "Info"})

	for _, row := range rows {
		t.AppendRow([]interface{}{
			row.ID, row.Name, row.Application, row.Status, row.Info,
		})
	}

	t.Render()
	fmt.Println(text.FgHiRed.Sprint(statusNote))
}

// MigrateStatusRow is the status of one applied/local migration
type MigrateStatusRow struct {
	ID          string
	Name        string
	Application string
	Status      string
	Info        string
}

// MigrateStatusHeader is the header for the status table
var MigrateStatusHeader = []string{"ID", "Name", "Application", "Status", "Info"}

// GetMigrationStatus returns the status of the current migrations
// This is done by combining the information from the changelog and the local files
func GetMigrationStatus(fileMigrations []FileMigration, appliedMigrations []AppliedMigration) (
	rows []MigrateStatusRow, statusNote string, err error,
) {
	ids := map[string]bool{}
	fileLookup := map[string]FileMigration{}
	dbLookup := map[string]AppliedMigration{}

	var localNotFound bool
	var inconsistentLog bool

	for _, fileMig := range fileMigrations {
		fileLookup[fileMig.ID] = fileMig
		ids[fileMig.ID] = true
	}
	for _, dbMig := range appliedMigrations {
		dbLookup[dbMig.ID] = dbMig
		ids[dbMig.ID] = true
	}

	for id := range ids {
		row := MigrateStatusRow{ID: id}
		_, fileExists := fileLookup[id]
		_, applied := dbLookup[id]

		if fileExists {
			row.Name = fileLookup[id].Description
			row.Application = fileLookup[id].Application
		} else {
			row.Name = dbLookup[id].Name
			row.Info = "Migration not found locally"
			localNotFound = true
		}

		if applied {
			timeString := dbLookup[id].AppliedAt.Format("2006-01-02 15:04:05")
			row.Status = fmt.Sprintf("applied at %s UTC", timeString)
		} else {
			row.Status = "not applied"
		}

		rows = append(rows, row)
	}

	sort.Slice(rows, func(i, j int) bool {
		return rows[i].ID < rows[j].ID
	})

	// we ignore the last entry as it cannot be inconsistent
	for idx := 0; idx < len(rows)-1; idx++ {
		if rows[idx].Status == "not applied" && rows[idx+1].Status != "not applied" {
			inconsistentLog = true
			rows[idx].Info = "Gap in migrations - inconsistency"
		}
	}

	if localNotFound {
		statusNote += "\nAn applied migration was not found locally"
	}

	if inconsistentLog {
		statusNote += "\nThere was a gap in the changelog, making it inconsistent"
	}

	return rows, statusNote, nil
}
