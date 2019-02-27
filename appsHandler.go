package main

import (
	"encoding/json"
	"net/http"
	"usery/db"
	"usery/requests"
	"usery/responses"

	"github.com/hacknights/negotiator"
)

type appsHandler struct {
	negotiator negotiator.Factory
	db         *db.DB
}

func newAppsHandler(n negotiator.Factory, db *db.DB) *appsHandler {
	return &appsHandler{
		negotiator: n,
		db:         db,
	}
}

func (h *appsHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case "POST":
		h.handlePostApp(w, r)
	default:
		h.negotiator(w, r).NotFound()
	}
}

func (h *appsHandler) handlePostApp(w http.ResponseWriter, r *http.Request) {
	n := h.negotiator(w, r)
	decoder := json.NewDecoder(r.Body)
	var ar requests.NewApp
	if err := decoder.Decode(&ar); err != nil {
		n.InternalServerError(err)
		return
	}

	if err := ar.Validate(); err != nil {
		n.BadRequestError(err)
		return
	}

	e := db.NewApp(ar.Appname)
	i, err := h.db.Insert(e)
	if err != nil {
		n.InternalServerError(err)
		return
	}

	n.OK(responses.NewNewApp(i.Identify(), e.Appname))
}
