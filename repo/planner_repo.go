package repo

import (
	"catatuangbackend/models"
	"database/sql"
)

type PlannerInterface interface {
	CreatePlannerBook(input models.PlannerBookModels, id string) error
	GetPlannerCount(id string) (int, error)
	GetPlannerList(id string) (*sql.Rows, error)
	RemovePlannerBook(bookId int, userId string) error
	UpdateBookData(data models.PlannerBookModels, userId string) error
	CreateItemPlanner(data models.DetailBookItemModels, userId string) error
	CheckIfItemPlannerWasAdded(bookId int, plannerName string) (bool, error)
	CheckIfPlannerBookWasOwnered(bookId int, userId string) (bool, error)
	GetDetailItemCount(bookId int) (int, error)
	GetDetailItemList(bookId int) (*sql.Rows, error)
}

type PlannerRepo struct {
	db *sql.DB
}

func PlannerRepository(db *sql.DB) PlannerInterface {
	return &PlannerRepo{db: db}
}

func (h *PlannerRepo) CreatePlannerBook(input models.PlannerBookModels, id string) error {
	query := `INSERT INTO "PlannerBook" (name, target_start, target_end, target_amount, user_id) VALUES ($1, $2, $3, $4, $5)`
	_, err := h.db.Exec(query, input.BookName, input.TargetStartDate, input.TargetEndDate, input.TargetAmount, id)
	return err
}

func (h *PlannerRepo) GetPlannerCount(id string) (int, error) {
	var totalItem int
	query := `SELECT COUNT(*) FROM "PlannerBook" WHERE user_id = $1`
	err := h.db.QueryRow(query, id).Scan(&totalItem)
	return totalItem, err
}

func (h *PlannerRepo) GetPlannerList(id string) (*sql.Rows, error) {
	query := `SELECT id, name, target_start, target_end, target_amount FROM "PlannerBook" WHERE user_id = $1`
	rows, err := h.db.Query(query, id)
	return rows, err
}

func (h *PlannerRepo) RemovePlannerBook(bookId int, userId string) error {
	query := `DELETE FROM "PlannerBook" WHERE id = $1 AND user_id = $2`
	_, err := h.db.Exec(query, bookId, userId)
	return err
}

func (h *PlannerRepo) UpdateBookData(data models.PlannerBookModels, userId string) error {
	query := `UPDATE "PlannerBook" SET name = $1, target_start = $2, target_end = $3, target_amount = $4, 
		updated_at = NOW() WHERE id = $5 AND user_id = $6`
	_, err := h.db.Exec(query, data.BookName, data.TargetStartDate, data.TargetEndDate, data.TargetAmount, data.BookId, userId)
	return err
}

func (h *PlannerRepo) CheckIfItemPlannerWasAdded(bookId int, plannerName string) (bool, error) {
	var exists bool
	query := `SELECT EXISTS(SELECT 1 FROM "DetailPlannerBook" WHERE planner_book_id = $1 AND name = $2)`
	err := h.db.QueryRow(query, bookId, plannerName).Scan(&exists)
	return exists, err
}

func (h *PlannerRepo) CheckIfPlannerBookWasOwnered(bookId int, userId string) (bool, error) {
	var isOwner bool
	query := `SELECT EXISTS(SELECT 1 FROM "PlannerBook" WHERE id = $1 AND user_id = $2)`
	err := h.db.QueryRow(query, bookId, userId).Scan(&isOwner)
	return isOwner, err
}

func (h *PlannerRepo) CreateItemPlanner(data models.DetailBookItemModels, userId string) error {
	query := `INSERT INTO "DetailPlannerBook" (name, total_amount, planner_book_id) VALUES ($1, $2, $3)`
	_, err := h.db.Exec(query, data.PlannerItemName, data.ItemTotalAmount, data.PlannerBookId.BookId)
	return err
}

func (h *PlannerRepo) GetDetailItemCount(bookId int) (int, error) {
	var count int
	query := `SELECT COUNT(*) FROM "DetailPlannerBook" WHERE planner_book_id = $1`
	err := h.db.QueryRow(query, bookId).Scan(&count)
	return count, err
}

func (h *PlannerRepo) GetDetailItemList(bookId int) (*sql.Rows, error) {
	query := `SELECT id, name, total_amount FROM "DetailPlannerBook" WHERE planner_book_id = $1`
	rows, err := h.db.Query(query, bookId)
	return rows, err
}
