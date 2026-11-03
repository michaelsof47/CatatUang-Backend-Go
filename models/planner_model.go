package models

type PlannerBookId struct {
	BookId int `json:"book_id" form:"book_id"`
}

type PlannerBookModels struct {
	PlannerBookId
	BookName        string `json:"book_name" form:"book_name"`
	TargetStartDate string `json:"target_start_date" form:"target_start_date"`
	TargetEndDate   string `json:"target_end_date" form:"target_end_date"`
	TargetAmount    int64  `json:"target_amount" form:"target_amount"`
}

type DetailBookItemModels struct {
	ItemPlannerId   int    `json:"item_id"`
	PlannerItemName string `json:"item_planner_name" form:"item_planner_name"`
	ItemTotalAmount int    `json:"item_total_amount" form:"item_total_amount"`
	PlannerBookId
}
