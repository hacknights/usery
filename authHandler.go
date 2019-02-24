package main

import (
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"strings"

	"github.com/tidwall/buntdb"
	"golang.org/x/crypto/bcrypt"
)

type authenticatedUserResponse struct {
	Claims map[string]interface{} `json:"claims"`
}

func newAuthenticatedUserResponse() *authenticatedUserResponse {
	return &authenticatedUserResponse{}
}

func (a *apiHandler) handleAuthenticate(w http.ResponseWriter, r *http.Request) {
	authType, err := authType(r)
	if err != nil {
		unauthorizedError(w, err)
	}

	switch authType {
	case "Basic":
		a.handleBasicAuth(w, r)
		return
	}
}

func authType(r *http.Request) (string, error) {
	authType := strings.SplitN(r.Header.Get("Authorization"), " ", 2)[0]
	switch authType {
	case "Basic":
		return authType, nil
	case "Bearer":
		return authType, nil
	default:
		return "", errors.New("only basic and bearer auth are supported")
	}
}

func (a *apiHandler) handleBasicAuth(w http.ResponseWriter, r *http.Request) {
	user, pass, ok := r.BasicAuth()
	if !ok {
		unauthorizedError(w, errors.New("invalid basic auth value"))
		return
	}

	response := newAuthenticatedUserResponse()
	if err := a.getClaims(user, pass, response); err != nil {
		unauthorizedError(w, err)
		return
	}

	b, err := json.Marshal(response)
	if err != nil {
		internalServerError(w, err)
		return
	}

	w.Header().Set("content-type", "application/json")
	fmt.Fprintf(w, "%s", b)
}

func (a *apiHandler) getClaims(username, password string, r *authenticatedUserResponse) error {
	const invalidCredentials = "invalid credentials"
	if err := a.db.View(func(tx *buntdb.Tx) error {
		s, err := tx.Get(username)
		if err != nil {
			return errors.New(invalidCredentials)
		}

		d := dbUser{}
		if err := json.Unmarshal([]byte(s), &d); err != nil {
			return err
		}

		if valid := d.validatePassword(password); !valid {
			return errors.New(invalidCredentials)
		}

		r.Claims = d.Claims

		return nil
	}); err != nil {
		return err
	}
	return nil
}

func (u *dbUser) validatePassword(password string) bool {
	b := []byte(password)
	if err := bcrypt.CompareHashAndPassword([]byte(u.Password), b); err != nil {
		return false
	}
	return true
}
