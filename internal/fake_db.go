package internal

import (
	"os"
	"testing"
	"time"

	"github.com/jedib0t/go-pretty/v6/progress"
	"github.com/kylelemons/godebug/pretty"

	"go-migrations/database"
	"go-migrations/database/config"
	"go-migrations/internal/direction"
)

type applyMigrationsWithCountArgs struct {
	count     uint
	all       bool
	direction direction.MigrateDirection
}

type applySpecificMigrationArgs struct {
	filter    string
	direction direction.MigrateDirection
}

// FakeDbWithSpy implements the database interface and saves method calls
type FakeDbWithSpy struct {
	initCalls                       []bool
	bootstrapCalls                  []bool
	waitForStartCalls               []bool
	ensureMigrationsChangelogCalls  []bool
	ensureConsistentMigrationsCalls []bool
	applyAllUpMigrationsCalls       []bool
	getFileMigrationsCalls          []bool
	getAppliedMigrationsCalls       []bool
	generateSeedSQLCalls            []bool
	applyMigrationsWithCountCalls   []applyMigrationsWithCountArgs
	applySpecificMigrationCalls     []applySpecificMigrationArgs
}

// WaitForStart saves the call
func (db *FakeDbWithSpy) WaitForStart(pollInterval time.Duration, retryCount int) error {
	db.waitForStartCalls = append(db.waitForStartCalls, true)
	return nil
}

// AssertWaitForStartCalled checks for calls
func (db *FakeDbWithSpy) AssertWaitForStartCalled(t *testing.T, expectCalled bool) {
	wasCalled := len(db.waitForStartCalls) > 0

	if wasCalled && !expectCalled {
		t.Errorf("WaitForStart was called but shouldn't have been")
	} else if !wasCalled && expectCalled {
		t.Errorf("WaitForStart wasn't called but should have been")
	}
}

// Bootstrap saves the call
func (db *FakeDbWithSpy) Bootstrap() error {
	db.bootstrapCalls = append(db.bootstrapCalls, true)
	return nil
}

// AssertBootstrapCalled checks for calls
func (db *FakeDbWithSpy) AssertBootstrapCalled(t *testing.T, expectCalled bool) {
	wasCalled := len(db.bootstrapCalls) > 0

	if wasCalled && !expectCalled {
		t.Errorf("Bootstrap was called but shouldn't have been")
	} else if !wasCalled && expectCalled {
		t.Errorf("Bootstrap wasn't called but should have been")
	}
}

// GetFileMigrations saves the call
func (db *FakeDbWithSpy) GetFileMigrations() ([]database.FileMigration, error) {
	db.getFileMigrationsCalls = append(db.getFileMigrationsCalls, true)
	return nil, nil
}

// AssertGetFileMigrationsCalled checks for calls
func (db *FakeDbWithSpy) AssertGetFileMigrationsCalled(t *testing.T, expectCalled bool) {
	wasCalled := len(db.getFileMigrationsCalls) > 0

	if wasCalled && !expectCalled {
		t.Errorf("GetFileMigrations was called but shouldn't have been")
	} else if !wasCalled && expectCalled {
		t.Errorf("GetFileMigrations wasn't called but should have been")
	}
}

// GetAppliedMigrations saves the call
func (db *FakeDbWithSpy) GetAppliedMigrations() ([]database.AppliedMigration, error) {
	db.getAppliedMigrationsCalls = append(db.getAppliedMigrationsCalls, true)
	return nil, nil
}

// AssertGetAppliedMigrationsCalled checks for calls
func (db *FakeDbWithSpy) AssertGetAppliedMigrationsCalled(t *testing.T, expectCalled bool) {
	wasCalled := len(db.getAppliedMigrationsCalls) > 0

	if wasCalled && !expectCalled {
		t.Errorf("GetAppliedMigrations was called but shouldn't have been")
	} else if !wasCalled && expectCalled {
		t.Errorf("GetAppliedMigrations wasn't called but should have been")
	}
}

// ApplyAllUpMigrations saves the call
func (db *FakeDbWithSpy) ApplyAllUpMigrations(pw progress.Writer) error {
	db.applyAllUpMigrationsCalls = append(db.applyAllUpMigrationsCalls, true)
	return nil
}

// AssertApplyAllUpMigrationsCalled checks for calls
func (db *FakeDbWithSpy) AssertApplyAllUpMigrationsCalled(t *testing.T, expectCalled bool) {
	wasCalled := len(db.applyAllUpMigrationsCalls) > 0

	if wasCalled && !expectCalled {
		t.Errorf("ApplyAllUpMigrations was called but shouldn't have been")
	} else if !wasCalled && expectCalled {
		t.Errorf("ApplyAllUpMigrations wasn't called but should have been")
	}
}

// EnsureConsistentMigrations checks for inconsistencies in the changelog
func (db *FakeDbWithSpy) EnsureConsistentMigrations() error {
	db.ensureConsistentMigrationsCalls = append(db.ensureConsistentMigrationsCalls, true)
	return nil
}

// AssertEnsureConsistentMigrationsCalled checks for at least one call
func (db *FakeDbWithSpy) AssertEnsureConsistentMigrationsCalled(t *testing.T, expectCalled bool) {
	wasCalled := len(db.ensureConsistentMigrationsCalls) > 0

	if wasCalled && !expectCalled {
		t.Errorf("EnsureConsistentMigrations was called but shouldn't have been")
	} else if !wasCalled && expectCalled {
		t.Errorf("EnsureConsistentMigrations wasn't called but should have been")
	}
}

