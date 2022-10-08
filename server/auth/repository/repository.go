package repository

import (
	"crypto/sha256"
	"encoding/hex"
	"task_manager/auth/repository/psql"
	"task_manager/auth/repository/tokens"
	"task_manager/models"

	"github.com/jackc/pgx/v4/pgxpool"
)

type TokenRepository interface {
	CreateToken(user *models.User) (string, error)
	GetUserByToken(token string) (*models.User, error)
}

type UserRepository interface {
	CreateUser(user *models.User) error
	CheckUser(username, password string) (bool, error)
}

type Repository struct {
	tokens TokenRepository
	users  UserRepository
}

func NewRepository(conn *pgxpool.Pool) *Repository {
	return &Repository{
		tokens: tokens.NewTokenRepository(),
		users:  psql.NewUserRepository(conn),
	}
}

func (r *Repository) CreateToken(user *models.User) (string, error) {
	return r.tokens.CreateToken(user)
}

func (r *Repository) GetUserByToken(token string) (*models.User, error) {
	return r.tokens.GetUserByToken(token)
}

func (r *Repository) CreateUser(user *models.User) error {
	h := sha256.New()
	hashpass := hex.EncodeToString(h.Sum([]byte(user.Login)))
	user.Password = hashpass
	return r.users.CreateUser(user)
}

func (r *Repository) CheckUser(username, password string) (bool, error) {
	h := sha256.New()
	hashpass := hex.EncodeToString(h.Sum([]byte(password)))
	return r.users.CheckUser(username, hashpass)
}
