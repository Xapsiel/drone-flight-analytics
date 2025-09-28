package model

type User struct {
	ID       string   `json:"id"`
	Email    string   `json:"email"`
	Token    string   `json:"token"`
	Name     string   `json:"name"`
	Username string   `json:"username"`
	Roles    []string `json:"roles"`
}
