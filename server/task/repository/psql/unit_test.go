package psql

import (
	"context"
	"task_manager/config"
	"task_manager/models"
	"testing"
	"time"

	"github.com/jackc/pgx/v4/pgxpool"
	"github.com/stretchr/testify/assert"
)

func TestCreate(t *testing.T) {
	cfg := config.GetConfig()
	ctx, cancel := context.WithTimeout(context.Background(), cfg.GetConnectionTimeout())
	defer cancel()

	conn, err := connect(ctx, cfg.GetConnectionString())
	assert.NoError(t, err)

	db := NewTaskRepository(conn)
	defer db.Close()

	us := &models.User{Login: "test", Password: "1234"}
	tsk := &models.Task{
		Title:       "tirle",
		Description: "description",
		Status:      "todo",
	}

	err = db.CreateTask(us, tsk)
	assert.NoError(t, err)

	err = db.DeleteTask(us, tsk.ID)
	assert.NoError(t, err)
}

func TestGet(t *testing.T) {
	cfg := config.GetConfig()
	ctx, cancel := context.WithTimeout(context.Background(), cfg.GetConnectionTimeout())
	defer cancel()

	conn, err := connect(ctx, cfg.GetConnectionString())
	assert.NoError(t, err)

	db := NewTaskRepository(conn)
	defer db.Close()

	us := &models.User{Login: "test", Password: "1234"}
	tsk := &models.Task{
		Title:       "title",
		Description: "description",
		Status:      "todo",
	}

	err = db.CreateTask(us, tsk)
	assert.NoError(t, err)

	res, err := db.GetTask(us, tsk.ID)
	assert.NoError(t, err)
	assert.Equal(t, tsk.Title, res.Title)
	assert.Equal(t, tsk.Description, res.Description)
	assert.Equal(t, tsk.Status, res.Status)
	assert.Equal(t, tsk.UserID, res.UserID)
	t.Log("asserted\n")

	err = db.DeleteTask(us, tsk.ID)
	assert.NoError(t, err)
}

func TestEdit(t *testing.T) {
	cfg := config.GetConfig()
	ctx, cancel := context.WithTimeout(context.Background(), cfg.GetConnectionTimeout())
	defer cancel()

	conn, err := connect(ctx, cfg.GetConnectionString())
	assert.NoError(t, err)

	db := NewTaskRepository(conn)
	defer db.Close()

	us := &models.User{Login: "test", Password: "1234"}
	tsk := &models.Task{
		Title:       "tirle",
		Description: "description",
		Status:      "todo",
	}

	err = db.CreateTask(us, tsk)
	assert.NoError(t, err)

	tsk.Title = "new title"
	tsk.Status = "in progress"
	err = db.EditTask(us, tsk, tsk.ID)
	assert.NoError(t, err)

	res, err := db.GetTask(us, tsk.ID)
	assert.NoError(t, err)
	assert.Equal(t, tsk.Title, res.Title)
	assert.Equal(t, tsk.Description, res.Description)
	assert.Equal(t, tsk.Status, res.Status)
	assert.Equal(t, tsk.UserID, res.UserID)

	err = db.DeleteTask(us, tsk.ID)
	assert.NoError(t, err)
}

func TestDelete(t *testing.T) {
	cfg := config.GetConfig()
	ctx, cancel := context.WithTimeout(context.Background(), cfg.GetConnectionTimeout())
	defer cancel()

	conn, err := connect(ctx, cfg.GetConnectionString())
	assert.NoError(t, err)

	db := NewTaskRepository(conn)
	defer db.Close()

	us := &models.User{Login: "test", Password: "1234"}
	tsk := &models.Task{
		Title:       "tirle",
		Description: "description",
		Status:      "todo",
	}

	err = db.CreateTask(us, tsk)
	assert.NoError(t, err)

	err = db.DeleteTask(us, tsk.ID)
	assert.NoError(t, err)

	_, err = db.GetTask(us, tsk.ID)
	assert.Error(t, err)
}

func TestGetTasks(t *testing.T) {
	cfg := config.GetConfig()
	ctx, cancel := context.WithTimeout(context.Background(), cfg.GetConnectionTimeout())
	defer cancel()

	conn, err := connect(ctx, cfg.GetConnectionString())
	assert.NoError(t, err)

	db := NewTaskRepository(conn)
	defer db.Close()

	us := &models.User{Login: "test", Password: "1234"}
	tsk1 := &models.Task{
		Title:       "title1",
		Description: "description1",
		Status:      "todo",
	}

	err = db.CreateTask(us, tsk1)
	assert.NoError(t, err)

	tsk2 := &models.Task{
		Title:       "title2",
		Description: "description2",
		Status:      "todo",
	}

	err = db.CreateTask(us, tsk2)
	assert.NoError(t, err)

	res, err := db.GetTasks(us)
	assert.NoError(t, err)
	assert.Equal(t, 2, len(res))

	assert.Equal(t, tsk1.Title, res[0].Title)
	assert.Equal(t, tsk1.Description, res[0].Description)
	assert.Equal(t, tsk1.Status, res[0].Status)
	assert.Equal(t, tsk1.UserID, res[0].UserID)

	assert.Equal(t, tsk2.Title, res[1].Title)
	assert.Equal(t, tsk2.Description, res[1].Description)
	assert.Equal(t, tsk2.Status, res[1].Status)
	assert.Equal(t, tsk2.UserID, res[1].UserID)

	err = db.DeleteTask(us, tsk1.ID)
	assert.NoError(t, err)

	err = db.DeleteTask(us, tsk2.ID)
	assert.NoError(t, err)
}

func TestGetTasksAfter(t *testing.T) {
	actualTime := time.Now()
	cfg := config.GetConfig()
	ctx, cancel := context.WithTimeout(context.Background(), cfg.GetConnectionTimeout())
	defer cancel()

	conn, err := connect(ctx, cfg.GetConnectionString())
	assert.NoError(t, err)

	db := NewTaskRepository(conn)
	defer db.Close()

	us := &models.User{Login: "test", Password: "1234"}
	tsk1 := &models.Task{
		Title:       "title1",
		Description: "description1",
		Status:      "todo",
	}

	err = db.CreateTask(us, tsk1)
	assert.NoError(t, err)

	tsk2 := &models.Task{
		Title:       "title2",
		Description: "description2",
		Status:      "todo",
	}

	err = db.CreateTask(us, tsk2)
	assert.NoError(t, err)

	res, err := db.GetTasksAfter(us, actualTime)
	assert.NoError(t, err)
	assert.Equal(t, 2, len(res))

	assert.Equal(t, tsk1.Title, res[0].Title)
	assert.Equal(t, tsk1.Description, res[0].Description)
	assert.Equal(t, tsk1.Status, res[0].Status)
	assert.Equal(t, tsk1.UserID, res[0].UserID)

	assert.Equal(t, tsk2.Title, res[1].Title)
	assert.Equal(t, tsk2.Description, res[1].Description)
	assert.Equal(t, tsk2.Status, res[1].Status)
	assert.Equal(t, tsk2.UserID, res[1].UserID)

	err = db.DeleteTask(us, tsk1.ID)
	assert.NoError(t, err)

	err = db.DeleteTask(us, tsk2.ID)
	assert.NoError(t, err)
}

func connect(ctx context.Context, connectionStr string) (conn *pgxpool.Pool, err error) {
	for {
		select {
		case <-ctx.Done():
			return
		default:
			conn, err = pgxpool.Connect(context.Background(), connectionStr)
			if err == nil {
				return
			}
		}
	}
}
