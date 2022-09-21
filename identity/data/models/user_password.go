package models

import "github.com/devesh2997/consequent/identity/data/constants"

type UserPassword struct {
	ID       int64  `json:"id" gorm:"column:id"`
	UserID   int64  `json:"user_id" gorm:"column:user_id"`
	Password string `json:"password" gorm:"column:password"`
	Status   string `json:"status" gorm:"column:status"`
}

func (UserPassword) TableName() string {
	return constants.TABLE_NAME_USER_PASSWORDS
}
