package models

type Response struct {
	Data  []Task `json:"data"`
	Page  int    `json:"page"`
	Total int    `json:"total"`
	Limit int    `json:"limit"`
}
