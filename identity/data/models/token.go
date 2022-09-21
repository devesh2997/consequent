package models

import (
	"time"

	"github.com/devesh2997/consequent/identity/data/constants"
)

type JWT struct {
	Token    string    `json:"token"`
	ExpiryAt time.Time `json:"expiry_at"`
}

type RefreshToken struct {
	ID        int64     `json:"-" gorm:"column:id"`
	Token     string    `json:"token" gorm:"column:token"`
	Status    string    `json:"-" gorm:"column:status"`
	CreatedAt time.Time `json:"-" gorm:"column:created_at"`
	ExpiryAt  time.Time `json:"expiry_at" gorm:"column:created_at"`
	UpdatedAt time.Time `json:"-" gorm:"column:updated_at"`
}

func (RefreshToken) TableName() string {
	return constants.TABLE_NAME_REFRESH_TOKENS
}

type Token struct {
	JWT          JWT          `json:"jwt"`
	RefreshToken RefreshToken `json:"refresh_token"`
}
