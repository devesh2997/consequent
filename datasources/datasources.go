package datasources

import (
	"sync"

	"github.com/devesh2997/consequent/sqlclients"
)

var dataSources DataSources
var dataSourcesErrors []error

// DataSources is...
type DataSources struct {
	SQLClients *sqlclients.SQLClients
}

var dsOnce sync.Once

func Get() (*DataSources, error) {
	dsOnce.Do(func() {
		sqlClients, err := getDroomSQLClients()
		if err == nil {
			dataSources.SQLClients = sqlClients
		} else {
			dataSourcesErrors = append(dataSourcesErrors, err)
		}
	})

	if len(dataSourcesErrors) > 0 {
		return &dataSources, dataSourcesErrors[0]
	}

	return &dataSources, nil
}
