package psql

import (
	"context"
	"task_manager/models"

	"github.com/jackc/pgx/v4/pgxpool"
)

type UserRepository struct {
	conn *pgxpool.Pool
}

func NewUserRepository(conn *pgxpool.Pool) *UserRepository {
	return &UserRepository{
		conn: conn,
	}
}

func (u *UserRepository) Close() {
	u.conn.Close()
}

func (u *UserRepository) CreateUser(user *models.User) error {
	tr, err := u.conn.Begin(context.Background())
	if err != nil {
		tr.Rollback(context.Background())
		return err
	}

	rows, err := tr.Query(
		context.Background(),
		`SELECT * FROM public.users
		WHERE login LIKE $1;`,
		user.Login,
	)
	if err != nil {
		tr.Rollback(context.Background())
		return err
	}
	for rows.Next() {
		usr := &models.User{}

		err = rows.Scan(&usr.Login, &usr.Password)
		if err != nil {
			rows.Close()
			tr.Rollback(context.Background())
			return err
		}

		if usr.Login == user.Login {
			rows.Close()
			tr.Rollback(context.Background())
			return ErrUserAlreadyExists
		}
	}

	rows, err = tr.Query(context.Background(),
		`INSERT INTO public.users ("login", "hashpass") 
		VALUES ($1, $2);`,
		user.Login, user.Password,
	)

	if err != nil {
		tr.Rollback(context.Background())
		return err
	}

	rows.Close()

	return tr.Commit(context.Background())
}

func (u *UserRepository) CheckUser(username, password string) (bool, error) {
	tr, err := u.conn.Begin(context.Background())
	if err != nil {
		tr.Rollback(context.Background())
		return false, err
	}

	rows, err := tr.Query(
		context.Background(),
		`SELECT * FROM public.users
		WHERE login LIKE $1;`,
		username,
	)
	if err != nil {
		tr.Rollback(context.Background())
		return false, err
	}

	user := &models.User{}
	if rows.Next() {
		err = rows.Scan(&user.Login, &user.Password)
		if err != nil {
			rows.Close()
			tr.Rollback(context.Background())
			return false, err
		}
	}

	rows.Close()
	tr.Commit(context.Background())

	if user.Login != username {
		return false, ErrUserNotFound
	}

	if user.Password != password {
		return false, ErrWrongPassword
	}

	return true, nil
}
