package model

type Balance struct {
	Current   int `json:"current"`
	Withdrawn int `json:"withdrawn"`
}
