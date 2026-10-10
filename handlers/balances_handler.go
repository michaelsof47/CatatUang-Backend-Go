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
	_, err := h.DB.Exec(query, id, input.BalanceAmount)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Deposit Akun Berhasil Ditambahkan"})
}
