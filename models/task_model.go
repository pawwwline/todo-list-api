package models

type Task struct {
	Id          int    `json:"id"`
	UserId      int    `json:"userId"`
	Title       string `json:"title"`
	Description string `json:"decription"`
}
