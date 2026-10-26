package handlers

import (
	"catatuangbackend/models"
	"catatuangbackend/repo"
	"net/http"

	"github.com/gin-gonic/gin"
)

type BalancesHandler struct {
	BalanceRepo repo.BalanceInterface
}

func BalancesHandlerFunc(balanceRepo repo.BalanceInterface) *BalancesHandler {
	return &BalancesHandler{BalanceRepo: balanceRepo}
}

func (h *BalancesHandler) CreateOrUpdateBalancesAmount(c *gin.Context) {
	id := c.GetString("id")
	var input models.BalanceInputOrUpdate

	if err := c.ShouldBind(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if err := h.BalanceRepo.CreateOrUpdateBalance(input.BalanceAmount, id); err != nil {
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

	if err := h.BalanceRepo.AddMoreBalance(input.BalanceAmount, id); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Deposit Akun anda berhasil ditambahkan"})
}

func (h *BalancesHandler) GetBalanceAmount(c *gin.Context) {
	id := c.GetString("id")

	output, err := h.BalanceRepo.GetBalanceAmount(id)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, output)
}
