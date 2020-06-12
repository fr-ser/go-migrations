package driver

import (
	"testing"

	"go-migrations/database/config"
	"go-migrations/internal"
)

var loadConfigCall []string

func fakeLoadConfigWithSpy(configPath, migrationsPath, environment string) (config.Config, error) {
	loadConfigCall = []string{configPath, migrationsPath, environment}
	return config.Config{}, nil
}
func TestLoadDb(t *testing.T) {
	loadConfig = fakeLoadConfigWithSpy

	db, err := LoadDB("./mig_path", "my_env")
	if err != nil {
		t.Errorf("Returned error: %v", err)
	}
	if db == nil {
		t.Error("Returned database was nil")
	}

	expected := []string{"./mig_path/_environments/my_env.yaml", "./mig_path", "my_env"}
	if !internal.StrSliceEqual(loadConfigCall, expected) {
		t.Errorf("Expected arguments '%v', but got %v", expected, loadConfigCall)
	}
}
