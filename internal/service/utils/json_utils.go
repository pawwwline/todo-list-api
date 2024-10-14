package utils

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"strconv"
	"todo-list-api/internal/server/middleware"
	"todo-list-api/models"
)

func ParseJson(r *http.Request, v any) error {
	if r.Body == nil {
		return fmt.Errorf("missing request body")
	}
	defer r.Body.Close()
	if err := json.NewDecoder(r.Body).Decode(v); err != nil {
		return fmt.Errorf("failed to decode json: %w", err)
	}
	return nil
}

func WriteJSON(w http.ResponseWriter, status int, v any) error {
	w.Header().Add("Content-Type", "application/json")
	w.WriteHeader(status)
	return json.NewEncoder(w).Encode(v)
}

func PaginationRequest(w http.ResponseWriter, r *http.Request) (models.PaginationRequest, error) {
	userId, err := UserIdfromCtx(r.Context())
	if err != nil {
		http.Error(w, "Unable to retrieve user ID", http.StatusUnauthorized)
		return models.PaginationRequest{}, err
	}
	pageStr := r.URL.Query().Get("page")
	limitStr := r.URL.Query().Get("limit")
	limit, err := strconv.Atoi(limitStr)
	if err != nil {
		http.Error(w, "Invalid limit format", http.StatusBadRequest)
		return models.PaginationRequest{}, err
	}
	page, err := strconv.Atoi(pageStr)
	if err != nil {
		http.Error(w, "Invalid request format", http.StatusBadRequest)
		return models.PaginationRequest{}, err
	}
	pagintaionReq := models.PaginationRequest{
		Limit:  limit,
		Page:   page,
		UserId: userId,
	}
	return pagintaionReq, nil
}

func UserIdfromCtx(ctx context.Context) (int, error) {
	userId, ok := ctx.Value(middleware.UserIdKey).(int)
	if !ok {
		return 0, errors.New("user ID not found in context")
	}
	return userId, nil
}
