package requests

import (
	"errors"
	"fmt"
	"net/mail"
	"strings"
)

type NewUser struct {
	Username string
	Email    string
	Password string
	Claims   map[string]interface{}
}

func (r *NewUser) validateRequest() error {
	var sb strings.Builder
	if err := r.validateUsername(); err != nil {
		writeError(&sb, err)
	}

	if err := r.validateEmail(); err != nil {
		writeError(&sb, err)
	}

	if 0 < sb.Len() {
		return errors.New(sb.String())
	}
	return nil
}

func (r *NewUser) validateUsername() error {
	if length := len(r.Username); length < 3 || 40 < length {
		return errors.New("username must be greater than 3 and less than 40")
	}
	return nil
}

func (r *NewUser) validateEmail() error {
	if _, err := mail.ParseAddress(r.Email); err != nil {
		return errors.New("invalid email address")
	}
	return nil
}

func writeError(sb *strings.Builder, err error) {
	sb.WriteString(fmt.Sprintf("%s\n", err.Error()))
}
