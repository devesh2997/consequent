package sqlclients

import (
	"database/sql"

	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"gorm.io/plugin/dbresolver"
)

func getGORMClientDatabase(driverName string, sqlDBs []*sql.DB) (*gorm.DB, error) {
	cfg := &gorm.Config{
		SkipDefaultTransaction: true,
	}
	db, err := gorm.Open(mysql.New(mysql.Config{DriverName: driverName, Conn: sqlDBs[0]}), cfg)
	if err != nil {
		return nil, err
	}

	var replicas []gorm.Dialector
	for i := 1; i < len(sqlDBs); i++ {
		replica := mysql.New(mysql.Config{DriverName: driverName, Conn: sqlDBs[i]})
		replicas = append(replicas, replica)
	}

	db.Use(dbresolver.Register(dbresolver.Config{
		Replicas: replicas,
		Policy:   dbresolver.RandomPolicy{},
	}))

	return db, nil
}
