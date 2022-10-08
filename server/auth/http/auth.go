package http

import (
	"task_manager/auth/repository"

	"github.com/jackc/pgx/v4/pgxpool"
)

type Auth struct {
	repository *repository.Repository
	handlers   *Handlers
}

func NewAuth(conn *pgxpool.Pool) *Auth {
	return &Auth{
		repository: repository.NewRepository(conn),
		handlers:   NewHandlers(),
	}
}
