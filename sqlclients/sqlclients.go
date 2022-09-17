package sqlclients

import (
	"database/sql"
	"fmt"

	"gorm.io/gorm"
)

type SQLConnConfig struct {
	User      string
	Password  string
	Host      string
	DB        string
	ParseTime bool
	Loc       string
}

type SQLConfig struct {
	DriverName string
	Master     SQLConnConfig
	Slaves     []SQLConnConfig
}

type SQLClients struct {
	sqlxDB *SQLXClusterDB
	gormDB *gorm.DB
}

func (clients SQLClients) GetSQLXClusterDB() *SQLXClusterDB {
	return clients.sqlxDB
}

func (clients SQLClients) GetGormDB() *gorm.DB {
	return clients.gormDB
}

func Connect(sqlConfig SQLConfig) (*SQLClients, error) {
	masterFormattedDataSourceName := getFormattedDataSourceName(sqlConfig.Master)

	var slaveFormattedDataSourceNames []string
	for _, slaveConnConfig := range sqlConfig.Slaves {
		formatted := getFormattedDataSourceName(slaveConnConfig)
		slaveFormattedDataSourceNames = append(slaveFormattedDataSourceNames, formatted)
	}

	dataSourceNames := append([]string{masterFormattedDataSourceName}, slaveFormattedDataSourceNames...)

	dbs := make([]*sql.DB, len(dataSourceNames))

	err := scatter(len(dbs), func(i int) (err error) {
		sqlDB, e := sql.Open(sqlConfig.DriverName, dataSourceNames[i])
		if e != nil {
			return e
		}
		if err := sqlDB.Ping(); err != nil {
			return err
		}

		dbs[i] = sqlDB
		return err
	})
	if err != nil {
		return nil, err
	}

	sqlxDB, err := getSQLXClientDatabase(sqlConfig.DriverName, dbs)
	if err != nil {
		return nil, err
	}

	gormDB, err := getGORMClientDatabase(sqlConfig.DriverName, dbs)
	if err != nil {
		return nil, err
	}

	return &SQLClients{sqlxDB: sqlxDB, gormDB: gormDB}, nil
}

func getFormattedDataSourceName(sqlConnConfig SQLConnConfig) string {
	dataSourceNameFormat := "%s:%s@tcp(%s)/%s"
	if sqlConnConfig.ParseTime {
		dataSourceNameFormat += "?parseTime=true&loc=" + sqlConnConfig.Loc
	}
	user := sqlConnConfig.User
	password := sqlConnConfig.Password
	host := sqlConnConfig.Host
	db := sqlConnConfig.DB

	dataSourceName := fmt.Sprintf(dataSourceNameFormat, user, password, host, db)

	return dataSourceName
}