// EnsureMigrationsChangelog saves the call
func (db *FakeDbWithSpy) EnsureMigrationsChangelog() (bool, error) {
	db.ensureMigrationsChangelogCalls = append(db.ensureMigrationsChangelogCalls, true)
	return false, nil
}

// AssertEnsureMigrationsChangelogCalled checks for calls
func (db *FakeDbWithSpy) AssertEnsureMigrationsChangelogCalled(t *testing.T, expectCalled bool) {
	wasCalled := len(db.ensureMigrationsChangelogCalls) > 0

	if wasCalled && !expectCalled {
		t.Errorf("EnsureMigrationsChangelog was called but shouldn't have been")
	} else if !wasCalled && expectCalled {
		t.Errorf("EnsureMigrationsChangelog wasn't called but should have been")
	}
}

// GenerateSeedSQL saves the call
func (db *FakeDbWithSpy) GenerateSeedSQL(f *os.File) error {
	db.generateSeedSQLCalls = append(db.generateSeedSQLCalls, true)
	return nil
}

// AssertGenerateSeedSQLCalled checks for calls
func (db *FakeDbWithSpy) AssertGenerateSeedSQLCalled(t *testing.T, expectCalled bool) {
	wasCalled := len(db.generateSeedSQLCalls) > 0

	if wasCalled && !expectCalled {
		t.Errorf("GenerateSeedSQL was called but shouldn't have been")
	} else if !wasCalled && expectCalled {
		t.Errorf("GenerateSeedSQL wasn't called but should have been")
	}
}

// Init saves the call
func (db *FakeDbWithSpy) Init(_ config.Config) error {
	db.initCalls = append(db.initCalls, true)
	return nil
}

// ApplySpecificMigration applies one up migration by a filter
func (db *FakeDbWithSpy) ApplySpecificMigration(filter string, direction direction.MigrateDirection) error {
	db.applySpecificMigrationCalls = append(
		db.applySpecificMigrationCalls,
		applySpecificMigrationArgs{filter: filter, direction: direction},
	)
	return nil
}

// AssertApplySpecificMigrationCalled checks for at least one call
func (db *FakeDbWithSpy) AssertApplySpecificMigrationCalled(t *testing.T, expectCalled bool) (
	wasCalled bool,
) {
	wasCalled = len(db.applySpecificMigrationCalls) > 0

	if wasCalled && !expectCalled {
		t.Errorf("ApplySpecificMigration was called but shouldn't have been")
	} else if !wasCalled && expectCalled {
		t.Errorf("ApplySpecificMigration wasn't called but should have been")
	}
	return wasCalled
}

// AssertApplySpecificMigrationCalledWith checks the arguments of the last call
func (db *FakeDbWithSpy) AssertApplySpecificMigrationCalledWith(
	t *testing.T, filter string, direction direction.MigrateDirection,
) {
	if !db.AssertApplySpecificMigrationCalled(t, true) {
		return
	}
	lastCall := db.applySpecificMigrationCalls[len(db.applySpecificMigrationCalls)-1]
	expectedCall := applySpecificMigrationArgs{filter: filter, direction: direction}
	if lastCall != expectedCall {
		t.Errorf(
			"ApplySpecificMigration was called with '%v' instead of '%v'",
			lastCall, expectedCall,
		)
	}
}

// ApplyMigrationsWithCount applies Up migration by a count
func (db *FakeDbWithSpy) ApplyMigrationsWithCount(
	count uint, all bool, dir direction.MigrateDirection,
) error {
	db.applyMigrationsWithCountCalls = append(
		db.applyMigrationsWithCountCalls,
		applyMigrationsWithCountArgs{count: count, all: all, direction: dir},
	)
	return nil
}

// AssertApplyMigrationsWithCountCalled checks for at least one call
func (db *FakeDbWithSpy) AssertApplyMigrationsWithCountCalled(t *testing.T, expectCalled bool) (
	wasCalled bool,
) {
	wasCalled = len(db.applyMigrationsWithCountCalls) > 0

	if wasCalled && !expectCalled {
		t.Errorf("ApplyMigrationsWithCount was called but shouldn't have been")
	} else if !wasCalled && expectCalled {
		t.Errorf("ApplyMigrationsWithCount wasn't called but should have been")
	}
	return wasCalled

}

// AssertApplyMigrationsWithCountCalledWith checks the arguments of the last call
func (db *FakeDbWithSpy) AssertApplyMigrationsWithCountCalledWith(
	t *testing.T, count uint, all bool, dir direction.MigrateDirection,
) {
	if !db.AssertApplyMigrationsWithCountCalled(t, true) {
		return
	}
	lastCall := db.applyMigrationsWithCountCalls[len(db.applyMigrationsWithCountCalls)-1]
	expectedCall := applyMigrationsWithCountArgs{
		count: count, all: all, direction: dir,
	}
	if diff := pretty.Compare(lastCall, expectedCall); diff != "" {
		t.Errorf(
			"ApplyMigrationsWithCount was called with %+v instead of %+v",
			lastCall, expectedCall,
		)
	}
}
