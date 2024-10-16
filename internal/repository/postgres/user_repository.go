package postgres

import (
	"database/sql"
	"errors"
	"todo-list-api/lib/e"
	"todo-list-api/models"

	"github.com/lib/pq"
)

type PostgresUserRepo struct {
	DB *sql.DB
}

func NewPostgresUserRepo(db *sql.DB) *PostgresUserRepo {
	return &PostgresUserRepo{
		DB: db,
	}
}

func (repo *PostgresUserRepo) CreateUser(user *models.User) (int64, error) {
	var user_id int64
	err := repo.DB.QueryRow("INSERT INTO users (name, email, password) VALUES ($1, $2, $3) RETURNING id;", user.Name, user.Email, user.Password).Scan(&user_id)
	if err != nil {
		//return specific error if email already exist
		if pgErr, ok := err.(*pq.Error); ok {
			if pgErr.Code == "23505" {
				return 0, e.UniqueViolationErr
			}
		}
		return 0, err
	}

	return user_id, nil
}

// returns user model, if user not found returns nil
func (repo *PostgresUserRepo) GetUserByEmail(email string) (*models.User, error) {
	var user models.User
	row := repo.DB.QueryRow("SELECT id, name, email, password FROM users WHERE email=$1", email)
	err := row.Scan(&user.Id, &user.Name, &user.Email, &user.Password)
	if errors.Is(err, sql.ErrNoRows) {
		return nil, nil
	} else if err != nil {
		return nil, err
	}

	return &user, nil
}
