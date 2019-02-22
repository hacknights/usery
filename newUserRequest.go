package main

import (
	"errors"
	"net/mail"
)

type newUserRequest struct {
	Username string
	Email    string
	Password string
	Claims   map[string]interface{}
}

func (r *newUserRequest) validateRequest() error {
	if err := r.validateUsername(); err != nil {
		return err
	}

	if err := r.validateEmail(); err != nil {
		return err
	}

	return nil
}

func (r *newUserRequest) validateUsername() error {
	if length := len(r.Username); length < 3 || 40 < length {
		return errors.New("username must be greater than 3 and less than 40")
	}

	return nil
}

func (r *newUserRequest) validateEmail() error {
	if _, err := mail.ParseAddress(r.Email); err != nil {
		return errors.New("invalid email address")
	}

	return nil
}
