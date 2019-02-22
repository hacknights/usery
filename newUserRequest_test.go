package main

import "testing"

const userNameError = "username must be greater than 3 and less than 40"
const emailError = "invalid email address"
const passwordError = "password must contain 1 uppercase, 1 special character, 1 number and be 6 characters long"

func TestValidateUsername_Empty(t *testing.T) {
	r := newUserRequest{Username: ""}
	if err := r.validateUsername(); err.Error() != userNameError {
		errorf(t, userNameError, err.Error())
	}
}

func TestValidateUsername_2Characters(t *testing.T) {
	r := newUserRequest{Username: "jo"}
	if err := r.validateUsername(); err.Error() != userNameError {
		errorf(t, userNameError, err.Error())
	}
}

func TestValidateUsername_41Characters(t *testing.T) {
	r := newUserRequest{Username: "Loremipsumdolorametphoto/boothseitanshabb"}
	if err := r.validateUsername(); err.Error() != userNameError {
		errorf(t, userNameError, err.Error())
	}
}

func TestValidateUsername(t *testing.T) {
	r := newUserRequest{Username: "john"}
	if err := r.validateUsername(); err != nil {
		errorf(t, nil, err.Error())
	}
}

func TestValidateEmail_Empty(t *testing.T) {
	r := newUserRequest{Email: ""}
	if err := r.validateEmail(); err.Error() != emailError {
		errorf(t, emailError, err.Error())
	}
}

func TestValidateEmail_NoAt(t *testing.T) {
	r := newUserRequest{Email: "test"}
	if err := r.validateEmail(); err.Error() != emailError {
		errorf(t, emailError, err.Error())
	}
}

func TestValidateEmail(t *testing.T) {
	r := newUserRequest{Email: "test@test.com"}
	if err := r.validateEmail(); err != nil {
		errorf(t, emailError, err.Error())
	}
}

func errorf(t *testing.T, expected, actual interface{}) {
	t.Errorf("expected: %s, actual: %s", expected, actual)
}
