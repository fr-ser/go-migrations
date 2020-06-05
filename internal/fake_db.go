package internal

import (
	"time"

	"go-migrations/database/config"
)

// FakeDbWithSpy implements the database interface and saves method calls
type FakeDbWithSpy struct {
	WaitForStartCalled              bool
	BootstrapCalled                 bool
	ApplyAllUpMigrationsCalled      bool
	EnsureMigrationsChangelogCalled bool
	InitCalled                      bool
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
