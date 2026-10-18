package handlers

import (
	"catatuangbackend/models"
	"database/sql"
	"net/http"

	"github.com/gin-gonic/gin"
)

type BalancesHandler struct {
	DB *sql.DB
}

func BalancesHandlerFunc(db *sql.DB) *BalancesHandler {
	return &BalancesHandler{DB: db}
}

func (h *BalancesHandler) CreateOrUpdateBalancesAmount(c *gin.Context) {
	id := c.GetString("id")
	var input models.BalanceInputOrUpdate

	if err := c.ShouldBind(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	query := `INSERT INTO "Balance" (user_id, balance_amount)
		VALUES ($1, $2)
		ON CONFLICT (user_id) DO UPDATE
		SET balance_amount = $2,
		updated_at = NOW()`
	if _, err := h.DB.Exec(query, id, input.BalanceAmount); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Deposit Akun Berhasil Ditambahkan"})
}

func (h *BalancesHandler) AddBalanceAmount(c *gin.Context) {
	id := c.GetString("id")
	var input models.BalanceInputOrUpdate

	if err := c.ShouldBind(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	query := `UPDATE "Balance" SET balance_amount = balance_amount + $1 WHERE user_id = $2`
	if _, err := h.DB.Exec(query, input.BalanceAmount, id); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Deposit Akun anda berhasil ditambahkan"})
}

func (h *BalancesHandler) GetBalanceAmount(c *gin.Context) {
	id := c.GetString("id")

	var output models.BalanceResult
	query := `SELECT user_id, balance_amount FROM "Balance" WHERE user_id = $1`
	if err := h.DB.QueryRow(query, id).Scan(&output.User_Id, &output.BalanceAmount); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, output)
}
