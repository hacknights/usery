package main

import (
	"net/http"
	"path"
	"strings"

	"github.com/hacknights/middleware"
	"github.com/tidwall/buntdb"
)

type app struct {
	db  *buntdb.DB
	api *apiHandler
}

func newAppHandler() http.Handler {
	db, err := buntdb.Open(":memory:")
	if err != nil {

	}

	a := &app{
		db:  db,
		api: newApiHandler(db),
	}

	return middleware.Use(
		middleware.TraceIDs,
		middleware.Logging,
		middleware.Recuperate)(a.ServeHTTP)
}

func (a *app) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	var head string
	head, r.URL.Path = shiftPath(r.URL.Path)

	if head == "api" {
		a.api.ServeHTTP(w, r)
		return
	}

	notFoundError(w)
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
