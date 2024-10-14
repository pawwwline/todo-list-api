package handlers

import (
	"errors"
	"net/http"
	"strconv"

	"todo-list-api/internal/logger"
	"todo-list-api/internal/service/task"
	"todo-list-api/internal/service/utils"
	"todo-list-api/lib/e"
	"todo-list-api/models"
)

type TaskServer struct {
	Service task.TaskService
}

func NewTaskServer(service task.TaskService) *TaskServer {
	return &TaskServer{
		Service: service,
	}
}

func (ts *TaskServer) CreateTaskHandler(w http.ResponseWriter, r *http.Request) {
	var task models.Task
	userId, err := utils.UserIdfromCtx(r.Context())
	logger.Logger.Debug("user id from ctx", "id", userId)
	if err := utils.ParseJson(r, &task); err != nil {
		http.Error(w, "Invalid request format", http.StatusBadRequest)
		return
	}
	task.UserId = userId

	createdTask, err := ts.Service.CreateTask(task)
	if err != nil {
		http.Error(w, "Internal server error", http.StatusInternalServerError)
		return
	}

	err = utils.WriteJSON(w, http.StatusCreated, createdTask)
	if err != nil {
		http.Error(w, "Internal server error", http.StatusInternalServerError)
		return
	}
}

func (ts *TaskServer) GetAllTasksHandler(w http.ResponseWriter, r *http.Request) {
	paginationReq, err := utils.PaginationRequest(w, r)
	if err != nil {
		logger.Logger.Error("failed to parse pagination request", "error", err)
		http.Error(w, "Bad request", http.StatusBadRequest)
		return
	}
	logger.Logger.Debug("pagination req", "limit", paginationReq.Limit, "page", paginationReq.Page, "userId", paginationReq.UserId)
	response, err := ts.Service.GetTasks(paginationReq)
	if err != nil {
		logger.Logger.Error("failed to get tasks", "error", err)
		http.Error(w, "error creating task", http.StatusInternalServerError)
		return
	}
	err = utils.WriteJSON(w, http.StatusOK, response)
	if err != nil {
		logger.Logger.Error("failed to write JSON response", "error", err)
		http.Error(w, "Internal server error", http.StatusInternalServerError)
		return
	}
}

func (ts *TaskServer) DeleteTaskHandler(w http.ResponseWriter, r *http.Request) {
	//getting userId from context
	userId, err := utils.UserIdfromCtx(r.Context())
	id, err := strconv.Atoi(r.PathValue("id"))
	if err != nil {
		http.Error(w, "Invalid request format", http.StatusBadRequest)
		return
	}
	err = ts.Service.DeleteTask(id, userId)
	if err != nil {
		if errors.Is(e.ItemIdNotFound, err) {
			http.Error(w, "task id not found", http.StatusNotFound)
			return
		} else {
			http.Error(w, "error deleting task", http.StatusInternalServerError)
			return
		}
	}
	w.WriteHeader(http.StatusNoContent)
}

func (ts *TaskServer) UpdateTaskHandler(w http.ResponseWriter, r *http.Request) {
	var task models.Task
	userId, err := utils.UserIdfromCtx(r.Context())
	//getting id from url
	id, err := strconv.Atoi(r.PathValue("id"))
	if err != nil {
		http.Error(w, "Error getting UserId", http.StatusInternalServerError)
		return
	}
	if err := utils.ParseJson(r, &task); err != nil {
		http.Error(w, "Invalid request format", http.StatusBadRequest)
		return
	}
	task.UserId = userId
	response, err := ts.Service.UpdateTask(id, task)
	if err != nil {
		if errors.Is(e.ItemIdNotFound, err) {
			http.Error(w, "task id not found", http.StatusNotFound)
			return
		} else {
			http.Error(w, "error deleting task", http.StatusInternalServerError)
			return
		}
	}

	utils.WriteJSON(w, http.StatusOK, response)
}
