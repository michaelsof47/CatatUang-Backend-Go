package handlers

import (
	"catatuangbackend/models"
	"catatuangbackend/repo"
	"catatuangbackend/utils"
	"database/sql"
	"log"
	"net/http"
	"os"
	"path/filepath"

	"github.com/gin-gonic/gin"
)

type TransactionHandler struct {
	transactionRepo repo.TransactionInterface
	balanceRepo     repo.BalanceInterface
}

func TransactionHandlerFunc(transactionRepo repo.TransactionInterface) *TransactionHandler {
	return &TransactionHandler{transactionRepo: transactionRepo}
}

func (h *TransactionHandler) CreateCategory(c *gin.Context) {
	var input models.CreateCategory
	id := c.GetString("id")

	if err := c.ShouldBind(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Check if category is available
	exists, err := h.transactionRepo.IsCategoryDataAvailable(&input.CategoryName)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	if exists {
		c.JSON(http.StatusConflict, gin.H{"error": "Nama kategori sudah ada"})
		return
	}

	file, err := c.FormFile("cat_url_image")
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Wajib upload foto kategori"})
		return
	}

	folderPath := filepath.Join("uploads", "category_photo", id)
	if _, err := os.Stat(folderPath); os.IsNotExist(err) {
		os.MkdirAll(folderPath, 0755)
	}

	savePath := filepath.Join(folderPath, file.Filename)
	log.Print(savePath)
	if err := c.SaveUploadedFile(file, savePath); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	if err := h.transactionRepo.InsertNewCategory(&input.CategoryName, &input.CategoryDesc, savePath, id); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Tambah kategori berhasil dilakukan"})
}

func (h *TransactionHandler) GetCategoriesList(c *gin.Context) {
	id := c.GetString("id")

	_page := c.DefaultQuery("page", "1")
	_limit := c.DefaultQuery("limit", "10")

	params := utils.ParsePagination(_page, _limit)

	// Take total of data in Category
	totalItem, err := h.transactionRepo.GetCategoryCount(id)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	rows, err := h.transactionRepo.GetCategoriesList(id, params.Limit, params.Offset)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	categories := make([]models.CategoryItem, 0)
	for rows.Next() {
		var cat models.CategoryItem
		if err := rows.Scan(&cat.CategoryId, &cat.CategoryName, &cat.CategoryDesc, &cat.CategoryUrlPhoto); err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}
		categories = append(categories, cat)
	}

	if err = rows.Err(); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"pagination": utils.BuildPaginationData(params, totalItem),
		"data":       categories,
	})
}

func (h *TransactionHandler) GetCategoryPhotoItem(c *gin.Context) {
	id := c.GetString("id")
	var input models.CategoryIdItem

	if err := c.ShouldBind(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	urlPhoto, err := h.transactionRepo.GetCategoryPhoto(&input.CategoryId, id)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	path := filepath.Join("uploads", "category_photo", id, filepath.Base(urlPhoto))
	if _, err := os.Stat(path); os.IsNotExist(err) {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "gambar tidak ditemukan"})
		return
	}

	c.File(path)
}

func (h *TransactionHandler) CreateTransaction(c *gin.Context) {
	id := c.GetString("id")
	var input models.TransactionInput

	if err := c.ShouldBind(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	output, err := h.balanceRepo.GetBalanceAmount(id)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	currentBalance := output.BalanceAmount

	diffBalanceAmount := currentBalance - input.TransactionTotalPrice

	if diffBalanceAmount < 0 {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Saldo tidak mencukupi untuk melakukan transaksi"})
		return
	}

	//Jika saldo telah berkurang dan server tiba" down. maka, lakukan rollback commit
	if err := h.transactionRepo.CreateTransaction(diffBalanceAmount, id, input, output.Id); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Transaksi telah berhasil dibuat"})
}

func (h *TransactionHandler) GetTransactionList(c *gin.Context) {
	id := c.GetString("id")

	_page := c.DefaultQuery("page", "1")
	_limit := c.DefaultQuery("limit", "10")

	params := utils.ParsePagination(_page, _limit)

	totalItem, err := h.transactionRepo.GetTransactionCount(id)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	rows, err := h.transactionRepo.GetTransactionList(id, params.Limit, params.Offset)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	transactions := make([]models.TransactionResult, 0)
	for rows.Next() {
		var transaction models.TransactionResult
		if err := rows.Scan(&transaction.TransactionId, &transaction.TransactionName,
			&transaction.TransactionDate, &transaction.TransactionAmount, &transaction.MerchantName,
			&transaction.TransactionPrice, &transaction.TransactionTotalPrice, &transaction.CategoryId); err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}
		transactions = append(transactions, transaction)
	}

	if err = rows.Err(); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"pagination": utils.BuildPaginationData(params, totalItem),
		"data":       transactions,
	})
}

func (h *TransactionHandler) RemoveCategory(c *gin.Context) {
	id := c.GetString("id")
	var input models.CategoryIdItem

	if err := c.ShouldBind(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	photoPath, err := h.transactionRepo.GetCategoryPhoto(&input.CategoryId, id)
	if err == sql.ErrNoRows {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Kategori tidak ditemukan"})
		return
	} else if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	if err := h.transactionRepo.RemoveCategoryItem(&input.CategoryId, id); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	if photoPath != "" {
		go func(path string) {
			os.Remove(path)
		}(photoPath)
	}

	c.JSON(http.StatusOK, gin.H{"message": "Kategori sudah berhasil dihapus"})

}
