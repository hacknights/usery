package responses

import "encoding/json"

type NewApp struct {
	ID      string `json:"id"`
	Appname string `json:"appname"`
}

func NewNewApp(id, appname string) *NewApp {
	return &NewApp{
		ID:      id,
		Appname: appname,
	}
}

func (r *NewApp) Serialize() (string, error) {
	b, err := json.Marshal(r)
	if err != nil {
		return "", err
	}
	return string(b), nil
}
