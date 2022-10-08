package http

import (
	"task_manager/models"
	"task_manager/task/repository"
)

type Task struct {
	ID          int    `json:"id"`
	Title       string `json:"title"`
	Description string `json:"description"`
	Status      string `json:"status"`
	UserID      string `json:"owner"`
	UpdateDate  string `json:"last_update"`
}

type Handler struct {
	repository repository.Repository
}

func NewTaskHandler(repository repository.Repository) *Handler {
	return &Handler{
		repository: repository,
	}
}

type createInput struct {
	Title       string `json:"title"`
	Description string `json:"description"`
	Status      string `json:"status"`
}

func (h *Handler) Create() {
	// var inp createInput
}

type getResponse struct {
	Bookmarks []*Task `json:"bookmarks"`
}

func (h *Handler) Get() {

}

type deleteInput struct {
	ID string `json:"id"`
}

func (h *Handler) Delete() {
}

func convertToTasks(bs []*models.Task) []*Task {
	out := make([]*Task, len(bs))

	for i, b := range bs {
		out[i] = convertToTask(b)
	}

	return out
}

func convertToTask(b *models.Task) *Task {
	return &Task{
		// todo
	}
}
