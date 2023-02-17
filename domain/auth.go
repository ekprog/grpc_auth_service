package domain

import (
	"time"
)

// Models

type User struct {
	Id        int64     `json:"id"`
	Username  string    `json:"username"`
	PwdHash   string    `json:"-"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

type UserToken struct {
	Id                    int64     `json:"id"`
	PairUUID              string    `json:"pair_uuid"`
	UserId                int64     `json:"user_id"`
	AccessToken           string    `json:"access_token"`
	RefreshToken          string    `json:"refresh_token"`
	IsValid               bool      `json:"is_valid_access_token"`
	AccessTokenExpiredAt  time.Time `json:"access_token_expired_at"`
	RefreshTokenExpiredAt time.Time `json:"refresh_token_expired_at"`

	UpdatedAt time.Time `json:"updated_at"`
	CreatedAt time.Time `json:"created_at"`
}

type UsersRepository interface {
	Insert(*User) error
	UpdateById(User) error
	FindById(id int64) (*User, error)
	FindByUsername(username string) (*User, error)
}

type UserTokensRepository interface {
	Insert(*UserToken) error
	FindValidPair(pairUUID string) (*UserToken, error)
	RevokePair(pairUUID string) error
	UpdateTime(pairUUID string) error
}

type AuthInteractor interface {
	Register(username, password string) (RegisterResponse, error)
	Login(username, password string) (LoginResponse, error)
	Revoke(token string) (RevokeResponse, error)
	RefreshToken(token string) (RefreshTokenResponse, error)
	Extract(token string) (ExtractResponse, error)
}

// Responses

type RegisterResponse struct {
	StatusCode string
}

type LoginResponse struct {
	StatusCode string
	UserToken  *UserToken
}

type RevokeResponse struct {
	StatusCode string
}

type RefreshTokenResponse struct {
	StatusCode string
	UserToken  *UserToken
}

type ExtractResponse struct {
	StatusCode string
	User       *User
}
