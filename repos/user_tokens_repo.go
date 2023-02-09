package repos

import (
	"Portfolio_Nodes/domain"
	"database/sql"
)

type UserTokensRepo struct {
	db *sql.DB
}

func NewUserTokensRepo(db *sql.DB) domain.UserTokensRepository {
	return &UserTokensRepo{db: db}
}

func (r *UserTokensRepo) Insert(token *domain.UserToken) error {
	query := `INSERT INTO user_tokens (user_id, token, expired_at) 
				VALUES ($1, $2, $3) returning id;`
	err := r.db.QueryRow(query, token.UserId, token.Token, token.ExpiredAt).Scan(&token.Id)
	if err != nil {
		return err
	}
	return nil
}

func (r *UserTokensRepo) FindValid(tokenString string) (*domain.UserToken, error) {
	var token = &domain.UserToken{
		Token: tokenString,
	}
	query := `select 
    			id, 
    			user_id, 
    			created_at,
    			expired_at
			from user_tokens
			where token=$1 and now() <= expired_at and is_valid=true and token <> ''
			limit 1`
	err := r.db.QueryRow(query, tokenString).Scan(&token.Id,
		&token.UserId,
		&token.CreatedAt,
		&token.ExpiredAt)
	switch err {
	case nil:
		return token, nil
	case sql.ErrNoRows:
		return nil, nil
	default:
		return nil, err
	}
}

func (r *UserTokensRepo) Revoke(s string) error {
	query := `UPDATE user_tokens 
				SET is_valid=false
				WHERE token=$1`
	_, err := r.db.Exec(query, s)
	if err != nil {
		return err
	}
	return nil
}

func (r *UserTokensRepo) UpdateTime(s string) error {
	query := `UPDATE user_tokens 
				SET updated_at=now()
				WHERE token=$1`
	_, err := r.db.Exec(query, s)
	if err != nil {
		return err
	}
	return nil
}
