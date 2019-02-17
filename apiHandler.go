package main

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/tidwall/buntdb"
)

type apiHandler struct {
	db *buntdb.DB
}

func newApiHandler(db *buntdb.DB) *apiHandler {
	return &apiHandler{
		db: db,
	}
}

func (a *apiHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	var head string
	head, r.URL.Path = shiftPath(r.URL.Path)

	switch head {
	case "users":
		a.handleUsers(w, r)
	case "authenticate":
		a.handleAuthenticate(w, r)
	default:
		notFoundError(w)
	}
}

type user struct {
	Username string                 `json:"username"`
	Email    string                 `json:"email"`
	Password string                 `json:"password"`
	Claims   map[string]interface{} `json:"claims"`
}

func (a *apiHandler) handleAuthenticate(w http.ResponseWriter, r *http.Request) {
	decoder := json.NewDecoder(r.Body)
	var u user
	if err := decoder.Decode(&u); err != nil {
		internalServerError(w, err)
		return
	}

	if err := a.db.View(func(tx *buntdb.Tx) error {
		s, err := tx.Get(u.Username)
		if err != nil {
			return err
		}
		w.Header().Set("content-type", "application/json")
		defer fmt.Fprintf(w, s)
		return nil
	}); err != nil {
		internalServerError(w, err)
		return
	}
}

func (a *apiHandler) handleUsers(w http.ResponseWriter, r *http.Request) {
	var head string
	head, r.URL.Path = shiftPath(r.URL.Path)
	p := fmt.Sprintf("%s%s", r.URL.Path, r.Method)

	switch p {
	case "/POST":
		a.handlePostUser(w, r)
	case "/claimsPOST":
		a.handlePostClaims(head).ServeHTTP(w, r)
	case "/claimsDELETE":
		a.handleDeleteClaims(head).ServeHTTP(w, r)
	default:
		notFoundError(w)
		return
	}
}

func (a *apiHandler) handlePostUser(w http.ResponseWriter, r *http.Request) {
	decoder := json.NewDecoder(r.Body)
	var u user
	if err := decoder.Decode(&u); err != nil {
		internalServerError(w, err)
		return
	}

	b, err := json.Marshal(u)
	if err != nil {
		internalServerError(w, err)
		return
	}

	if err := a.db.Update(func(tx *buntdb.Tx) error {
		_, _, err := tx.Set(u.Username, string(b), nil)
		if err != nil {
			return err
		}
		return nil
	}); err != nil {
		internalServerError(w, err)
		return
	}
}

func (a *apiHandler) handlePostClaims(id string) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprintf(w, id)
	}
}

func (a *apiHandler) handleDeleteClaims(id string) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {

	}
}
