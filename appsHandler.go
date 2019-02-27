package main

import (
	"encoding/json"
	"net/http"

	"github.com/google/uuid"

	"github.com/hacknights/negotiator"

	"github.com/tidwall/buntdb"
)

type appsHandler struct {
	negotiator negotiator.Factory
	db         *buntdb.DB
}

func newAppsHandler(n negotiator.Factory, db *buntdb.DB) *appsHandler {
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

type dbApp struct {
	ID      string   `json:"id"`
	Appname string   `json:"appname"`
	Users   []dbUser `json:"users"`
}

func newDbApp(appname string) *dbApp {
	return &dbApp{
		ID:      uuid.New().String(),
		Appname: appname,
	}
}

func (h *appsHandler) handlePostApp(w http.ResponseWriter, r *http.Request) {
	n := h.negotiator(w, r)
	decoder := json.NewDecoder(r.Body)
	var ar newAppRequest
	if err := decoder.Decode(&ar); err != nil {
		n.InternalServerError(err)
		return
	}

	if err := ar.validateRequest(); err != nil {
		n.BadRequestError(err)
		return
	}

	dbApp := newDbApp(ar.Appname)
	b, err := json.Marshal(dbApp)
	if err != nil {
		n.InternalServerError(err)
		return
	}

	if err := h.db.Update(func(tx *buntdb.Tx) error {
		if _, _, err := tx.Set(dbApp.ID, string(b), nil); err != nil {
			return err
		}
		n.OK(dbApp.ID)
		return nil
	}); err != nil {
		n.InternalServerError(err)
		return
	}
}
