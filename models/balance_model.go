package models

type BalanceInputOrUpdate struct {
	BalanceAmount int64 `json:"balance_amount" form:"balance_amount"`
}

type BalanceResult struct {
	Id            string `json:"id"`
	BalanceAmount int64  `json:"balance_amount"`
}
