package main

// import (
// 	"encoding/json"
// 	"errors"
// 	"fmt"
// 	"net/http"
// 	"strings"
// 	"usery/db"

// 	"github.com/hacknights/negotiator"

// 	"github.com/tidwall/buntdb"
// 	"golang.org/x/crypto/bcrypt"
// )

// type authHandler struct {
// 	negotiator negotiator.Factory
// 	db         *db.DB
// }

// func newAuthHandler(n negotiator.Factory, db *db.DB) *authHandler {
// 	return &authHandler{
// 		negotiator: n,
// 		db:         db,
// 	}
// }

// type authenticatedUserResponse struct {
// 	Claims map[string]interface{} `json:"claims"`
// }

// func newAuthenticatedUserResponse() *authenticatedUserResponse {
// 	return &authenticatedUserResponse{}
// }

// func (h *authHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
// 	authType, err := authType(r)
// 	if err != nil {
// 		h.negotiator(w, r).InternalServerError(err)
// 		return
// 	}

// 	switch authType {
// 	case "Basic":
// 		h.handleBasicAuth(w, r)
// 		return
// 	case "Bearer":
// 		return
// 	}
// }

// func authType(r *http.Request) (string, error) {
// 	match := firstCaseInsensitivePrefixMatch(r.Header.Get("Authorization"), "Basic ", "Bearer ")
// 	if len(match) < 1 {
// 		return "", errors.New("only basic and bearer auth are supported")
// 	}
// 	return match[:len(match)-1], nil
// }

// func firstCaseInsensitivePrefixMatch(value string, prefixes ...string) string {
// 	for _, prefix := range prefixes {
// 		if len(value) >= len(prefix) && strings.EqualFold(value[:len(prefix)], prefix) {
// 			return prefix
// 		}
// 	}
// 	return ""
// }

// func (h *authHandler) handleBasicAuth(w http.ResponseWriter, r *http.Request) {
// 	user, pass, ok := r.BasicAuth()
// 	if !ok {
// 		h.negotiator(w, r).Unauthorized("invalid basic auth value")
// 		return
// 	}

// 	response := newAuthenticatedUserResponse()
// 	if err := h.getClaims(user, pass, response); err != nil {
// 		h.negotiator(w, r).UnauthorizedError(err)
// 		return
// 	}

// 	b, err := json.Marshal(response)
// 	if err != nil {
// 		h.negotiator(w, r).UnauthorizedError(err)
// 		return
// 	}

// 	w.Header().Set("content-type", "application/json")
// 	fmt.Fprintf(w, "%s", b)
// }

// func (h *authHandler) getClaims(username, password string, r *authenticatedUserResponse) error {
// 	const invalidCredentials = "invalid credentials"
// 	if err := h.db.View(func(tx *buntdb.Tx) error {
// 		s, err := tx.Get(username)
// 		if err != nil {
// 			return errors.New(invalidCredentials)
// 		}

// 		d := dbUser{}
// 		if err := json.Unmarshal([]byte(s), &d); err != nil {
// 			return err
// 		}

// 		if valid := d.validatePassword(password); !valid {
// 			return errors.New(invalidCredentials)
// 		}

// 		r.Claims = d.Claims

// 		return nil
// 	}); err != nil {
// 		return err
// 	}
// 	return nil
// }

// func (u *dbUser) validatePassword(password string) bool {
// 	b := []byte(password)
// 	if err := bcrypt.CompareHashAndPassword([]byte(u.Password), b); err != nil {
// 		return false
// 	}
// 	return true
// }
