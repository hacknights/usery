package db

import "encoding/json"

type App struct {
	Id      string `json:"id"`
	Appname string `json:"appname"`
	Users   []User `json:"users"`
}

func NewApp(appname string) *App {
	return &App{
		Appname: appname,
	}
}

func (a *App) Identify() string {
	return a.Id
}

func (a *App) identityAssign(id string) {
	a.Id = id
}

func (a *App) serialize() (string, error) {
	b, err := json.Marshal(a)
	if err != nil {
		return "", err
	}
	return string(b), nil
}

func (a *App) deserialize(s string) error {
	if err := json.Unmarshal([]byte(s), a); err != nil {
		return err
	}
	return nil
}
