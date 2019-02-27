package requests

import "errors"

type NewApp struct {
	Appname string
}

func (r *NewApp) Validate() error {
	if length := len(r.Appname); length < 3 || 40 < length {
		return errors.New("appname must be greater than 3 and less than 40")
	}
	return nil
}
