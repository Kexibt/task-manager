package tokens

import (
	"task_manager/models"
	"testing"

	"github.com/stretchr/testify/assert"
)

// func TestCreate(t *testing.T) {
// 	rep := NewUserLocalStore()
// 	usr := &models.User{Login: "test", Password: "pass"}

// 	err := rep.CreateUser(usr)
// 	assert.NoError(t, err)

// 	assert.Equal(t, usr, rep.users[usr.Login])
// 	assert.NotSame(t, usr, rep.users[usr.Login])
// }

// func TestNegativeCreate(t *testing.T) {
// 	rep := NewUserLocalStore()
// 	usr := &models.User{Login: "test", Password: "pass"}

// 	err := rep.CreateUser(usr)
// 	assert.NoError(t, err)

// 	err = rep.CreateUser(usr)
// 	assert.Error(t, err)
// }

// func TestCheck(t *testing.T) {
// 	rep := NewUserLocalStore()
// 	usr := &models.User{Login: "test", Password: "pass"}

// 	rep.users[usr.Login] = usr
// 	ok, err := rep.CheckUser(usr.Login, usr.Password)

// 	assert.NoError(t, err)
// 	assert.True(t, ok)
// }

// func TestNegativeCheck(t *testing.T) {
// 	rep := NewUserLocalStore()
// 	usr := &models.User{Login: "test", Password: "pass"}
// 	usr1 := &models.User{Login: "test1", Password: "pass1"}

// 	rep.users[usr.Login] = usr
// 	ok, err := rep.CheckUser(usr1.Login, usr1.Password)

// 	assert.Error(t, err)
// 	assert.False(t, ok)
// }

func TestCreateToken(t *testing.T) {
	rep := NewTokenRepository()
	usr := &models.User{Login: "test", Password: "pass"}

	token, err := rep.CreateToken(usr)

	assert.NoError(t, err)
	assert.NotEqual(t, 0, len(token))
	assert.Equal(t, usr, rep.tokens[token])
	assert.NotSame(t, usr, rep.tokens[token])
}

func TestGetUserByToken(t *testing.T) {
	rep := NewTokenRepository()
	usr := &models.User{Login: "test", Password: "pass"}

	token, err := rep.CreateToken(usr)

	assert.NoError(t, err)
	assert.NotEqual(t, 0, len(token))

	u, err := rep.GetUserByToken(token)
	assert.NoError(t, err)
	assert.Equal(t, usr, u)
	assert.NotSame(t, usr, u)
}
