package domain

import "time"

type User struct {
	Id        int64     `json:"id"`
	Username  string    `json:"username"`
	PwdHash   string    `json:"-"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

type UserToken struct {
	Id        int64     `json:"id"`
	UserId    int64     `json:"user_id"`
	Token     string    `json:"token"`
	IsValid   bool      `json:"is_valid"`
	UpdatedAt time.Time `json:"updated_at"`
	CreatedAt time.Time `json:"created_at"`
	ExpiredAt time.Time `json:"expired_at"`
}

type UsersRepository interface {
	Insert(*User) error
	UpdateById(User) error
	FindById(id int64) (*User, error)
	FindByUsername(username string) (*User, error)
}

type UserTokensRepository interface {
	Insert(*UserToken) error
	FindValid(string) (*UserToken, error)
	Revoke(string) error
	UpdateTime(string) error
}

type AuthInteractor interface {
	Register(username, password string) error
	Login(username, password string) (*UserToken, error)
	ValidateAndExtract(token string) (*User, error)
}
