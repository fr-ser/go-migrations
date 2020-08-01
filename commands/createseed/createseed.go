package createseed

import (
	"fmt"
	"os"

	"github.com/urfave/cli/v2"

	"go-migrations/commands"
	"go-migrations/database/driver"
	"go-migrations/utils"
)

var (
	mockableLoadDB = driver.LoadDB
)

var flags = []cli.Flag{
	&cli.StringFlag{
		Name: "migrations-path", Aliases: []string{"p"}, Value: "./migrations/zlab",
		Usage: "(relative) path to the folder containing the database migrations",
	},
	&cli.StringFlag{
		Name: "environment", Aliases: []string{"e"}, Value: "development",
		Usage: "Name of the environment and the corresponding configuration",
	},
	&cli.StringFlag{
		Name: "target", Aliases: []string{"t"}, Value: "seed.sql",
		Usage: "Name and path of the file containing the seed sql",
	},
}

// CreateSeedCommand creates an SQL file that can be used to seed the database directly
var CreateSeedCommand = &cli.Command{
	Name:   "create-seed",
	Usage:  "creates an SQL file that can be used to seed the database directly",
	Flags:  flags,
	Before: commands.NoArguments,
	Action: func(c *cli.Context) (err error) {
		exists, err := utils.FileExists(c.String("target"))
		if err != nil {
			return fmt.Errorf("Could not check existence of %s: %v", c.String("target"), err)
		} else if exists {
			return fmt.Errorf("The file %s already exists", c.String("target"))
		}

		db, err := mockableLoadDB(c.String("migrations-path"), c.String("environment"))
		if err != nil {
			return err
		}

		target, err := os.Create(c.String("target"))
		defer target.Close()
		if err != nil {
			return fmt.Errorf("Could not create the target file: %v", err)
		}

		err = db.GenerateSeedSQL(target)
		if err != nil {
			return err
		}

		return nil
	},
}
