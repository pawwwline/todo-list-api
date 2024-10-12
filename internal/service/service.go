package service

import (
	"todo-list-api/models"
)

type TaskServiceInterface interface {
	GetTasks(req models.PaginationRequest) (*models.Response, error)
	CreateTask(task models.Task) (*models.Task, error)
	UpdateTask(id int, task models.Task) (*models.Task, error)
	DeleteTask(id int) error
}

type UserServiceInterface interface {
	SignIn(user models.User) (string, error)
	LogIn(user models.User) error
}

type Service struct {
	TaskService TaskServiceInterface
	UserService UserServiceInterface
}
