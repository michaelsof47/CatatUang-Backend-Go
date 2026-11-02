package models

type PlannerBookModels struct {
	Id              int    `json:"book_id" form:"book_id"`
	BookName        string `json:"book_name" form:"book_name"`
	TargetStartDate string `json:"target_start_date" form:"target_start_date"`
	TargetEndDate   string `json:"target_end_date" form:"target_end_date"`
	TargetAmount    int64  `json:"target_amount" form:"target_amount"`
}
