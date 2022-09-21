package models

import "github.com/devesh2997/consequent/user/data/constants"

type User struct {
	ID     int64  `json:"id" gorm:"column:id"`
	Mobile string `json:"mobile" gorm:"column:mobile"`
	Email  string `json:"email" gorm:"column:email"`
	Name   string `json:"name" gorm:"column:name"`
	Gender string `json:"gender" gorm:"column:gender"`
}

func (user User) TableName() string {
	return constants.TABLE_NAME_USERS
}
