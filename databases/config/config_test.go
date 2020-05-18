package config

import (
	"fmt"
	"io/ioutil"
	"regexp"
	"syscall"
	"testing"

	"github.com/kylelemons/godebug/pretty"
)

func assertEqual(t *testing.T, a, b interface{}) {
	if diff := pretty.Compare(a, b); diff != "" {
		t.Errorf("The data was not the same:\n%s", diff)
	}
}

var validConfigYaml = `
db_type: postgres
prepare: True
host: localhost
port: 35432
db_name: zlab
user: db_admin
password: pass
`

func TestLoadValidConfig(t *testing.T) {
	f, _ := ioutil.TempFile("", "tmp_file")
	defer syscall.Unlink(f.Name())

	f.WriteString(validConfigYaml)

	expectedConfig := Config{}
	expectedConfig.ChangelogName = "migrations_changelog"
	expectedConfig.DockerComposeFile = "my-compose.yaml"
	expectedConfig.MigrationsPath = "./migrations"
	expectedConfig.Environment = "test_env"
	expectedConfig.ApplyPrepareScripts = true
	expectedConfig.Db.Type = "postgres"
	expectedConfig.Db.Host = "localhost"
	expectedConfig.Db.Port = 35432
	expectedConfig.Db.Name = "zlab"
	expectedConfig.Db.User = "db_admin"
	expectedConfig.Db.Password = "pass"

	config, err := LoadConfig(f.Name(), "my-compose.yaml", "./migrations", "test_env")
	if err != nil {
		t.Errorf("Returned error: %v", err)
	}

	assertEqual(t, config, expectedConfig)

}

var invalidConfigFiles = []struct{ name, file string }{
	{"missing port", configWithoutLineFor("port")},
	{"missing host", configWithoutLineFor("host")},
	{"missing database name", configWithoutLineFor("db_name")},
	{"missing user", configWithoutLineFor("user")},
}

func TestInvalidConfigFile(t *testing.T) {
	for _, configFile := range invalidConfigFiles {
		f, _ := ioutil.TempFile("", "tmp_file")
		defer syscall.Unlink(f.Name())
		f.WriteString(configFile.file)

		t.Run(configFile.name, func(t *testing.T) {
			_, err := LoadConfig(f.Name(), "", "", "")
			if err == nil {
				t.Errorf("Got no error for: %s", configFile.name)
			}
		})
	}
}

func configWithoutLineFor(key string) string {
	regexMatcher := fmt.Sprintf("(?m)%s:.+$", key)
	return regexp.MustCompile(regexMatcher).ReplaceAllString(validConfigYaml, "")
}
