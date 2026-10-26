package models

type CreateCategory struct {
	CategoryName string `json:"category_name" form:"cat_name"`
	CategoryDesc string `json:"category_desc" form:"cat_desc"`
}

type CategoryIdItem struct {
	CategoryId int `json:"category_id" form:"cat_id"`
}

type CategoryItem struct {
	CategoryIdItem
	CreateCategory
	CategoryUrlPhoto string `json:"category_url_photo"`
}

type TransactionInput struct {
	TransactionName       string `json:"trans_name" form:"trans_name"`
	TransactionDate       string `json:"trans_date" form:"trans_date"`
	TransactionAmount     int    `json:"trans_amount" form:"trans_amount"`
	MerchantName          string `json:"merch_name" form:"merch_name"`
	TransactionPrice      int    `json:"trans_price" form:"trans_price"`
	TransactionTotalPrice int64  `json:"trans_total_price" form:"trans_total_price"`
	CategoryIdItem
}

type TransactionResult struct {
	TransactionId         int    `json:"trans_id"`
	TransactionName       string `json:"trans_name"`
	TransactionDate       string `json:"trans_date"`
	TransactionAmount     int    `json:"trans_amount"`
	MerchantName          string `json:"merch_name"`
	TransactionPrice      int    `json:"trans_price"`
	TransactionTotalPrice int    `json:"trans_total_price"`
	CategoryIdItem
}
