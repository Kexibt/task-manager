package repository

import (
	"task_manager/models"
	"time"
)

type Repository interface {
	CreateTask(user *models.User, tsk *models.Task) error
	EditTask(user *models.User, newTsk *models.Task, oldID int) error
	DeleteTask(user *models.User, id int) error

	GetTask(user *models.User, taskID int) (*models.Task, error)
	GetTasks(user *models.User) ([]*models.Task, error)
	GetTasksAfter(user *models.User, date time.Time) ([]*models.Task, error)
}
