package models

type Token struct {
	JWT          string `json:"jwt"`
	RefreshToken string `json:"refresh_token"`
}
