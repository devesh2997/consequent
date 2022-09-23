package services

import "github.com/devesh2997/consequent/errorx"

var errUserNotFound = func() error {
	return errorx.NewNotFoundError(-1, "user", "sql")
}
