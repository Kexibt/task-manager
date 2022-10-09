package http

import (
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"task_manager/models"
	"task_manager/task"
	"time"

	"github.com/jackc/pgx/v4/pgxpool"
)

type Handler struct {
	tasksRepository *task.TasksRepository
}

func NewHandler(conn *pgxpool.Pool) *Handler {
	return &Handler{
		tasksRepository: task.NewTasksRepository(conn),
	}
}

type createInput struct {
	Title       string `json:"title"`
	Description string `json:"description"`
	Status      string `json:"status"`
}

func (c *createInput) convertToModelsTask() *models.Task {
	return &models.Task{
		Title:       c.Title,
		Description: c.Description,
		Status:      c.Status,
	}
}

func (c *Handler) CreateTask(w http.ResponseWriter, r *http.Request, user *models.User) {
	b, err := io.ReadAll(r.Body)
	if err != nil {
		log.Println(err)
		return
	}
	defer r.Body.Close()

	tsk := &createInput{}
	err = json.Unmarshal(b, tsk)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		err = sendErr(w, err)
		if err != nil {
			log.Println(err)
		}
		return
	}

	modelsTask := tsk.convertToModelsTask()
	err = c.tasksRepository.CreateTask(user, modelsTask)
	if err != nil {
		sendErr(w, err)
		return
	}

	sendResult(w, fmt.Sprintf("successfully added, id: %d", modelsTask.ID))
}

type getInput struct {
	ID int `json:"id"`
}

func (c *Handler) GetTask(w http.ResponseWriter, r *http.Request, user *models.User) {
	b, err := io.ReadAll(r.Body)
	if err != nil {
		log.Println(err)
		return
	}
	defer r.Body.Close()

	input := &getInput{}
	err = json.Unmarshal(b, input)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		err = sendErr(w, err)
		if err != nil {
			log.Println(err)
		}
		return
	}

	tsk, err := c.tasksRepository.GetTask(user, input.ID)
	if err != nil {
		sendErr(w, err)
		return
	}

	by, err := json.Marshal(NewTaskOutput(tsk))
	if err != nil {
		log.Println(err)
		return
	}
	sendBytes(w, by)
}

func (c *Handler) GetTasks(w http.ResponseWriter, r *http.Request, user *models.User) {
	defer r.Body.Close()

	tsks, err := c.tasksRepository.GetTasks(user)
	if err != nil {
		sendErr(w, err)
		return
	}

	res := make([]*TaskOutput, 0, len(tsks))
	for _, tsk := range tsks {
		res = append(res, NewTaskOutput(tsk))
	}

	by, err := json.Marshal(res)
	if err != nil {
		log.Println(err)
		return
	}
	sendBytes(w, by)
}

type getInputDate struct {
	Date time.Time `json:"date"`
}

func (c *Handler) GetTasksAfter(w http.ResponseWriter, r *http.Request, user *models.User) {
	b, err := io.ReadAll(r.Body)
	if err != nil {
		log.Println(err)
		return
	}
	defer r.Body.Close()

	input := &getInputDate{}
	err = json.Unmarshal(b, input)
	if err != nil {
		err = sendErr(w, err)
		if err != nil {
			log.Println(err)
		}
		return
	}

	// loc := time.FixedZone("UTC-8", 3*60*60) // todo переделать
	// if err != nil {
	// 	err = sendErr(w, err)
	// 	if err != nil {
	// 		log.Println(err)
	// 	}
	// 	return
	// }

	// t, err := time.ParseInLocation(time.R, input.Date, loc)
	// if err != nil {
	// 	err = sendErr(w, err)
	// 	if err != nil {
	// 		log.Println(err)
	// 	}
	// 	return
	// }

	tsks, err := c.tasksRepository.GetTasksAfter(user, input.Date)
	if err != nil {
		sendErr(w, err)
		return
	}

	res := make([]*TaskOutput, 0, len(tsks))
	for _, tsk := range tsks {
		res = append(res, NewTaskOutput(tsk))
	}

	by, err := json.Marshal(res)
	if err != nil {
		log.Println(err)
		return
	}
	sendBytes(w, by)
}

type TaskOutput struct {
	ID          int    `json:"id"`
	Title       string `json:"title"`
	Description string `json:"description"`
	Status      string `json:"status"`
	UserID      string `json:"userid"`
	UpdateDate  string `json:"last update date"`
}

func NewTaskOutput(tsk *models.Task) *TaskOutput {
	return &TaskOutput{
		ID:          tsk.ID,
		Title:       tsk.Title,
		Description: tsk.Description,
		Status:      tsk.Status,
		UserID:      tsk.UserID,
		UpdateDate:  tsk.UpdateDate.Local().Format(time.UnixDate),
	}
}

func sendBytes(w http.ResponseWriter, b []byte) error {
	w.Header().Add("Content-Type", "application/json")
	_, err := w.Write(b)
	if err != nil {
		return err
	}
	return nil
}

type ResultResponse struct {
	Result string `json:"result"`
}

func sendResult(w http.ResponseWriter, msg string) error {
	r := ResultResponse{
		Result: msg,
	}

	b, err := json.Marshal(r)
	if err != nil {
		return err
	}

	w.Header().Add("Content-Type", "application/json")
	_, err = w.Write(b)
	if err != nil {
		return err
	}
	return nil
}

type ErrorResponse struct {
	Error string `json:"error"`
}

func sendErr(w http.ResponseWriter, err error) error {
	e := ErrorResponse{
		Error: err.Error(),
	}

	b, err := json.Marshal(e)
	if err != nil {
		return err
	}

	w.Header().Add("Content-Type", "application/json")
	w.WriteHeader(http.StatusBadRequest)
	_, err = w.Write(b)
	if err != nil {
		return err
	}
	return nil
}
