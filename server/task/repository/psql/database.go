package psql

import (
	"context"
	"fmt"
	"task_manager/models"
	"time"

	"github.com/jackc/pgx/v4/pgxpool"
)

type TaskRepository struct {
	conn *pgxpool.Pool
}

func NewTaskRepository(conn *pgxpool.Pool) *TaskRepository {
	return &TaskRepository{
		conn: conn,
	}
}

func (t *TaskRepository) Close() {
	t.conn.Close()
}

func (t *TaskRepository) CreateTask(user *models.User, tsk *models.Task) error {
	tsk.UserID = user.Login

	tr, err := t.conn.Begin(context.Background())
	if err != nil {
		tr.Rollback(context.Background())
		return err
	}

	rows, err := tr.Query(context.Background(), fmt.Sprintf(
		`INSERT INTO public.tasks ("title", "description", "status", "user_id") 
		VALUES ('%s', '%s', '%s', '%s')
		RETURNING id;`,
		tsk.Title, tsk.Description, tsk.Status, user.Login,
	))

	if err != nil {
		tr.Rollback(context.Background())
		return err
	}

	id := -1
	for rows.Next() {
		err := rows.Scan(&id)
		if err != nil {
			// rows.Close()
			tr.Rollback(context.Background())
			return err
		}
	}

	// rows.Close()
	tr.Commit(context.Background())

	tsk.ID = id
	return nil
}

func (t *TaskRepository) EditTask(user *models.User, newTsk *models.Task, oldID int) error {
	tr, err := t.conn.Begin(context.Background())
	if err != nil {
		tr.Rollback(context.Background())
		return err
	}

	rows, err := tr.Query(
		context.Background(),
		`UPDATE public.tasks
		SET title = $1,
			description = $2,
			status = $3,
			user_id = $4,
			update_date = $5
		WHERE id = $6 AND user_id = $7;`,
		newTsk.Title,
		newTsk.Description,
		newTsk.Status,
		newTsk.UserID,
		time.Now(),
		oldID, user.Login,
	)
	if err != nil {
		tr.Rollback(context.Background())
		return err
	}

	rows.Close()
	tr.Commit(context.Background())

	return nil
}

func (t *TaskRepository) DeleteTask(user *models.User, id int) error {
	tr, err := t.conn.Begin(context.Background())
	if err != nil {
		tr.Rollback(context.Background())
		return err
	}

	rows, err := tr.Query(
		context.Background(),
		`DELETE FROM public.tasks
		WHERE id = $1 AND user_id = $2;`,
		id, user.Login,
	)
	if err != nil {
		tr.Rollback(context.Background())
		return err
	}

	rows.Close()
	tr.Commit(context.Background())

	return nil
}

func (t *TaskRepository) GetTask(user *models.User, id int) (*models.Task, error) {
	tsk := &models.Task{}

	tr, err := t.conn.Begin(context.Background())
	if err != nil {
		tr.Rollback(context.Background())
		return nil, err
	}

	rows, err := tr.Query(
		context.Background(),
		`SELECT * FROM public.tasks WHERE id = $1`, id,
	)

	if err != nil {
		tr.Rollback(context.Background())
		return nil, err
	}

	for rows.Next() {
		err := rows.Scan(&tsk.ID, &tsk.Title, &tsk.Description, &tsk.Status, &tsk.UserID, &tsk.UpdateDate)
		if err != nil {
			rows.Close()
			tr.Rollback(context.Background())
			return nil, err
		}
	}

	tr.Commit(context.Background())

	if tsk.ID != id {
		return nil, ErrTaskNotFound
	}

	if tsk.UserID != user.Login {
		return nil, ErrDontHavePermission
	}
	return tsk, nil
}

func (t *TaskRepository) GetTasks(user *models.User) ([]*models.Task, error) {
	tr, err := t.conn.Begin(context.Background())
	if err != nil {
		tr.Rollback(context.Background())
		return nil, err
	}

	rows, err := tr.Query(
		context.Background(),
		fmt.Sprintf(`SELECT * FROM public.tasks WHERE user_id LIKE '%s'`, user.Login),
	)

	if err != nil {
		tr.Rollback(context.Background())
		return nil, err
	}

	res := make([]*models.Task, 0, 10)
	for rows.Next() {
		tsk := &models.Task{}
		err := rows.Scan(&tsk.ID, &tsk.Title, &tsk.Description, &tsk.Status, &tsk.UserID, &tsk.UpdateDate)
		if err != nil {
			// rows.Close()
			tr.Rollback(context.Background())
			return nil, err
		}
		res = append(res, tsk)
	}

	// rows.Close()
	tr.Commit(context.Background())
	return res, nil
}

func (t *TaskRepository) GetTasksAfter(user *models.User, date time.Time) ([]*models.Task, error) {
	tr, err := t.conn.Begin(context.Background())
	if err != nil {
		tr.Rollback(context.Background())
		return nil, err
	}

	rows, err := tr.Query(
		context.Background(),
		fmt.Sprintf(`SELECT * FROM public.tasks WHERE user_id LIKE '%s'`, user.Login),
		// `SELECT * FROM public.tasks WHERE user_id LIKE $1 AND update_date >= $2`, user.Login, date,
	)

	if err != nil {
		tr.Rollback(context.Background())
		return nil, err
	}

	res := make([]*models.Task, 0, 10)
	for rows.Next() {
		tsk := &models.Task{}
		err := rows.Scan(&tsk.ID, &tsk.Title, &tsk.Description, &tsk.Status, &tsk.UserID, &tsk.UpdateDate)
		if err != nil {
			// rows.Close()
			tr.Rollback(context.Background())
			return nil, err
		}
		if tsk.UpdateDate.After(date) {
			res = append(res, tsk)
		}
		// res = append(res, tsk)
	}

	// rows.Close()
	tr.Commit(context.Background())

	return res, nil
}
