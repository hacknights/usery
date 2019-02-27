package main

import "errors"

type newAppRequest struct {
	Appname string
}

func (r *newAppRequest) validateRequest() error {
	if length := len(r.Appname); length < 3 || 40 < length {
		return errors.New("appname must be greater than 3 and less than 40")
	}
	return nil
}
