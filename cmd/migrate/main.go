package main

import (
	"database/sql"
	"errors"
	"flag"
	"fmt"

	"github.com/devesh2997/consequent/cmd/flags"
	"github.com/devesh2997/consequent/config"
	"github.com/devesh2997/consequent/datasources"

	"github.com/devesh2997/consequent/migrate"
)

var errNoValidOptionsOrArgs = errors.New("no valid args or options found")

const actionUp = "up"
const actionDown = "down"
const actionForceUp = "forceUp"
const actionForceDown = "forceDown"

func main() {
	env := flags.GetEnvironment()
	config.LoadConfig(env, ".")
	sqlDB, err := getSQLDB()
	if err != nil {
		fmt.Println(err)
		return
	}

	migrationsConfig := getMigrationsConfig(sqlDB)

	migrator, err := migrate.New(migrationsConfig)
	if err != nil {
		fmt.Println(err)
		return
	}

	flags.Parse()

	actionArg := flag.Arg(0)

	if *flags.Create != "" {
		err = migrator.Create(*flags.Create)
	} else if *flags.Steps != 0 {
		err = migrator.Steps(*flags.Steps)
	} else if actionArg == actionUp {
		err = migrator.Up()
	} else if actionArg == actionDown {
		err = migrator.Down()
	} else if actionArg == actionForceUp {
		err = migrator.ForceUp()
	} else if actionArg == actionForceDown {
		err = migrator.ForceDown()
	} else {
		err = errNoValidOptionsOrArgs
	}

	if err != nil {
		fmt.Println(err)
	}
}

func getSQLDB() (*sql.DB, error) {
	ds, err := datasources.Get()
	if err != nil {
		return nil, err
	}

	sqlxClient := ds.SQLClients.GetSQLXClusterDB()
	masterDB := sqlxClient.Master().DB

	return masterDB, nil
}

func getMigrationsConfig(sqlDB *sql.DB) migrate.Config {
	return migrate.Config{
		MigrationsDir:      "migrations",
		DBInstance:         sqlDB,
		MigrationsDatabase: config.Config.SQL.Master.DB,
	}
}
