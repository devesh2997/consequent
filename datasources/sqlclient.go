package datasources

import (
	"fmt"
	"net/url"

	"github.com/devesh2997/consequent/config"
	"github.com/devesh2997/consequent/sqlclients"
	_ "github.com/go-sql-driver/mysql"
)

// getDroomSQLClient connects to a mysql server and returns ...
func getDroomSQLClients() (*sqlclients.SQLClients, error) {
	inputConfig := config.Config.SQL

	url, err := url.Parse("Asia/Calcutta")
	if err != nil {
		panic(err)
	}

	master := sqlclients.SQLConnConfig{
		User:      inputConfig.Master.User,
		Password:  inputConfig.Master.Password,
		Host:      inputConfig.Master.Host,
		DB:        inputConfig.Master.DB,
		ParseTime: true,
		Loc:       url.Query().Encode(),
	}

	slaves := []sqlclients.SQLConnConfig{}
	for i := 0; i < len(inputConfig.Slaves); i++ {
		slave := sqlclients.SQLConnConfig{
			User:      inputConfig.Slaves[i].User,
			Password:  inputConfig.Slaves[i].Password,
			Host:      inputConfig.Slaves[i].Host,
			DB:        inputConfig.Slaves[i].DB,
			ParseTime: true,
			Loc:       url.Query().Encode(),
		}

		slaves = append(slaves, slave)
	}

	sqlConfig := sqlclients.SQLConfig{
		DriverName: inputConfig.DriverName,
		Master:     master,
		Slaves:     slaves,
	}

	sqlClients, err := sqlclients.Connect(sqlConfig)

	if err != nil {
		return nil, err
	}

	fmt.Println("Connected to SQL!")

	return sqlClients, nil
}
