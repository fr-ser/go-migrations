package internal

import (
	"testing"
	"time"

	"go-migrations/database/config"
)

type applyUpMigrationsWithCountArgs struct {
	count int
	all   bool
}

// FakeDbWithSpy implements the database interface and saves method calls
type FakeDbWithSpy struct {
	WaitForStartCalled              bool
	BootstrapCalled                 bool
	ApplyAllUpMigrationsCalled      bool
	EnsureMigrationsChangelogCalled bool
	InitCalled                      bool

	applyUpMigrationsWithCountCalls []applyUpMigrationsWithCountArgs
	applySpecificUpMigrationCalls   []string
}

// WaitForStart saves the call
func (db *FakeDbWithSpy) WaitForStart(pollInterval time.Duration, retryCount int) error {
	db.WaitForStartCalled = true
	return nil
}

// Bootstrap saves the call
func (db *FakeDbWithSpy) Bootstrap() error {
	db.BootstrapCalled = true
	return nil
}

// ApplyAllUpMigrations saves the call
func (db *FakeDbWithSpy) ApplyAllUpMigrations() error {
	db.ApplyAllUpMigrationsCalled = true
	return nil
}

// EnsureMigrationsChangelog saves the call
func (db *FakeDbWithSpy) EnsureMigrationsChangelog() (bool, error) {
	db.EnsureMigrationsChangelogCalled = true
	return false, nil
}

// Init saves the call
func (db *FakeDbWithSpy) Init(_ config.Config) error {
	db.InitCalled = true
	return nil
}

// ApplySpecificUpMigration applies one up migration by a filter
func (db *FakeDbWithSpy) ApplySpecificUpMigration(filter string) error {
	db.applySpecificUpMigrationCalls = append(db.applySpecificUpMigrationCalls, filter)
	return nil
}

// ApplySpecificUpMigrationCalled checks for at least one call
func (db *FakeDbWithSpy) ApplySpecificUpMigrationCalled(t *testing.T) bool {
	if len(db.applySpecificUpMigrationCalls) == 0 {
		t.Errorf("ApplySpecificUpMigration was not called")
		return false
	}
	return true
}

// ApplySpecificUpMigrationNotCalled checks for at least one call
func (db *FakeDbWithSpy) ApplySpecificUpMigrationNotCalled(t *testing.T) bool {
	if len(db.applySpecificUpMigrationCalls) != 0 {
		t.Errorf("ApplySpecificUpMigration was called")
		return false
	}
	return true
}

// ApplySpecificUpMigrationCalledWith checks the arguments of the last call
func (db *FakeDbWithSpy) ApplySpecificUpMigrationCalledWith(t *testing.T, filter string) bool {
	if !db.ApplySpecificUpMigrationCalled(t) {
		return false
	}
	lastCall := db.applySpecificUpMigrationCalls[len(db.applySpecificUpMigrationCalls)-1]
	if lastCall != filter {
		t.Errorf(
			"ApplySpecificUpMigration was called with '%s' instead of '%s'", lastCall, filter,
		)
		return false
	}
	return true
}

// ApplyUpMigrationsWithCount applies up migration by a count
func (db *FakeDbWithSpy) ApplyUpMigrationsWithCount(count int, all bool) error {
	db.applyUpMigrationsWithCountCalls = append(
		db.applyUpMigrationsWithCountCalls,
		applyUpMigrationsWithCountArgs{count: count, all: all},
	)
	return nil
}

// ApplyUpMigrationsWithCountCalled checks for at least one call
func (db *FakeDbWithSpy) ApplyUpMigrationsWithCountCalled(t *testing.T) bool {
	if len(db.applyUpMigrationsWithCountCalls) == 0 {
		t.Errorf("ApplyUpMigrationsWithCount was not called")
		return false
	}
	return true
}

// ApplyUpMigrationsWithCountNotCalled checks for at least one call
func (db *FakeDbWithSpy) ApplyUpMigrationsWithCountNotCalled(t *testing.T) bool {
	if len(db.applyUpMigrationsWithCountCalls) != 0 {
		t.Errorf("ApplyUpMigrationsWithCount was called")
		return false
	}
	return true
}

// ApplyUpMigrationsWithCountCalledWith checks the arguments of the last call
func (db *FakeDbWithSpy) ApplyUpMigrationsWithCountCalledWith(t *testing.T, count int, all bool) bool {
	if !db.ApplyUpMigrationsWithCountCalled(t) {
		return false
	}
	lastCall := db.applyUpMigrationsWithCountCalls[len(db.applyUpMigrationsWithCountCalls)-1]
	if lastCall.all != all || lastCall.count != count {
		t.Errorf(
			"ApplyUpMigrationsWithCount was called with %+v instead of %+v",
			lastCall,
			applyUpMigrationsWithCountArgs{all: all, count: count},
		)
		return false
	}
	return true
}
