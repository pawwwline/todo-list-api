package handlers

import (
	"net/http"
	"strconv"

	"todo-list-api/internal/service/task"
	"todo-list-api/internal/service/utils"
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
	paginationReq := utils.PaginationRequest(w, r)
	response, err := ts.Service.GetTasks(paginationReq)
	err = utils.WriteJSON(w, http.StatusOK, response)
	if err != nil {
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
	ts.Service.DeleteTask(id, userId)
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
	utils.WriteJSON(w, http.StatusOK, response)
}
