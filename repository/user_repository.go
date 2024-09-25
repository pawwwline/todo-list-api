package repository

import (
	"database/sql"
	"errors"
	"todo-list-api/models"
)

type PostgresUserRepo struct {
	DB *sql.DB
}

func NewPostgresUserRepo(db *sql.DB) *PostgresUserRepo {
	return &PostgresUserRepo{
		DB: db,
	}
}

func (repo *PostgresUserRepo) CreateUser(user *models.User) error {
	_, err := repo.DB.Exec("INSERT INTO users (name, username, email, password) VALUES ($1, $2, $3, $4)", user.Name, user.Username, user.Email, user.Password)
	if err != nil {
		return err
	}
	return nil
}

func (repo *PostgresUserRepo) GetUserByEmail(email string) (*models.User, error) {
	var user models.User
	row := repo.DB.QueryRow("SELECT id, name, username, email, password FROM users WHERE email=$1", email)
	err := row.Scan(&user.Id, &user.Name, &user.Username, &user.Email, &user.Password)
	if errors.Is(err, sql.ErrNoRows) {
		return nil, nil
	} else if err != nil {
		return nil, err
	}
	return &user, nil
}
