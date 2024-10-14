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

func (repo *PostgresTaskRepo) GetAllTasks(req models.PaginationRequest) ([]models.Task, error) {
	var tasks []models.Task
	offset := (req.Limit * req.Page) - req.Limit
	rows, err := repo.DB.Query("SELECT * FROM tasks WHERE user_id = $1 ORDER BY id LIMIT $2 OFFSET $3;", req.UserId, req.Limit, offset)
	if err != nil {
		return nil, fmt.Errorf("query get all tasks failed: %v", err)
	}
	defer rows.Close()
	for rows.Next() {
		var task models.Task
		if err := rows.Scan(&task.Id, &task.UserId, &task.Title, &task.Description); err != nil {
			return nil, fmt.Errorf("error scannig rows in task struct: %v", err)
		}
		tasks = append(tasks, task)
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}
	return tasks, err
}

func (repo *PostgresTaskRepo) GetRowsCount(userId int) (int, error) {
	//getting total for pagination response
	var total int
	rows, err := repo.DB.Query("SELECT COUNT(*) FROM tasks WHERE user_id = $1;", userId)
	if err != nil {
		return 0, fmt.Errorf("error counting rows task: %v", err)
	}
	defer rows.Close()
	if rows.Next() {
		if err := rows.Scan(&total); err != nil {
			return 0, fmt.Errorf("error scannig rows in total: %v", err)
		}
	}
	return total, nil
}

func (repo *PostgresTaskRepo) CreateTask(task *models.Task) (int, error) {
	var id int
	res := repo.DB.QueryRow("INSERT INTO tasks(user_id, title, description) VALUES($1, $2, $3) RETURNING id", task.UserId, task.Title, task.Description)
	//return id for response
	err := res.Scan(&id)
	if err != nil {
		return 0, fmt.Errorf("error creating task: %v", err)
	}

	return id, nil
}

func (repo *PostgresTaskRepo) UpdateTask(task *models.Task) error {
	res, err := repo.DB.Exec("UPDATE tasks SET title=$1, description=$2 WHERE user_id=$3 AND id=$4", task.Title, task.Description, task.UserId, task.Id)
	count, err := res.RowsAffected()
	if err != nil {
		return fmt.Errorf("error updating task: %v", err)
	}
	if count == 0 {
		return e.ItemIdNotFound
	}
	return nil
}

func (repo *PostgresTaskRepo) DeleteTask(id int, userId int) error {
	res, err := repo.DB.Exec("DELETE from tasks WHERE user_id = $1 AND id=$2", userId, id)
	count, err := res.RowsAffected()
	if err != nil {
		return fmt.Errorf("error deleting task: %v", err)
	}
	if count == 0 {
		return e.ItemIdNotFound
	}

	return nil
}
