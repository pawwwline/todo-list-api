package repository

import "todo-list-api/models"

type TaskRepository interface {
	GetAllTasks(req models.PaginationRequest) (*[]models.Task, error)
	GetRowsCount(userId int) (int, error)
	CreateTask(task *models.Task) (int, error)
	UpdateTask(task *models.Task) error
	DeleteTask(id, userId int) error
}

type UserRepository interface {
	CreateUser(user *models.User) (int64, error)
	GetUserByEmail(email string) (*models.User, error)
}

type Repository struct {
	Task TaskRepository
	User UserRepository
}
