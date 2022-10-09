package localstore

import (
	"sync"
	"task_manager/models"
	"time"
)

type TaskLocalStorage struct {
	tasks map[int]*models.Task
	mutex *sync.RWMutex
}

func NewLocalStorage() *TaskLocalStorage {
	return &TaskLocalStorage{
		tasks: make(map[int]*models.Task),
		mutex: new(sync.RWMutex),
	}
}

func (t *TaskLocalStorage) CreateTask(user *models.User, tsk *models.Task) error {
	t.mutex.Lock()
	defer t.mutex.Unlock()

	tsk.UserID = user.Login
	if _, exists := t.tasks[tsk.ID]; exists {
		return ErrTaskAlreadyExists
	}

	if tsk.Status != "todo" &&
		tsk.Status != "in progress" &&
		tsk.Status != "done" &&
		tsk.Status != "" {
		return ErrInvalidStatus
	}

	t.tasks[tsk.ID] = tsk
	return nil
}

func (t *TaskLocalStorage) EditTask(user *models.User, newTsk *models.Task, oldID int) error {
	t.mutex.Lock()
	defer t.mutex.Unlock()

	if _, exists := t.tasks[oldID]; !exists {
		return ErrTaskNotFound
	}

	// в задании явно не требуется ограничивать доступ,
	// но я ограничу, потому что должен быть смысл в авториазации ¯\_( ͡° ͜ʖ ͡°)_/¯
	if t.tasks[oldID].UserID != user.Login {
		return ErrDontHavePermission
	}
	if newTsk.Status != "todo" &&
		newTsk.Status != "in progress" &&
		newTsk.Status != "done" &&
		newTsk.Status != "" {
		return ErrInvalidStatus
	}

	if newTsk.Title != "" {
		t.tasks[oldID].Title = newTsk.Title
	}
	if newTsk.Description != "" {
		t.tasks[oldID].Description = newTsk.Description
	}
	if newTsk.Status != "" {
		t.tasks[oldID].Status = newTsk.Status
	}
	if newTsk.UserID != "" {
		t.tasks[oldID].UserID = newTsk.UserID
	}
	t.tasks[oldID].UpdateDate = newTsk.UpdateDate
	return nil

}

func (t *TaskLocalStorage) DeleteTask(user *models.User, id int) error {
	t.mutex.Lock()
	defer t.mutex.Unlock()

	if _, exists := t.tasks[id]; !exists {
		return ErrTaskNotFound
	}
	if t.tasks[id].UserID != user.Login {
		return ErrDontHavePermission
	}

	delete(t.tasks, id)
	return nil
}

func (t *TaskLocalStorage) GetTask(user *models.User, taskID int) (*models.Task, error) {
	t.mutex.RLock()
	defer t.mutex.RUnlock()

	tsk, exists := t.tasks[taskID]
	if !exists {
		return nil, ErrTaskNotFound
	}

	// в задании явно не требуется ограничивать доступ,
	// но я ограничу, потому что должен быть смысл в авториазации ¯\_( ͡° ͜ʖ ͡°)_/¯
	if tsk.UserID != user.Login {
		return nil, ErrDontHavePermission
	}

	copy := copyTask(tsk)
	return copy, nil
}

func (t *TaskLocalStorage) GetTasks(user *models.User) ([]*models.Task, error) {
	// в задании явно не указывается какие таски нужно возвращать,
	// поэтому я буду возвращать таски пользователя, потому что должен быть смысл в авториазации ¯\_( ͡° ͜ʖ ͡°)_/¯
	res := make([]*models.Task, 0, 10)
	t.mutex.RLock()
	defer t.mutex.RUnlock()

	for _, value := range t.tasks {
		if value.UserID == user.Login {
			res = append(res, copyTask(value))
		}
	}
	return res, nil
}

func (t *TaskLocalStorage) GetTasksAfter(user *models.User, date time.Time) ([]*models.Task, error) {
	// в задании явно не указывается какие таски нужно возвращать,
	// поэтому я буду возвращать таски пользователя, да-да, потому что должен быть смысл в авториазации ¯\_( ͡° ͜ʖ ͡°)_/¯
	res := make([]*models.Task, 0, 10)
	t.mutex.RLock()
	defer t.mutex.RUnlock()

	for _, value := range t.tasks {
		if value.UserID == user.Login && value.UpdateDate.After(date) {
			res = append(res, copyTask(value))
		}
	}
	return res, nil
}

func copyTask(tsk *models.Task) *models.Task {
	copy := &models.Task{
		ID:          tsk.ID,
		Title:       tsk.Title,
		Description: tsk.Description,
		Status:      tsk.Status,
		UserID:      tsk.UserID,
		UpdateDate:  tsk.UpdateDate,
	}
	return copy
}
