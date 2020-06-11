package internal

import (
	"testing"
	"time"

	"go-migrations/database/config"
)

// TODO: change to "assertWaitForStartCalled(true)"

type applyUpMigrationsWithCountArgs struct {
	count uint
	all   bool
}

// FakeDbWithSpy implements the database interface and saves method calls
type FakeDbWithSpy struct {
	initCalls                       []bool
	bootstrapCalls                  []bool
	waitForStartCalls               []bool
	ensureMigrationsChangelogCalls  []bool
	ensureConsistentMigrationsCalls []bool
	applyAllUpMigrationsCalls       []bool
	applyUpMigrationsWithCountCalls []applyUpMigrationsWithCountArgs
	applySpecificUpMigrationCalls   []string
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

// ApplyAllUpMigrations saves the call
func (db *FakeDbWithSpy) ApplyAllUpMigrations() error {
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

// Init saves the call
func (db *FakeDbWithSpy) Init(_ config.Config) error {
	db.initCalls = append(db.initCalls, true)
	return nil
}

// ApplySpecificUpMigration applies one up migration by a filter
func (db *FakeDbWithSpy) ApplySpecificUpMigration(filter string) error {
	db.applySpecificUpMigrationCalls = append(db.applySpecificUpMigrationCalls, filter)
	return nil
}

// AssertApplySpecificUpMigrationCalled checks for at least one call
func (db *FakeDbWithSpy) AssertApplySpecificUpMigrationCalled(t *testing.T, expectCalled bool) (
	wasCalled bool,
) {
	wasCalled = len(db.applySpecificUpMigrationCalls) > 0

	if wasCalled && !expectCalled {
		t.Errorf("ApplySpecificUpMigration was called but shouldn't have been")
	} else if !wasCalled && expectCalled {
		t.Errorf("ApplySpecificUpMigration wasn't called but should have been")
	}
	return wasCalled
}

// AssertApplySpecificUpMigrationCalledWith checks the arguments of the last call
func (db *FakeDbWithSpy) AssertApplySpecificUpMigrationCalledWith(t *testing.T, filter string) {
	if !db.AssertApplySpecificUpMigrationCalled(t, true) {
		return
	}
	lastCall := db.applySpecificUpMigrationCalls[len(db.applySpecificUpMigrationCalls)-1]
	if lastCall != filter {
		t.Errorf(
			"ApplySpecificUpMigration was called with '%s' instead of '%s'", lastCall, filter,
		)
	}
}

// ApplyUpMigrationsWithCount applies up migration by a count
func (db *FakeDbWithSpy) ApplyUpMigrationsWithCount(count uint, all bool) error {
	db.applyUpMigrationsWithCountCalls = append(
		db.applyUpMigrationsWithCountCalls,
		applyUpMigrationsWithCountArgs{count: count, all: all},
	)
	return nil
}

// AssertApplyUpMigrationsWithCountCalled checks for at least one call
func (db *FakeDbWithSpy) AssertApplyUpMigrationsWithCountCalled(t *testing.T, expectCalled bool) (
	wasCalled bool,
) {
	wasCalled = len(db.applyUpMigrationsWithCountCalls) > 0

	if wasCalled && !expectCalled {
		t.Errorf("ApplyUpMigrationsWithCount was called but shouldn't have been")
	} else if !wasCalled && expectCalled {
		t.Errorf("ApplyUpMigrationsWithCount wasn't called but should have been")
	}
	return wasCalled

}

// AssertApplyUpMigrationsWithCountCalledWith checks the arguments of the last call
func (db *FakeDbWithSpy) AssertApplyUpMigrationsWithCountCalledWith(t *testing.T, count uint,
	all bool,
) {
	if !db.AssertApplyUpMigrationsWithCountCalled(t, true) {
		return
	}
	lastCall := db.applyUpMigrationsWithCountCalls[len(db.applyUpMigrationsWithCountCalls)-1]
	if lastCall.all != all || lastCall.count != count {
		t.Errorf(
			"ApplyUpMigrationsWithCount was called with %+v instead of %+v",
			lastCall,
			applyUpMigrationsWithCountArgs{all: all, count: count},
		)
	}
}
