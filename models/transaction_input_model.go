package models

type CreateCategory struct {
	CategoryName string `json:"category_name" form:"cat_name"`
	CategoryDesc string `json:"category_desc" form:"cat_desc"`
}
