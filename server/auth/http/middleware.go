package http

import (
	"net/http"
	"strings"
	"task_manager/models"
)

type InternalHandler interface {
	Handler() func(w http.ResponseWriter, r *http.Request, user *models.User)
}

func (a *Auth) MiddleAuth(w http.ResponseWriter, r *http.Request) {
	tokenCookie, err := r.Cookie("token")
	if err != nil {
		sendErr(w, ErrAuthorization)
		return
	}

	usr, err := a.repository.GetUserByToken(tokenCookie.Value)
	if err != nil {
		sendErr(w, err)
		return
	}

	str := r.URL.String()
	str = str[strings.LastIndex(str, "/"):]

	a.internalHandlers[str](w, r, usr)
}
