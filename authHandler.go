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
		return
	}

	switch authType {
	case "Basic":
		a.handleBasicAuth(w, r)
		return
	case "Bearer":
		return
	}
}

func authType(r *http.Request) (string, error) {
	match := firstCaseInsensitivePrefixMatch(r.Header.Get("Authorization"), "Basic ", "Bearer ")
	if len(match) < 1 {
		return "", errors.New("only basic and bearer auth are supported")
	}
	return match[:len(match)-1], nil
}

func firstCaseInsensitivePrefixMatch(value string, prefixes ...string) string {
	for _, prefix := range prefixes {
		if len(value) >= len(prefix) && strings.EqualFold(value[:len(prefix)], prefix) {
			return prefix
		}
	}
	return ""
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
