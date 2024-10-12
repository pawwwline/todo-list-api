package models

type Task struct {
	Id          int    `json:"id"`
	UserId      int    `json:"user_id"`
	Title       string `json:"title"`
	Description string `json:"decription"`
}
