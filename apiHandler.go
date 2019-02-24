package main

import (
	"encoding/json"
	"errors"
	"fmt"
	"net/http"

	"github.com/tidwall/buntdb"
	"golang.org/x/crypto/bcrypt"
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

type dbUser struct {
	Username string                 `json:"username"`
	Email    string                 `json:"email"`
	Password string                 `json:"password"`
	Claims   map[string]interface{} `json:"claims"`
}

func newDbUser(username string, email string, password string, claims map[string]interface{}) *dbUser {
	return &dbUser{
		Username: username,
		Email:    email,
		Password: password,
		Claims:   claims,
	}
}

func hashAndSalt(password string) (string, error) {
	hash, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.MinCost)
	if err != nil {
		return "", err
	}
	return string(hash), nil
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
	var u newUserRequest
	if err := decoder.Decode(&u); err != nil {
		internalServerError(w, err)
		return
	}

	encryptedPassword, err := hashAndSalt(u.Password)
	if err != nil {
		internalServerError(w, err)
		return
	}

	dbUser := newDbUser(u.Username, u.Email, encryptedPassword, u.Claims)
	b, err := json.Marshal(dbUser)
	if err != nil {
		internalServerError(w, err)
		return
	}

	if err := a.db.Update(func(tx *buntdb.Tx) error {
		if _, err := tx.Get(dbUser.Username); err == nil {
			return errors.New("username unavailable")
		}

		_, _, err = tx.Set(dbUser.Username, string(b), nil)
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
