package handlers

import (
	"catatuangbackend/models"
	"catatuangbackend/repo"
	"catatuangbackend/utils"
	"net/http"
	"strconv"
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
		if err := rows.Scan(&item.BookId, &item.BookName, &item.TargetStartDate, &item.TargetEndDate, &item.TargetAmount); err != nil {
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

	if err := h.plannerRepo.RemovePlannerBook(input.BookId, id); err != nil {
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

func (h *PlannerHandler) CreateItemPlanner(c *gin.Context) {
	id := c.GetString("id")
	var input models.DetailBookItemModels

	if err := c.ShouldBind(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	//Check if book was an owner or not
	isOwner, err := h.plannerRepo.CheckIfPlannerBookWasOwnered(input.BookId, id)
	if !isOwner {
		c.JSON(http.StatusForbidden, gin.H{"error": "Anda tidak dapat membuat item buku karena bukan punya anda. Silahkan menghubungi admin terlebih dahulu."})
		return
	} else if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	//Check if book name was similar in database. if similar, user get error
	exists, err := h.plannerRepo.CheckIfItemPlannerWasAdded(input.PlannerBookId.BookId, input.PlannerItemName)
	if exists {
		c.JSON(http.StatusConflict, gin.H{"error": "Terjadi duplikasi item. Mohon gunakan item lainnya."})
		return
	} else if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	if err := h.plannerRepo.CreateItemPlanner(input, id); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Item Buku telah berhasil dibuat"})
}

func (h *PlannerHandler) GetDetailItemList(c *gin.Context) {
	id := c.GetString("id")

	bookId := c.Param("id")
	_bookId, _ := strconv.Atoi(bookId)

	_page := c.DefaultQuery("page", "1")
	_limit := c.DefaultQuery("limit", "10")

	params := utils.ParsePagination(_page, _limit)

	//Check if book was an owner or not
	isOwner, err := h.plannerRepo.CheckIfPlannerBookWasOwnered(_bookId, id)
	if !isOwner {
		c.JSON(http.StatusForbidden, gin.H{"error": "Anda tidak dapat membuat item buku karena bukan punya anda. Silahkan menghubungi admin terlebih dahulu."})
		return
	} else if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	totalItem, err := h.plannerRepo.GetDetailItemCount(_bookId)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	rows, err := h.plannerRepo.GetDetailItemList(_bookId)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	itemList := make([]models.DetailBookItemModels, 0)
	for rows.Next() {
		var item models.DetailBookItemModels
		if err := rows.Scan(&item.ItemPlannerId, &item.PlannerItemName, &item.ItemTotalAmount); err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}
		itemList = append(itemList, item)
	}

	if err := rows.Err(); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"pagination": utils.BuildPaginationData(params, totalItem),
		"data":       itemList,
	})
}
