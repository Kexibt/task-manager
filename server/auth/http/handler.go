package http

import (
	"encoding/json"
	"io"
	"log"
	"net/http"
	"strconv"
	"task_manager/models"
)

type User struct {
	Login    string `json:"login"`
	Password string `json:"password"`
}

func (a *Auth) signUp(w http.ResponseWriter, r *http.Request) {
	b, err := io.ReadAll(r.Body)
	defer r.Body.Close()
	if err != nil {
		log.Println(err)
		return
	}

	user := &User{}
	json.Unmarshal(b, user)
	if user.Login == "" {
		sendErr(w, ErrNoLogin)
		return
	}

	if user.Password == "" {
		sendErr(w, ErrNoPassword)
		return
	}

	err = a.repository.CreateUser(user.convertToModelsUser())
	if err != nil {
		sendErr(w, err)
		return
	}

	err = sendResult(w, "user successfully created")
	if err != nil {
		log.Println(err)
	}
}

func (a *Auth) signIn(w http.ResponseWriter, r *http.Request) {
	b, err := io.ReadAll(r.Body)
	defer r.Body.Close()
	if err != nil {
		log.Println(err)
		return
	}

	user := &User{}
	json.Unmarshal(b, user)
	if user.Login == "" {
		sendErr(w, ErrNoLogin)
		return
	}

	if user.Password == "" {
		sendErr(w, ErrNoPassword)
		return
	}

	ok, err := a.repository.CheckUser(user.convertToModelsUser())
	if err != nil {
		sendErr(w, err)
		return
	}
	if !ok {
		sendResult(w, "invalid login and password")
	}

	token, err := a.repository.CreateToken(user.convertToModelsUser())
	if err != nil {
		sendErr(w, err)
		return
	}

	cookie := http.Cookie{}
	cookie.Name = "token"
	cookie.Value = token
	http.SetCookie(w, &cookie)

	err = sendResult(w, "success")
	if err != nil {
		log.Println(err)
	}
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
	w.Header().Add("Status", strconv.Itoa(http.StatusOK))
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
	w.Header().Add("Status", strconv.Itoa(http.StatusOK))
	_, err = w.Write(b)
	if err != nil {
		return err
	}
	return nil
}

func (u *User) convertToModelsUser() *models.User {
	return &models.User{
		Login:    u.Login,
		Password: u.Password,
	}
}
