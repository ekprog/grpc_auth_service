package repos

import (
	"auth_service/app"
	"auth_service/domain"
	"database/sql"
)

type UserTokensRepo struct {
	log app.Logger
	db  *sql.DB
}

func NewUserTokensRepo(log app.Logger, db *sql.DB) domain.UserTokensRepository {
	return &UserTokensRepo{
		log: log,
		db:  db,
	}
}

func (r *UserTokensRepo) Insert(token *domain.UserToken) error {
	query := `INSERT INTO user_tokens (
                         user_id, 
                         pair_uuid,
                         access_token, 
                         refresh_token, 
                         access_token_expired_at, 
                         refresh_token_expired_at) 
				VALUES ($1, $2, $3, $4, $5, $6) returning id;`
	err := r.db.QueryRow(query,
		token.UserId,
		token.PairUUID,
		token.AccessToken,
		token.RefreshToken,
		token.AccessTokenExpiredAt,
		token.RefreshTokenExpiredAt).Scan(&token.Id)
	if err != nil {
		return err
	}
	return nil
}

func (r *UserTokensRepo) FindValidPair(pairUUID string) (*domain.UserToken, error) {
	var token = &domain.UserToken{
		PairUUID: pairUUID,
	}
	query := `select 
    			id, 
    			user_id, 
    			access_token,
    			refresh_token,
    			access_token_expired_at,
    			refresh_token_expired_at,
    			updated_at,
    			created_at
			from user_tokens
			where pair_uuid=$1 and is_valid=true and access_token <> ''
			limit 1`
	err := r.db.QueryRow(query, pairUUID).Scan(&token.Id,
		&token.UserId,
		&token.AccessToken,
		&token.RefreshToken,
		&token.AccessTokenExpiredAt,
		&token.RefreshTokenExpiredAt,
		&token.UpdatedAt,
		&token.CreatedAt,
	)
	switch err {
	case nil:
		return token, nil
	case sql.ErrNoRows:
		return nil, nil
	default:
		return nil, err
	}
}

func (r *UserTokensRepo) RevokePair(pairUUID string) error {
	query := `UPDATE user_tokens
				SET is_valid=false
				WHERE pair_uuid=$1`
	_, err := r.db.Exec(query, pairUUID)
	if err != nil {
		return err
	}
	return nil
}

func (r *UserTokensRepo) UpdateTime(pairUUID string) error {
	query := `UPDATE user_tokens 
				SET updated_at=now()
				WHERE pair_uuid=$1`
	_, err := r.db.Exec(query, pairUUID)
	if err != nil {
		return err
	}
	return nil
}
