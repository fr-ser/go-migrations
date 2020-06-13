package direction

// MigrateDirection is the direction of a migration: up or down
type MigrateDirection bool

const (
	// Up applies the initial migration
	Up MigrateDirection = false
	// Down removes the initial migration
	Down MigrateDirection = true
)

// Directions is a list of all possible directions to iterate
var Directions = []struct {
	Direction MigrateDirection
	Name      string
}{
	{Up, "Up"},
	{Down, "Down"},
}
