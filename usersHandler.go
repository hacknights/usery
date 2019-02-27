package main

import (
	"encoding/json"
	"errors"
	"fmt"
	"net/http"

	"github.com/hacknights/negotiator"

	"github.com/tidwall/buntdb"
	"golang.org/x/crypto/bcrypt"
)

type usersHandler struct {
	negotiator negotiator.Factory
	db         *buntdb.DB
}

func newUsersHandler(n negotiator.Factory, db *buntdb.DB) *usersHandler {
	return &usersHandler{
		negotiator: n,
		db:         db,
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

func (u *dbUser) incrementRev() {
	if val := u.Claims["rev"]; val == nil {
		u.Claims["rev"] = 1
	} else {
		rev := val.(int)
		u.Claims["rev"] = rev + 1
	}
}

func hashAndSalt(password string) (string, error) {
	hash, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.MinCost)
	if err != nil {
		return "", err
	}
	return string(hash), nil
}

func (h *usersHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	var head string
	head, r.URL.Path = shiftPath(r.URL.Path)
	p := fmt.Sprintf("%s%s", r.URL.Path, r.Method)

	switch p {
	case "/POST":
		h.handlePostUser(w, r)
	case "/claimsPOST":
		h.handlePostClaims(head).ServeHTTP(w, r)
	case "/claimsDELETE":
		h.handleDeleteClaims(head).ServeHTTP(w, r)
	default:
		h.negotiator(w, r).NotFound()
		return
	}
}

func (h *usersHandler) handlePostUser(w http.ResponseWriter, r *http.Request) {
	n := h.negotiator(w, r)
	decoder := json.NewDecoder(r.Body)
	var ur newUserRequest
	if err := decoder.Decode(&ur); err != nil {
		n.InternalServerError(err)
		return
	}

	if err := ur.validateRequest(); err != nil {
		n.BadRequestError(err)
		return
	}

	encryptedPassword, err := hashAndSalt(ur.Password)
	if err != nil {
		n.InternalServerError(err)
		return
	}

	dbUser := newDbUser(ur.Username, ur.Email, encryptedPassword, ur.Claims)
	dbUser.incrementRev()
	b, err := json.Marshal(dbUser)
	if err != nil {
		n.InternalServerError(err)
		return
	}

	const usernameUnavailableErr = "username unavailable"
	if err := h.db.Update(func(tx *buntdb.Tx) error {
		if _, err := tx.Get(dbUser.Username); err == nil {
			return errors.New(usernameUnavailableErr)
		}

		if _, _, err := tx.Set(dbUser.Username, string(b), nil); err != nil {
			return err
		}
		return nil
	}); err != nil {
		if err.Error() == usernameUnavailableErr {
			n.BadRequestError(err)
			return
		}

		n.InternalServerError(err)
		return
	}
}

func (a *usersHandler) handlePostClaims(id string) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprintf(w, id)
	}
}

func (a *usersHandler) handleDeleteClaims(id string) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {

	}
}
