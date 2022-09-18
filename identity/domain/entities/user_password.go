package entities

type UserPassword struct {
	ID       int64
	UserID   int64
	Password string
	Status   string
}
