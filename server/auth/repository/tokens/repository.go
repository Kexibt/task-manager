package tokens

import (
	"crypto/sha256"
	"encoding/hex"
	"sync"
	"time"

	"task_manager/models"
)

type TokenRepository struct {
	tokens      map[string]*models.User
	tokensMutex sync.RWMutex

	// users      map[string]*models.User
	// usersMutex sync.RWMutex
}

func NewTokenRepository() *TokenRepository {
	return &TokenRepository{
		tokens: make(map[string]*models.User),
		// users:  make(map[string]*models.User),
	}
}

// func (u *UserLocalStore) CreateUser(user *models.User) error {
// 	u.usersMutex.Lock()
// 	defer u.usersMutex.Unlock()

// 	if _, exists := u.users[user.Login]; exists {
// 		return ErrUserAlreadyExists
// 	}

// 	u.users[user.Login] = &models.User{
// 		Login:    user.Login,
// 		Password: user.Password,
// 	}

// 	return nil
// }

// func (u *UserLocalStore) CheckUser(username, password string) (bool, error) {
// 	u.usersMutex.RLock()
// 	defer u.usersMutex.RUnlock()

// 	if _, exists := u.users[username]; !exists {
// 		return false, ErrUserNotFound
// 	}

// 	if u.users[username].Password != password {
// 		return false, ErrWrongPassword
// 	}
// 	return true, nil
// }

func (t *TokenRepository) CreateToken(user *models.User) (string, error) {
	userTimestamp := user.Login + user.Password + time.Now().String()
	h := sha256.New()
	token := hex.EncodeToString(h.Sum([]byte(userTimestamp)))

	t.tokensMutex.Lock()
	defer t.tokensMutex.Unlock()

	for _, exists := t.tokens[token]; exists; _, exists = t.tokens[token] {
		userTimestamp = user.Login + user.Password + time.Now().String()
		token = hex.EncodeToString(h.Sum([]byte(userTimestamp)))
	}

	t.tokens[token] = &models.User{
		Login:    user.Login,
		Password: user.Password,
	}

	return token, nil
}

func (t *TokenRepository) GetUserByToken(token string) (*models.User, error) {
	t.tokensMutex.RLock()
	defer t.tokensMutex.RUnlock()

	if _, exists := t.tokens[token]; exists {
		return &models.User{
			Login:    t.tokens[token].Login,
			Password: t.tokens[token].Password,
		}, nil
	}

	return nil, ErrInvalidAccessToken
}
