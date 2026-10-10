package models

type BalanceInputOrUpdate struct {
	BalanceAmount string `json:"balance_amount" form:"balance_amount"`
}
