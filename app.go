package main

import (
	"net/http"
	"path"
	"strings"
	"usery/db"

	"github.com/hacknights/middleware"
	"github.com/hacknights/negotiator"
	"github.com/tidwall/buntdb"
)

type app struct {
	negotiator negotiator.Factory
	db         *db.DB
	apps       *appsHandler
	// users      *usersHandler
	// auth       *authHandler
}

func newAppHandler() http.Handler {
	bdb, err := buntdb.Open(":memory:")
	if err != nil {

	}

	db := db.NewDB(bdb)
	n := negotiator.NewNegotiator
	a := &app{
		negotiator: n,
		db:         db,
		apps:       newAppsHandler(n, db),
		// users:      newUsersHandler(n, db),
		// auth:       newAuthHandler(n, db),
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
	case "apps":
		a.apps.ServeHTTP(w, r)
	// case "users":
	// 	a.users.ServeHTTP(w, r)
	// case "authenticate":
	// 	a.auth.ServeHTTP(w, r)
	default:
		a.negotiator(w, r).NotFound()
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
