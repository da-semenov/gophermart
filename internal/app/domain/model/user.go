package model

type User struct {
	ID    int
	Login string `json:"login"`
	Pass  string `json:"password"`
}
