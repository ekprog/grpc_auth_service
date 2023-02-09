package repos

import (
	"Portfolio_Nodes/domain"
	"database/sql"
)

type UsersRepo struct {
	db *sql.DB
}

func NewUsersRepo(db *sql.DB) domain.UsersRepository {
	return &UsersRepo{db: db}
}

func (r *UsersRepo) Insert(user *domain.User) error {
	query := `INSERT INTO users (username, pwd_hash) VALUES ($1, $2) returning id;`

	err := r.db.QueryRow(query, user.Username, user.PwdHash).Scan(&user.Id)
	if err != nil {
		return err
	}
	return nil
}

func (r *UsersRepo) UpdateById(user domain.User) error {
	query := `UPDATE users 
				SET username=$2, pwd_hash=$3, updated_at=now()
				WHERE id=$1`
	_, err := r.db.Exec(query, user.Id, user.Username, user.PwdHash)
	if err != nil {
		return err
	}
	return nil
}

func (r *UsersRepo) FindById(id int64) (*domain.User, error) {
	var user = &domain.User{
		Id: id,
	}
	query := `select 
    			username, 
    			pwd_hash,
    			created_at, 
    			updated_at 
			from users
			where id=$1
			limit 1`
	err := r.db.QueryRow(query, id).Scan(&user.Username, &user.PwdHash, &user.CreatedAt, &user.UpdatedAt)
	switch err {
	case nil:
		return user, nil
	case sql.ErrNoRows:
		return nil, nil
	default:
		return nil, err
	}
}

func (r *UsersRepo) FindByUsername(username string) (*domain.User, error) {
	var user = &domain.User{
		Username: username,
	}
	query := `select 
    			id, 
    			pwd_hash,
    			created_at, 
    			updated_at 
			from users
			where username=$1
			limit 1`
	err := r.db.QueryRow(query, username).Scan(&user.Id, &user.PwdHash, &user.CreatedAt, &user.UpdatedAt)
	switch err {
	case nil:
		return user, nil
	case sql.ErrNoRows:
		return nil, nil
	default:
		return nil, err
	}
}
