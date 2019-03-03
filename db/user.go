package db

type User struct {
	Username string                 `json:"username"`
	Email    string                 `json:"email"`
	Password string                 `json:"password"`
	Claims   map[string]interface{} `json:"claims"`
}

func newDbUser(username string, email string, password string, claims map[string]interface{}) *User {
	return &User{
		Username: username,
		Email:    email,
		Password: password,
		Claims:   claims,
	}
}
