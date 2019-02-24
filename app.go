package main

import (
	"net/http"
	"path"
	"strings"

	"github.com/hacknights/middleware"
	"github.com/tidwall/buntdb"
)

type app struct {
	db    *buntdb.DB
	users *usersHandler
	auth  *authHandler
}

func newAppHandler() http.Handler {
	db, err := buntdb.Open(":memory:")
	if err != nil {

	}

	a := &app{
		db:    db,
		users: newUsersHandler(db),
		auth:  newAuthHandler(db),
	}

	return middleware.Use(
		middleware.TraceIDs,
		middleware.Logging,
		middleware.Recuperate)(a.ServeHTTP)
}

func (a *app) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	var head string
	head, r.URL.Path = shiftPath(r.URL.Path)

	switch head {
	case "users":
		a.users.ServeHTTP(w, r)
	case "authenticate":
		a.auth.ServeHTTP(w, r)
	default:
		notFoundError(w)
	}
}

func shiftPath(p string) (head, tail string) {
	p = path.Clean("/" + p)
	i := strings.Index(p[1:], "/") + 1
	if i <= 0 {
		return p[1:], "/"
	}
	return p[1:i], p[i:]
}

func notFoundError(w http.ResponseWriter) {
	http.Error(w, "Not Found", http.StatusNotFound)
}

func internalServerError(w http.ResponseWriter, err error) {
	http.Error(w, err.Error(), http.StatusInternalServerError)
}

func unauthorizedError(w http.ResponseWriter, err error) {
	http.Error(w, err.Error(), http.StatusUnauthorized)
}
