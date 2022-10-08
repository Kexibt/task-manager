package psql

import (
	"context"
	"task_manager/config"
	"task_manager/models"
	"testing"

	"github.com/jackc/pgx/v4/pgxpool"
	"github.com/stretchr/testify/assert"
)

func TestCreate(t *testing.T) {
	cfg := config.GetConfig()
	ctx, cancel := context.WithTimeout(context.Background(), cfg.GetConnectionTimeout())
	defer cancel()

	conn, err := connect(ctx, cfg.GetConnectionString())
	assert.NoError(t, err)

	rep := NewUserRepository(conn)
	usr := &models.User{Login: "test1", Password: "pass1"}

	err = rep.CreateUser(usr)
	assert.NoError(t, err)

}

func TestNegativeCreate(t *testing.T) {
	cfg := config.GetConfig()
	ctx, cancel := context.WithTimeout(context.Background(), cfg.GetConnectionTimeout())
	defer cancel()

	conn, err := connect(ctx, cfg.GetConnectionString())
	assert.NoError(t, err)

	rep := NewUserRepository(conn)
	usr := &models.User{Login: "test2", Password: "pass2"}

	err = rep.CreateUser(usr)
	assert.NoError(t, err)

	usr = &models.User{Login: "test2", Password: "pass3"}
	err = rep.CreateUser(usr)
	assert.Error(t, err)
}

func TestCheck(t *testing.T) {
	cfg := config.GetConfig()
	ctx, cancel := context.WithTimeout(context.Background(), cfg.GetConnectionTimeout())
	defer cancel()

	conn, err := connect(ctx, cfg.GetConnectionString())
	assert.NoError(t, err)

	rep := NewUserRepository(conn)
	usr := &models.User{Login: "test3", Password: "pass3"}

	rep.CreateUser(usr)
	ok, err := rep.CheckUser(usr.Login, usr.Password)

	assert.NoError(t, err)
	assert.True(t, ok)
}

func TestNegativeCheck(t *testing.T) {
	cfg := config.GetConfig()
	ctx, cancel := context.WithTimeout(context.Background(), cfg.GetConnectionTimeout())
	defer cancel()

	conn, err := connect(ctx, cfg.GetConnectionString())
	assert.NoError(t, err)

	rep := NewUserRepository(conn)
	usr := &models.User{Login: "test4", Password: "pass4"}
	usr1 := &models.User{Login: "test5", Password: "pass5"}

	rep.CreateUser(usr)
	ok, err := rep.CheckUser(usr1.Login, usr1.Password)

	assert.Error(t, err)
	assert.False(t, ok)
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
