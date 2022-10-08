package repository

import (
	"task_manager/models"
	"time"
)

type Repository interface {
	CreateTask(user *models.User, tsk *models.Task) error
	EditTask(user *models.User, newTsk *models.Task, oldID int) error
	DeleteTask(user *models.User, id int) error

	GetTask(user *models.User, taskID int) (*models.Task, error)
	GetTasks(user *models.User) ([]*models.Task, error)
	GetTasksAfter(user *models.User, date time.Time) ([]*models.Task, error)
}

// type TokenRepository interface {
// 	CreateToken(user *models.User) (string, error)
// 	GetUserByToken(token string) (*models.User, error)
// }

// type UserRepository interface {
// 	CreateUser(user *models.User) error
// 	CheckUser(username, password string) (bool, error)
// }

// type Repository struct {
// 	tokens TokenRepository
// 	users  UserRepository
// }

// func NewRepository(conn *pgxpool.Pool) *Repository {
// 	return &Repository{
// 		tokens: tokens.NewTokenRepository(),
// 		users:  psql.NewUserRepository(conn),
// 	}
// }

// func (r *Repository) CreateToken(user *models.User) (string, error) {
// 	return r.tokens.CreateToken(user)
// }

// func (r *Repository) GetUserByToken(token string) (*models.User, error) {
// 	return r.tokens.GetUserByToken(token)
// }

// func (r *Repository) CreateUser(user *models.User) error {
// 	h := sha256.New()
// 	hashpass := hex.EncodeToString(h.Sum([]byte(user.Login)))
// 	user.Password = hashpass
// 	return r.users.CreateUser(user)
// }

// func (r *Repository) CheckUser(username, password string) (bool, error) {
// 	h := sha256.New()
// 	hashpass := hex.EncodeToString(h.Sum([]byte(password)))
// 	return r.users.CheckUser(username, hashpass)
// }
