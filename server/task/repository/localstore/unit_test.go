package localstore

import (
	"task_manager/models"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestCreateTask(t *testing.T) {
	st := NewLocalStorage()

	usr := &models.User{
		Login:    "login",
		Password: "pass",
	}

	tsk := &models.Task{
		ID:          0,
		Title:       "title",
		Description: "description",
		Status:      "in progress",
	}

	err := st.CreateTask(usr, tsk)
	assert.NoError(t, err)

	exp := &models.Task{
		ID:          0,
		Title:       "title",
		Description: "description",
		Status:      "in progress",
		UserID:      usr.Login,
	}

	assert.True(t, isEqual(st.tasks[0], exp))
	assert.NotSame(t, tsk, st.tasks[0])
}

func TestCreateInvalidStatus(t *testing.T) {
	st := NewLocalStorage()

	usr := &models.User{
		Login:    "login",
		Password: "pass",
	}

	tsk := &models.Task{
		ID:          0,
		Title:       "title",
		Description: "description",
		Status:      "invalid",
	}

	err := st.CreateTask(usr, tsk)
	assert.ErrorIs(t, err, ErrInvalidStatus)

	assert.Nil(t, st.tasks[0])
}

func TestGetTask(t *testing.T) {
	st := NewLocalStorage()

	usr := &models.User{
		Login:    "login",
		Password: "pass",
	}

	tsk := &models.Task{
		ID:          0,
		Title:       "title",
		Description: "description",
		Status:      "in progress",
		UserID:      usr.Login,
		UpdateDate:  time.Now(),
	}

	st.tasks[0] = tsk

	res, err := st.GetTask(usr, 0)
	assert.NoError(t, err)
	assert.Equal(t, res, tsk)
	assert.NotSame(t, res, tsk)
}

func TestGetNilTask(t *testing.T) {
	st := NewLocalStorage()

	usr := &models.User{
		Login:    "login",
		Password: "pass",
	}

	res, err := st.GetTask(usr, 0)
	assert.ErrorIs(t, err, ErrTaskNotFound)
	assert.Nil(t, res)
}

func TestGetInvalidUser(t *testing.T) {
	st := NewLocalStorage()

	usr1 := &models.User{
		Login:    "login1",
		Password: "pass",
	}

	usr2 := &models.User{
		Login:    "login2",
		Password: "pass",
	}

	tsk := &models.Task{
		ID:          0,
		Title:       "title",
		Description: "description",
		Status:      "in progress",
		UserID:      usr1.Login,
		UpdateDate:  time.Now(),
	}

	st.tasks[0] = tsk

	res, err := st.GetTask(usr2, 0)
	assert.ErrorIs(t, err, ErrDontHavePermission)
	assert.Nil(t, res)
}

func TestEdit(t *testing.T) {
	st := NewLocalStorage()

	usr := &models.User{
		Login:    "login",
		Password: "pass",
	}

	tsk := &models.Task{
		ID:          0,
		Title:       "title",
		Description: "description",
		Status:      "in progress",
		UserID:      usr.Login,
		UpdateDate:  time.Now(),
	}

	newTsk := &models.Task{
		ID:          0,
		Title:       "changed title",
		Description: "description",
		Status:      "todo",
		UserID:      usr.Login,
		UpdateDate:  time.Now(),
	}

	st.tasks[0] = tsk
	err := st.EditTask(usr, newTsk, 0)

	assert.NoError(t, err)
	assert.Equal(t, st.tasks[0], newTsk)
}

func TestEditInvalidStatus(t *testing.T) {
	st := NewLocalStorage()

	usr := &models.User{
		Login:    "login",
		Password: "pass",
	}

	tsk := &models.Task{
		ID:          0,
		Title:       "title",
		Description: "description",
		Status:      "in progress",
		UserID:      usr.Login,
		UpdateDate:  time.Now(),
	}

	newTsk := &models.Task{
		ID:          0,
		Title:       "changed title",
		Description: "description",
		Status:      "invalid",
		UserID:      usr.Login,
		UpdateDate:  time.Now(),
	}

	st.CreateTask(usr, tsk)
	err := st.EditTask(usr, newTsk, 0)

	assert.ErrorIs(t, err, ErrInvalidStatus)
	assert.Equal(t, st.tasks[0], tsk)
	assert.NotEqual(t, st.tasks[0], newTsk)
	assert.NotSame(t, tsk, st.tasks[0])
}

func TestEditInvalidUser(t *testing.T) {
	st := NewLocalStorage()

	usr1 := &models.User{
		Login:    "login1",
		Password: "pass",
	}

	usr2 := &models.User{
		Login:    "login2",
		Password: "pass",
	}

	tsk := &models.Task{
		ID:          0,
		Title:       "title",
		Description: "description",
		Status:      "in progress",
		UserID:      usr1.Login,
		UpdateDate:  time.Now(),
	}

	newTsk := &models.Task{
		ID:          0,
		Title:       "changed title",
		Description: "description",
		Status:      "todo",
		UserID:      usr2.Login,
		UpdateDate:  time.Now(),
	}

	st.CreateTask(usr1, tsk)
	err := st.EditTask(usr2, newTsk, 0)

	assert.ErrorIs(t, err, ErrDontHavePermission)
	assert.Equal(t, st.tasks[0], tsk)
	assert.NotEqual(t, st.tasks[0], newTsk)
	assert.NotSame(t, tsk, st.tasks[0])
}

func TestDelete(t *testing.T) {
	st := NewLocalStorage()

	usr := &models.User{
		Login:    "login",
		Password: "pass",
	}

	tsk := &models.Task{
		ID:          0,
		Title:       "title",
		Description: "description",
		Status:      "in progress",
		UserID:      usr.Login,
		UpdateDate:  time.Now(),
	}

	st.CreateTask(usr, tsk)
	err := st.DeleteTask(usr, 0)

	assert.NoError(t, err)
	assert.Nil(t, st.tasks[0])
}

func TestDeleteInvalidUser(t *testing.T) {
	st := NewLocalStorage()

	usr1 := &models.User{
		Login:    "login1",
		Password: "pass",
	}

	usr2 := &models.User{
		Login:    "login2",
		Password: "pass",
	}

	tsk := &models.Task{
		ID:          0,
		Title:       "title",
		Description: "description",
		Status:      "in progress",
		UserID:      usr1.Login,
		UpdateDate:  time.Now(),
	}

	st.CreateTask(usr1, tsk)
	err := st.DeleteTask(usr2, 0)

	assert.ErrorIs(t, err, ErrDontHavePermission)
	assert.NotNil(t, st.tasks[0])
	assert.Equal(t, st.tasks[0], tsk)
}

func TestGetTasks(t *testing.T) {
	st := NewLocalStorage()

	usr := &models.User{
		Login:    "login",
		Password: "pass",
	}
	usr1 := &models.User{
		Login:    "another login",
		Password: "pass",
	}

	tsk1 := &models.Task{
		ID:          1,
		Title:       "title1",
		Description: "description1",
		Status:      "in progress",
		UserID:      usr.Login,
		UpdateDate:  time.Now(),
	}
	tsk2 := &models.Task{
		ID:          2,
		Title:       "title2",
		Description: "description2",
		Status:      "done",
		UserID:      usr.Login,
		UpdateDate:  time.Now(),
	}
	tsk3 := &models.Task{
		ID:          3,
		Title:       "title3",
		Description: "description3",
		Status:      "todo",
		UserID:      usr.Login,
		UpdateDate:  time.Now(),
	}
	tsk4 := &models.Task{
		ID:          4,
		Title:       "title4",
		Description: "description4",
		Status:      "todo",
		UserID:      usr1.Login,
		UpdateDate:  time.Now(),
	}

	st.CreateTask(usr, tsk1)
	st.CreateTask(usr, tsk2)
	st.CreateTask(usr, tsk3)
	st.CreateTask(usr1, tsk4)

	res, err := st.GetTasks(usr)
	assert.NoError(t, err)
	assert.Equal(t, len(res), 3)

	for _, r := range res {
		assert.True(t, isEqual(r, tsk1) || isEqual(r, tsk2) || isEqual(r, tsk3))
	}
}

func TestGetNilTasks(t *testing.T) {
	st := NewLocalStorage()

	usr := &models.User{
		Login:    "login",
		Password: "pass",
	}

	res, _ := st.GetTasks(usr)
	// assert.ErrorIs(t, err, ErrTaskNotFound)
	assert.Equal(t, len(res), 0)
}

func TestGetTasksInvalidUser(t *testing.T) {
	st := NewLocalStorage()

	usr1 := &models.User{
		Login:    "login1",
		Password: "pass",
	}

	usr2 := &models.User{
		Login:    "login2",
		Password: "pass",
	}

	tsk1 := &models.Task{
		ID:          1,
		Title:       "title1",
		Description: "description1",
		Status:      "in progress",
	}
	tsk2 := &models.Task{
		ID:          2,
		Title:       "title2",
		Description: "description2",
		Status:      "done",
	}
	tsk3 := &models.Task{
		ID:          3,
		Title:       "title3",
		Description: "description3",
		Status:      "todo",
	}

	st.CreateTask(usr1, tsk1)
	st.CreateTask(usr1, tsk2)
	st.CreateTask(usr1, tsk3)

	res, _ := st.GetTasks(usr2)
	// assert.ErrorIs(t, err, ErrDontHavePermission)
	assert.Equal(t, 0, len(res))
}
func TestGetTasksAfter(t *testing.T) {
	st := NewLocalStorage()

	usr := &models.User{
		Login:    "login",
		Password: "pass",
	}
	usr1 := &models.User{
		Login:    "another login",
		Password: "pass",
	}

	tsk1 := &models.Task{
		ID:          1,
		Title:       "title1",
		Description: "description1",
		Status:      "in progress",
	}
	tsk2 := &models.Task{
		ID:          2,
		Title:       "title2",
		Description: "description2",
		Status:      "done",
	}
	tsk3 := &models.Task{
		ID:          3,
		Title:       "title3",
		Description: "description3",
		Status:      "todo",
	}
	tsk4 := &models.Task{
		ID:          4,
		Title:       "title4",
		Description: "description4",
		Status:      "todo",
	}

	st.CreateTask(usr, tsk1)
	st.CreateTask(usr, tsk2)

	now := time.Now()

	st.CreateTask(usr, tsk3)
	st.CreateTask(usr1, tsk4)

	res, err := st.GetTasksAfter(usr, now)
	assert.NoError(t, err)
	assert.Equal(t, 1, len(res))

	assert.True(t, isEqual(tsk3, res[0]))
}

func TestGetNilTasksAfter(t *testing.T) {
	st := NewLocalStorage()

	usr := &models.User{
		Login:    "login",
		Password: "pass",
	}

	tsk := &models.Task{
		ID: 0,
		// bla-bla
	}

	st.CreateTask(usr, tsk)

	res, _ := st.GetTasksAfter(usr, time.Now())
	// assert.ErrorIs(t, err, ErrTaskNotFound)
	assert.Equal(t, len(res), 0)
}

func TestGetTasksAfterInvalidUser(t *testing.T) {
	st := NewLocalStorage()

	usr1 := &models.User{
		Login:    "login1",
		Password: "pass",
	}

	usr2 := &models.User{
		Login:    "login2",
		Password: "pass",
	}

	tsk1 := &models.Task{
		ID:          1,
		Title:       "title1",
		Description: "description1",
		Status:      "in progress",
	}
	tsk2 := &models.Task{
		ID:          2,
		Title:       "title2",
		Description: "description2",
		Status:      "done",
	}
	tsk3 := &models.Task{
		ID:          3,
		Title:       "title3",
		Description: "description3",
		Status:      "todo",
	}

	st.CreateTask(usr1, tsk1)

	now := time.Now()

	st.CreateTask(usr1, tsk2)
	st.CreateTask(usr1, tsk3)

	res, _ := st.GetTasksAfter(usr2, now)
	// assert.ErrorIs(t, err, ErrDontHavePermission)
	assert.Equal(t, 0, len(res))
}

func isEqual(left, right *models.Task) bool {
	return left.Title == right.Title &&
		left.Description == right.Description &&
		left.ID == right.ID &&
		left.Status == right.Status &&
		left.UserID == right.UserID
}
