package sqlclients

import (
	"errors"

	"github.com/go-sql-driver/mysql"
)

func IsDuplicateError(err error) bool {
	var sqlError *mysql.MySQLError
	if ok := errors.As(err, &sqlError); ok {
		if sqlError.Number == 1062 {
			return true
		}
	}

	return false
}
