package repository

import "todo-list-api/models"

type TaskRepository interface {
	GetAllTasks(page int, limit int) (*[]models.Task, error)
	GetRowsCount() (int, error)
	CreateTask(task *models.Task) error
	UpdateTask(task *models.Task) error
	DeleteTask(id int) error
}

type UserRepository interface {
	CreateUser(user *models.User) error
	GetUserByEmail(email string) (*models.User, error)
}
