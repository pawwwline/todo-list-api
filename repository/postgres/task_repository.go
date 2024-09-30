package postgres

import (
	"database/sql"
	"fmt"
	"todo-list-api/lib/e"
	"todo-list-api/models"
)

type PostgresTaskRepo struct {
	DB *sql.DB
}

func NewPostgresTaskRepo(db *sql.DB) *PostgresTaskRepo {
	return &PostgresTaskRepo{
		DB: db,
	}
}

func (repo *PostgresTaskRepo) GetAllTasks(req models.PaginationRequest) (*[]models.Task, error) {
	var tasks []models.Task
	offset := (req.Limit * req.Page) - req.Limit
	rows, err := repo.DB.Query("SELECT FROM * tasks ORDER BY ID LIMIT $1 OFFSET $2", req.Limit, offset)
	if err != nil {
		return nil, e.WrapError("query get all tasks failed", err)
	}
	defer rows.Close()
	for rows.Next() {
		var task models.Task
		if err := rows.Scan(&task.Id, &task.UserId, &task.Title, &task.Description); err != nil {
			return nil, e.WrapError("failed to scan row into task struct", err)
		}
		tasks = append(tasks, task)
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}
	return &tasks, err
}

func (repo *PostgresTaskRepo) GetRowsCount() (int, error) {
	var total int
	rows, err := repo.DB.Query("SELECT FROM COUNT(*) tasks")
	if err != nil {
		return -1, e.WrapError("failed to count rows", err)
	}
	defer rows.Close()
	err = rows.Scan(&total)
	if err != nil {
		return -1, err
	}
	return total, nil
}

func (repo *PostgresTaskRepo) CreateTask(task *models.Task) (int, error) {
	var id int
	res := repo.DB.QueryRow("INSERT INTO tasks(user_id, title, description) VALUES($1, $2, $3) RETURNING id", task.UserId, task.Title, task.Description)
	//return id for response
	err := res.Scan(&id)
	if err != nil {
		return -1, e.WrapError("failed to create task", err)
	}

	return id, nil
}

func (repo *PostgresTaskRepo) UpdateTask(task *models.Task) error {
	res, err := repo.DB.Exec("UPDATE tasks SET title=$1, description=$2 WHERE id=$3", task.Title, task.Description, task.Id)
	count, err := res.RowsAffected()
	if err != nil {
		return e.WrapError("failed update task", err)
	}
	if count == 0 {
		return e.WrapError(fmt.Sprintf("no rows affected, id: %d not found", task.Id), err)
	}
	return nil
}

func (repo *PostgresTaskRepo) DeleteTask(id int) error {
	res, err := repo.DB.Exec("DELETE from tasks WHERE id=$1", id)
	count, err := res.RowsAffected()
	if err != nil {
		return e.WrapError("failed update task", err)
	}
	if count == 0 {
		return e.WrapError(fmt.Sprintf("no rows affected, id: %d not found", id), err)
	}

	return nil
}
