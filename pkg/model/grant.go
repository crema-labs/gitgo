package model

type Grant struct {
	GrantID       string             `json:"grantid"`
	GrantAmount   string             `json:"grant_amount"`
	Status        string             `json:"status"`
	Contributions map[string]float64 `json:"contribution"`
}
