package entities

import "time"

type JWT struct {
	Token    string
	ExpiryAt time.Time
}

type RefreshToken struct {
	ID        int64
	Token     string
	Status    string
	CreatedAt time.Time
	ExpiryAt  time.Time
	UpdatedAt time.Time
}

type Token struct {
	JWT          JWT
	RefreshToken RefreshToken
}
