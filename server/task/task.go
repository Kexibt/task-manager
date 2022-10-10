package task

import (
	"errors"
	"task_manager/models"
	"task_manager/task/repository"
	"task_manager/task/repository/localstore"
	"task_manager/task/repository/psql"
	"time"

	"github.com/jackc/pgx/v4/pgxpool"
)

type TasksRepository struct {
	Localstore repository.Repository
	Database   repository.Repository
}

func NewTasksRepository(conn *pgxpool.Pool) *TasksRepository {
	task := &TasksRepository{
		Localstore: localstore.NewLocalStorage(),
		Database:   psql.NewTaskRepository(conn),
	}

	return task
}

func (t *TasksRepository) CreateTask(user *models.User, tsk *models.Task) error {
	err := t.Database.CreateTask(user, tsk)
	if err != nil {
		return err
	}

	err = t.Localstore.CreateTask(user, tsk)
	return err
}

func (t *TasksRepository) EditTask(user *models.User, newTsk *models.Task, oldID int) error {
	original, err := t.GetTask(user, oldID)
	if err != nil {
		return err
	}

	changeOrigin(newTsk, original)
	err = t.Database.EditTask(user, original, oldID)
	if err != nil {
		return err
	}

	err = t.Localstore.EditTask(user, original, oldID)
	if errors.Is(err, localstore.ErrTaskNotFound) {
		err = t.Localstore.CreateTask(user, original)
		if err != nil {
			return err
		}
	}

	return err
}

func changeOrigin(left, right *models.Task) {
	if left.Title != "" {
		right.Title = left.Title
	}
	if left.Description != "" {
		right.Description = left.Description
	}
	if left.Status != "" {
		right.Status = left.Status
	}
	if left.UserID != "" {
		right.UserID = left.UserID
	}
}

func (t *TasksRepository) DeleteTask(user *models.User, id int) error {
	err := t.Database.DeleteTask(user, id)
	if err != nil {
		return err
	}

	err = t.Localstore.DeleteTask(user, id)
	return err
}

func (t *TasksRepository) GetTask(user *models.User, taskID int) (*models.Task, error) {
	tsk, err := t.Localstore.GetTask(user, taskID)
	if errors.Is(err, localstore.ErrTaskNotFound) {
		tsk, err = t.Database.GetTask(user, taskID)
		if err != nil {
			return nil, err
		}
	}
	return tsk, err
}

func (t *TasksRepository) GetTasks(user *models.User) ([]*models.Task, error) {
	tsks0, _ := t.Localstore.GetTasks(user)
	// if err != nil {
	// } // always nil, might be changed

	tsks1, _ := t.Database.GetTasks(user)
	// if err != nil {
	// } // always nil, might be changed

	if len(tsks0) == len(tsks1) {
		return tsks0, nil
	}

	if len(tsks1) > len(tsks0) {
		for i := 0; i < len(tsks1); i++ {
			t.Localstore.CreateTask(user, tsks1[i])
		}
	}
	return tsks1, nil
}

func (t *TasksRepository) GetTasksAfter(user *models.User, date time.Time) ([]*models.Task, error) {
	tsks0, _ := t.Localstore.GetTasksAfter(user, date)
	// if err != nil {
	// } // always nil, might be changed

	tsks1, _ := t.Database.GetTasksAfter(user, date)
	// if err != nil {
	// } // always nil, might be changed

	if len(tsks0) == len(tsks1) {
		return tsks0, nil
	}

	if len(tsks1) > len(tsks0) {
		for i := 0; i < len(tsks1); i++ {
			t.Localstore.CreateTask(user, tsks1[i])
		}
	}
	return tsks1, nil
}
