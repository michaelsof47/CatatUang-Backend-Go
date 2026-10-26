package repo

import (
	"catatuangbackend/models"
	"database/sql"
)

type BalanceInterface interface {
	CreateOrUpdateBalance(balanceAmount int64, id string) error
	AddMoreBalance(balanceAmount int64, id string) error
	GetBalanceAmount(id string) (*models.BalanceResult, error)
}

type BalanceRepo struct {
	db *sql.DB
}

func BalanceRepository(db *sql.DB) BalanceInterface {
	return &BalanceRepo{db: db}
}

func (h *BalanceRepo) CreateOrUpdateBalance(balanceAmount int64, id string) error {
	query := `INSERT INTO "Balance" (user_id, balance_amount)
		VALUES ($1, $2)
		ON CONFLICT (user_id) DO UPDATE
		SET balance_amount = $2,
		updated_at = NOW()`
	_, err := h.db.Exec(query, id, balanceAmount)
	return err
}

func (h *BalanceRepo) AddMoreBalance(balanceAmount int64, id string) error {
	query := `UPDATE "Balance" SET balance_amount = balance_amount + $1, updated_at = NOW() WHERE user_id = $2`
	_, err := h.db.Exec(query, balanceAmount, id)
	return err
}

func (h *BalanceRepo) GetBalanceAmount(id string) (*models.BalanceResult, error) {
	var output models.BalanceResult
	query := `SELECT user_id, balance_amount FROM "Balance" WHERE user_id = $1`
	err := h.db.QueryRow(query, id).Scan(&output.Id, &output.BalanceAmount)
	return &output, err
}
