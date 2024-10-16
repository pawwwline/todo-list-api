package postgres

import (
	"database/sql"
	r "todo-list-api/internal/repository"
)

func NewPostgresRepo(db *sql.DB) *r.Repository {
	return &r.Repository{
		Task: NewPostgresTaskRepo(db),
		User: NewPostgresUserRepo(db),
	}
}
