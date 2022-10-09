package http

import (
	"net/http"
	"task_manager/auth/repository"
	"task_manager/models"

	"github.com/jackc/pgx/v4/pgxpool"
)

type Auth struct {
	repository       *repository.Repository
	handlers         map[string]func(w http.ResponseWriter, r *http.Request)
	internalHandlers map[string]func(w http.ResponseWriter, r *http.Request, user *models.User)
}

func NewAuth(conn *pgxpool.Pool) *Auth {
	auth := &Auth{
		repository: repository.NewRepository(conn),
	}

	auth.handlers = map[string]func(w http.ResponseWriter, r *http.Request){
		"/sign_in": auth.signIn,
		"/sign_up": auth.signUp,
	}

	auth.internalHandlers = make(map[string]func(w http.ResponseWriter, r *http.Request, user *models.User))

	return auth
}

func (a *Auth) AddNewMiddlewaredHandler(pattern string, internal func(w http.ResponseWriter, r *http.Request, user *models.User)) {
	a.handlers[pattern] = a.MiddleAuth
	a.internalHandlers[pattern] = internal
}

func (a *Auth) AddHandlers(mux *http.ServeMux) {
	for pattern, handler := range a.handlers {
		mux.HandleFunc(pattern, handler)
	}
}
