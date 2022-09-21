package migrate

import (
	"database/sql"
	_ "database/sql/driver"
	"errors"
	"fmt"
	"os"
	"time"

	"github.com/golang-migrate/migrate/v4"
	"github.com/golang-migrate/migrate/v4/database/mysql"
	_ "github.com/golang-migrate/migrate/v4/source/file"
)

var ErrMigrationsDirRequired = errors.New("migrations directory is required")
var ErrDBInstanceRequired = errors.New("db instance is required")
var ErrMigrationsDatabaseRequired = errors.New("migrations database is required")

type Config struct {
	// MigrationsDir is the path of the directory where migration files are present
	MigrationsDir string
	// DBInstance is the client to be used for running migrations on
	DBInstance *sql.DB
	// MigrationsTable is the migrations table that will be used to store migration history
	// If not provided, "schema_migrations" will be used as default
	MigrationsTable string
	// MigrationsDatabase is the database on which migrations will be run
	MigrationsDatabase string
}

func (config Config) validate() error {
	if config.MigrationsDir == "" {
		return ErrMigrationsDirRequired
	}

	if config.DBInstance == nil {
		return ErrDBInstanceRequired
	}

	if config.MigrationsDatabase == "" {
		return ErrMigrationsDatabaseRequired
	}

	return nil
}

type Migrator struct {
	config  Config
	migrate *migrate.Migrate
}

func New(config Config) (*Migrator, error) {
	if err := config.validate(); err != nil {
		return nil, err
	}

	mysqlConfig := &mysql.Config{
		DatabaseName: config.MigrationsDatabase,
		NoLock:       false,
	}
	if config.MigrationsTable != "" {
		mysqlConfig.MigrationsTable = config.MigrationsTable
	}

	driver, _ := mysql.WithInstance(config.DBInstance, mysqlConfig)

	migrate, err := migrate.NewWithDatabaseInstance(
		fmt.Sprintf("file://%s", config.MigrationsDir),
		"mysql",
		driver,
	)

	if err != nil {
		return nil, err
	}

	return &Migrator{
		config:  config,
		migrate: migrate,
	}, nil
}

func (m Migrator) Up() error {
	return m.migrate.Up()
}

func (m Migrator) Down() error {
	return m.migrate.Down()
}

// Steps looks at the currently active migration version. It will migrate up if n > 0, and down if n < 0
func (m Migrator) Steps(n int) error {
	return m.migrate.Steps(n)
}

//first checks the cuurent version in which dirty is shown
//then using force function it changes dirty = 1 to dirty = 0
//then using steps function to StepDown to the previous stable version
func (m *Migrator) ForceUp() error {
	version, dirty, err := m.migrate.Version()
	if err != nil {
		return err
	}
	if err == nil && dirty {
		if err := m.migrate.Force(int(version)); err != nil {
			return err
		}

		if err = m.migrate.Steps(-1); err != nil {
			return err
		}

		return nil
	}

	return fmt.Errorf("no dirty version, invalid use of force command")
}

//first checks the cuurent version in which dirty is shown
//then using force function it changes dirty = 1 to dirty = 0
//then using steps function to stepUp to the previous stable version
func (m *Migrator) ForceDown() error {
	version, dirty, err := m.migrate.Version()
	if err != nil {
		return err
	}
	if err == nil && dirty {
		if errForce := m.migrate.Force(int(version)); errForce != nil {
			return errForce
		}

		if errSteps := m.migrate.Steps(1); errSteps != nil {
			return errSteps
		}

		return nil
	}

	return fmt.Errorf("no dirty version, invalid use of force command")
}

// Create will create 2 migration files (up and down) with the given title. It uses current timestamp
// for determining the version of the migration.
func (m Migrator) Create(title string) error {
	upFilePath, downFilePath := m.getNewMigrationFilePaths(title)

	upFile, err := os.Create(upFilePath)
	if err != nil {
		return err
	}
	defer upFile.Close()

	downFile, err := os.Create(downFilePath)
	if err != nil {
		return err
	}
	defer downFile.Close()

	return nil
}

func (m Migrator) getNewMigrationFilePaths(title string) (string, string) {
	currentTimeStamp := time.Now().UnixMilli()
	basePath := m.config.MigrationsDir
	upFilePath := fmt.Sprintf("%s/%d_%s.up.sql", basePath, currentTimeStamp, title)
	downFilePath := fmt.Sprintf("%s/%d_%s.down.sql", basePath, currentTimeStamp, title)

	return upFilePath, downFilePath

}
