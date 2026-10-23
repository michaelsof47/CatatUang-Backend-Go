package models

type BalanceInputOrUpdate struct {
	BalanceAmount string `json:"balance_amount" form:"balance_amount"`
}

type BalanceResult struct {
	Id            string `json:"id"`
	BalanceAmount string `json:"balance_amount"`
}
