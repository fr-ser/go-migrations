package config

import (
	"errors"
	"fmt"
	"io/ioutil"
	"os"

	"gopkg.in/yaml.v2"
)

// fileConfig is a direct translation of the config YAML into a struct
type fileConfig struct {
	DbType   string `yaml:"db_type"`
	Prepare  bool   `yaml:"prepare"`
	Host     string `yaml:"host"`
	Port     uint16 `yaml:"port"`
	DbName   string `yaml:"db_name"`
	User     string `yaml:"user"`
	Password string `yaml:"password"`
}

// Config stores configuration for database environment like host, port
// and also migration parameters like whether to apply prepare scripts
type Config struct {
	MigrationsPath      string
	Environment         string
	ApplyPrepareScripts bool
	ChangelogName       string
	Db                  struct {
		Type     string
		Host     string
		Port     uint16
		Name     string
		User     string
		Password string
	}
}

// LoadConfig takes a path to a configuration file reads it
// and performs validity checks
func LoadConfig(configPath, migrationsPath, environment string) (Config, error) {
	databaseConfig, err := unmarshalConfig(configPath)
	if err != nil {
		return databaseConfig, err
	}

	databaseConfig.ChangelogName = "migrations_changelog"
	databaseConfig.MigrationsPath = migrationsPath
	databaseConfig.Environment = environment

	if err := validateConfig(databaseConfig); err != nil {
		return databaseConfig, err
	}

	return databaseConfig, nil
}

func unmarshalConfig(path string) (Config, error) {
	databaseConfig := Config{}
	fConfig := fileConfig{}

	configFile, err := os.Open(path)
	if err != nil {
		return databaseConfig, fmt.Errorf("Couldn't open config file: %v", err)
	}
	defer configFile.Close()
	fileContent, err := ioutil.ReadAll(configFile)
	if err != nil {
		return databaseConfig, fmt.Errorf("Couldn't read config file: %v", err)
	}

	err = yaml.UnmarshalStrict(fileContent, &fConfig)
	if err != nil {
		return databaseConfig, fmt.Errorf("Couldn't unmarshal yaml: %v", err)
	}

	databaseConfig.ApplyPrepareScripts = fConfig.Prepare
	databaseConfig.Db.Type = fConfig.DbType
	databaseConfig.Db.Host = fConfig.Host
	databaseConfig.Db.Port = fConfig.Port
	databaseConfig.Db.Name = fConfig.DbName
	databaseConfig.Db.User = fConfig.User
	databaseConfig.Db.Password = fConfig.Password

	return databaseConfig, nil
}

func validateConfig(config Config) error {
	if config.Db.Port == 0 {
		return errors.New("No port specified, or invalid port of 0")
	}
	if config.Db.Host == "" {
		return errors.New("No host specified")
	}
	if config.Db.Name == "" {
		return errors.New("No database name specified")
	}
	if config.Db.User == "" {
		return errors.New("No username specified")
	}
	return nil
}
