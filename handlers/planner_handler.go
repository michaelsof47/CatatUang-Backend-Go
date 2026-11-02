package handlers

import (
	"catatuangbackend/models"
	"catatuangbackend/repo"
	"catatuangbackend/utils"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
)

type PlannerHandler struct {
	plannerRepo repo.PlannerInterface
}

func PlannerHandlerFunc(plannerRepo repo.PlannerInterface) *PlannerHandler {
	return &PlannerHandler{plannerRepo: plannerRepo}
}

// This function work for validate between target_start and target_end
func ValidatePlannerDates(startDate, endDate string) string {

	_endDate, err := time.Parse("2006-01-02", endDate)
	if err != nil {
		return "Format target_end_date tidak valid"
	}

	_startDate, err := time.Parse("2006-01-02", startDate)
	if err != nil {
		return "Format target_start_date tidak valid"
	}

	now := time.Now()
	today := time.Date(now.Year(), now.Month(), now.Day(), 0, 0, 0, 0, now.Location())

	if _endDate.Before(today) {
		return "Target tanggal selesai tidak boleh kurang dari hari ini"
	}

	if !_startDate.IsZero() && _endDate.Before(_startDate) {
		return "Target tanggal selesai harus lebih besar dari tanggal mulai"
	}

	return ""
}

func (h *PlannerHandler) CreatePlannerBook(c *gin.Context) {
	id := c.GetString("id")
	var input models.PlannerBookModels

	if err := c.ShouldBind(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if err := ValidatePlannerDates(input.TargetStartDate, input.TargetEndDate); err != "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": err})
		return
	}

	if err := h.plannerRepo.CreatePlannerBook(input, id); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Buku " + input.BookName + " sudah berhasil dibuat"})
}

func (h *PlannerHandler) GetPlannerBookList(c *gin.Context) {
	id := c.GetString("id")

	_page := c.DefaultQuery("page", "1")
	_limit := c.DefaultQuery("limit", "10")

	params := utils.ParsePagination(_page, _limit)

	totalItem, err := h.plannerRepo.GetPlannerCount(id)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	rows, err := h.plannerRepo.GetPlannerList(id)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	planners := make([]models.PlannerBookModels, 0)
	for rows.Next() {
		var item models.PlannerBookModels
		if err := rows.Scan(&item.Id, &item.BookName, &item.TargetStartDate, &item.TargetEndDate, &item.TargetAmount); err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}
		planners = append(planners, item)
	}

	if err := rows.Err(); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"pagination": utils.BuildPaginationData(params, totalItem),
		"data":       planners,
	})
}

func (h *PlannerHandler) RemoveBookItem(c *gin.Context) {
	id := c.GetString("id")
	var input models.PlannerBookModels

	if err := c.ShouldBind(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if err := h.plannerRepo.RemovePlannerBook(input.Id, id); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Buku telah berhasil dihapus"})
}

func (h *PlannerHandler) UpdatePlannerBook(c *gin.Context) {
	id := c.GetString("id")
	var input models.PlannerBookModels

	if err := c.ShouldBind(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if err := ValidatePlannerDates(input.TargetStartDate, input.TargetEndDate); err != "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": err})
		return
	}

	if err := h.plannerRepo.UpdateBookData(input, id); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Data Buku berhasil diubah."})
}
