package repo

import (
	"catatuangbackend/models"
	"database/sql"
)

type TransactionInterface interface {
	IsCategoryDataAvailable(categoryName *string) (bool, error)
	InsertNewCategory(categoryName, categoryDesc *string, savePath, id string) error
	GetCategoryCount(id string) (int, error)
	GetCategoriesList(id string, limit, offset int) (*sql.Rows, error)
	GetCategoryPhoto(catId *int, userId string) (string, error)
	CreateTransaction(diffBalanceAmount int64, id string, input models.TransactionInput, balanceId string) error
	GetTransactionCount(id string) (int, error)
	GetTransactionList(id string, limit, offset int) (*sql.Rows, error)
	RemoveCategoryItem(catId *int, id string) error
}

type TransactionRepo struct {
	db *sql.DB
}

func TransactionRepository(db *sql.DB) TransactionInterface {
	return &TransactionRepo{db: db}
}

func (h *TransactionRepo) IsCategoryDataAvailable(categoryName *string) (bool, error) {
	var exists bool
	isAccAvailableQuery := `SELECT EXISTS(SELECT 1 FROM "Category" WHERE LOWER(name) = LOWER($1))`
	err := h.db.QueryRow(isAccAvailableQuery, categoryName).Scan(&exists)
	return exists, err
}

func (h *TransactionRepo) InsertNewCategory(categoryName, categoryDesc *string, savePath, id string) error {
	query := `INSERT INTO "Category" (name, description, url_image, user_id) VALUES ($1, $2, $3, $4)`
	_, err := h.db.Exec(query, categoryName, categoryDesc, savePath, id)
	return err
}

func (h *TransactionRepo) GetCategoryCount(id string) (int, error) {
	var totalItem int
	countQuery := `SELECT COUNT(*) FROM "Category" WHERE user_id = $1`
	err := h.db.QueryRow(countQuery, id).Scan(&totalItem)
	return totalItem, err
}

func (h *TransactionRepo) GetCategoriesList(id string, limit, offset int) (*sql.Rows, error) {
	query := `SELECT id, name, description, url_image FROM "Category" WHERE user_id = $1 ORDER BY id DESC LIMIT $2 OFFSET $3`
	rows, err := h.db.Query(query, id, limit, offset)
	return rows, err
}

func (h *TransactionRepo) GetCategoryPhoto(catId *int, userId string) (string, error) {
	var urlPhoto string
	query := `SELECT url_image FROM "Category" WHERE user_id = $1 AND id = $2`
	err := h.db.QueryRow(query, userId, catId).Scan(&urlPhoto)
	return urlPhoto, err
}

func (h *TransactionRepo) CreateTransaction(diffBalanceAmount int64, id string, input models.TransactionInput, balanceId string) error {
	rollback, err := h.db.Begin()
	if err != nil {
		return err
	}
	defer rollback.Rollback()

	updateBalanceQuery := `UPDATE "Balance" SET balance_amount = $1 WHERE user_id = $2`
	if _, err := rollback.Exec(updateBalanceQuery, diffBalanceAmount, id); err != nil {
		return err
	}

	query := `INSERT INTO "Transaction" (name, transaction_date, amount, merchant_name, item_price, total_price, category_id, balances_id, user_id) 
		VALUES($1, $2, $3, $4, $5, $6, $7, $8, $9)`
	if _, err := rollback.Exec(query, input.TransactionName, input.TransactionDate,
		input.TransactionAmount, input.MerchantName, input.TransactionPrice, input.TransactionTotalPrice,
		input.CategoryId, balanceId, id); err != nil {
		return err
	}

	return rollback.Commit()
}

func (h *TransactionRepo) GetTransactionCount(id string) (int, error) {
	var totalItem int
	findTotalItem := `SELECT COUNT(*) FROM "Transaction" WHERE user_id = $1`
	err := h.db.QueryRow(findTotalItem, id).Scan(&totalItem)
	return totalItem, err
}

func (h *TransactionRepo) GetTransactionList(id string, limit, offset int) (*sql.Rows, error) {
	query := `SELECT id, name, transaction_date, amount, merchant_name, item_price,
		total_price, category_id FROM "Transaction" WHERE user_id = $1 ORDER BY id DESC LIMIT $2 OFFSET $3`
	rows, err := h.db.Query(query, id, limit, offset)
	return rows, err
}

func (h *TransactionRepo) RemoveCategoryItem(catId *int, id string) error {
	deleteQuery := `DELETE FROM "Category" WHERE id = $1 AND user_id = $2`
	_, err := h.db.Exec(deleteQuery, catId, id)
	return err
}
