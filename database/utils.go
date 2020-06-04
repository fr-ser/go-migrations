package database

import (
	"database/sql"
	"fmt"
	"time"
)

// WaitForStart tries to connect to the database
// parameters are the number of retries and the sleep interval in milliseconds between the retries
func WaitForStart(db *sql.DB, pollInterval time.Duration, retries int) error {
	var err error

	for retry := 0; retry < retries; retry++ {
		_, err = db.Exec("SELECT 1")
		if err == nil {
			return nil
		}
		time.Sleep(pollInterval)
	}

	return fmt.Errorf("Timed out connecting to database: %v", err)
}
