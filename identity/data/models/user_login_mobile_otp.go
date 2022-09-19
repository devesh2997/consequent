package models

import "time"

type UserLoginMobileOTP struct {
	ID             int64     `json:"id" gorm:"column:id"`
	VerificationID string    `json:"verification_id" gorm:"column:id"`
	Mobile         int64     `json:"mobile" gorm:"column:mobile"`
	OTP            int       `json:"otp" gorm:"column:otp"`
	Status         string    `json:"status" gorm:"column:status"`
	CreatedAt      time.Time `json:"created_at" gorm:"column:created_at"`
	ExpiryAt       time.Time `json:"expiry_at" gorm:"column:expiry_at"`
	UpdatedAt      time.Time `json:"updated_at" gorm:"column:updated_at"`
}
