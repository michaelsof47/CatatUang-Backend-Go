package models

type CreateCategory struct {
	CategoryName string `json:"category_name" form:"cat_name"`
	CategoryDesc string `json:"category_desc" form:"cat_desc"`
}

type CategoryItem struct {
	CategoryId int `json:"category_id"`
	CreateCategory
	CategoryUrlPhoto string `json:"category_url_photo"`
}
