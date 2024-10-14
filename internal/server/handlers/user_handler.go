package handlers

import (
	"errors"
	"net/http"
	"todo-list-api/internal/logger"
	"todo-list-api/internal/service/user"
	"todo-list-api/internal/service/utils"
	"todo-list-api/lib/e"
	"todo-list-api/models"
)

type UserServer struct {
	Service user.UserService
}

func NewUserServer(service user.UserService) *UserServer {
	return &UserServer{
		Service: service,
	}
}

func (us *UserServer) LoginHandler(w http.ResponseWriter, r *http.Request) {
	var user models.User
	err := utils.ParseJson(r, &user)
	if err != nil {
		http.Error(w, "Invalid request format", http.StatusBadRequest)
		logger.Logger.Debug("user", "user_json", user)
		return
	}
	token, err := us.Service.LogIn(user)
	if err != nil {
		if errors.Is(err, e.InvalidCredentialsErr) {
			http.Error(w, "Invalid password or email", http.StatusUnauthorized)
			return
		}
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return

	}

	utils.WriteJSON(w, http.StatusOK, map[string]string{"token": token})
}

func (us *UserServer) SignupHandler(w http.ResponseWriter, r *http.Request) {
	var user models.User
	err := utils.ParseJson(r, &user)
	if err != nil {
		http.Error(w, "Invalid request format", http.StatusBadRequest)
	}
	token, err := us.Service.SignUp(user)
	if err != nil {
		if errors.Is(err, e.UniqueViolationErr) {
			http.Error(w, "user with this email already exists", http.StatusConflict)
			logger.Logger.Error("error creating token", "error", err)
			return
		} else {
			http.Error(w, "error creating token", http.StatusInternalServerError)
			logger.Logger.Error("error creating token", "error", err)
			return
		}
	}
	logger.Logger.Debug("token string", "token", token)
	utils.WriteJSON(w, http.StatusOK, map[string]string{"token": token})

}
