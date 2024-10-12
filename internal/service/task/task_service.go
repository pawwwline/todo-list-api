package task

import (
	"fmt"
	"math"
	"todo-list-api/internal/repository"
	"todo-list-api/models"
)

type TaskService struct {
	taskRepo repository.Repository
}

func NewTaskService(taskRepo *repository.Repository) *TaskService {
	return &TaskService{
		taskRepo: *taskRepo,
	}
}

func (s *TaskService) CreateTask(task models.Task) (*models.Task, error) {
	//adding id to response
	id, err := s.taskRepo.Task.CreateTask(&task)
	if err != nil {
		return nil, err
	}
	createdTask := models.Task{
		Id:          id,
		Title:       task.Title,
		Description: task.Description,
	}
	return &createdTask, nil
}

func (s *TaskService) GetTasks(req models.PaginationRequest) (*models.Response, error) {
	tasks, err := s.taskRepo.Task.GetAllTasks(req)
	if err != nil {
		return nil, err
	}
	//count total pages for pagination response
	totalRows, err := s.taskRepo.Task.GetRowsCount(req.UserId)
	if err != nil {
		return nil, fmt.Errorf("error getting pagination response: %v", err)
	}
	total := math.Ceil(float64(totalRows) / float64(req.Limit))
	totalInt := int(total)

	response := models.Response{
		Data:  *tasks,
		Page:  req.Page,
		Limit: req.Limit,
		Total: totalInt,
	}
	return &response, nil
}

func (s *TaskService) UpdateTask(id int, task models.Task) (*models.Task, error) {
	task.Id = id
	err := s.taskRepo.Task.UpdateTask(&task)
	if err != nil {
		return nil, err
	}
	return &task, nil
}

func (s *TaskService) DeleteTask(id, userId int) error {
	return s.taskRepo.Task.DeleteTask(id, userId)
}
