package domain

type Balance struct {
	Current   int `json:"current"`
	Withdrawn int `json:"withdrawn"`
}
