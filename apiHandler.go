package main

import (
	"encoding/base64"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"strings"

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

func (u *dbUser) validatePassword(password string) bool {
	b := []byte(password)
	if err := bcrypt.CompareHashAndPassword([]byte(u.Password), b); err != nil {
		return false
	}
	return true
}

func (a *apiHandler) handleAuthenticate(w http.ResponseWriter, r *http.Request) {
	u := authenticateUserRequest{}
	if err := decodeBasicAuth(&u, r.Header.Get("Authorization")); err != nil {
		unauthorizedError(w)
		return
	}

	if err := a.db.View(func(tx *buntdb.Tx) error {
		s, err := tx.Get(u.username)
		if err != nil {
			return err
		}

		d := dbUser{}
		if err := json.Unmarshal([]byte(s), &d); err != nil {
			return err
		}

		if valid := d.validatePassword(u.password); !valid {
			return errors.New("unauthorized")
		}

		w.Header().Set("content-type", "application/json")
		response := authenticatedUserResponse{
			Claims: d.Claims,
		}

		b, err := json.Marshal(response)
		if err != nil {
			return err
		}

		defer fmt.Fprintf(w, "%s", b)
		return nil
	}); err != nil {
		unauthorizedError(w)
		return
	}
}

type authenticateUserRequest struct {
	username string
	password string
}

type authenticatedUserResponse struct {
	Claims map[string]interface{} `json:"claims"`
}

func decodeBasicAuth(u *authenticateUserRequest, encodedCredentials string) error {
	auth := strings.SplitN(encodedCredentials, " ", 2)
	if auth[0] != "Basic" {
		return errors.New("must be basic auth")
	}

	b, err := base64.StdEncoding.DecodeString(auth[1])
	if err != nil {
		return err
	}

	decoded := strings.SplitN(string(b), ":", 2)
	u.username = decoded[0]
	u.password = decoded[1]

	return nil
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
