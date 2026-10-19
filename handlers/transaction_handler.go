package handlers

import (
	"catatuangbackend/models"
	"database/sql"
	"log"
	"math"
	"net/http"
	"os"
	"path/filepath"
	"strconv"

	"github.com/gin-gonic/gin"
)

type TransactionHandler struct {
	DB *sql.DB
}

func TransactionHandlerFunc(db *sql.DB) *TransactionHandler {
	return &TransactionHandler{DB: db}
}

func (h *TransactionHandler) CreateCategory(c *gin.Context) {
	var input models.CreateCategory
	id := c.GetString("id")

	if err := c.ShouldBind(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Check if category is available
	var exists bool
	isAccAvailableQuery := `SELECT EXISTS(SELECT * FROM "Category" WHERE LOWER(name) = LOWER($1))`
	if err := h.DB.QueryRow(isAccAvailableQuery, &input.CategoryName).Scan(&exists); err != nil {
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

	query := `INSERT INTO "Category" (name, description, url_image, user_id) VALUES ($1, $2, $3, $4)`
	_, err = h.DB.Exec(query, input.CategoryName, input.CategoryDesc, savePath, id)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Tambah kategori berhasil dilakukan"})
}

func (h *TransactionHandler) GetCategoriesList(c *gin.Context) {
	id := c.GetString("id")

	_page := c.DefaultQuery("page", "1")
	_limit := c.DefaultQuery("limit", "10")

	page, err := strconv.Atoi(_page)
	if err != nil || page < 1 {
		page = 1
	}

	limit, err := strconv.Atoi(_limit)
	if err != nil || limit < 1 {
		limit = 10
	}

	offset := (page - 1) * limit

	// Take total of data in Category
	var totalItem int
	countQuery := `SELECT COUNT(*) FROM "Category" WHERE user_id = $1`
	if err := h.DB.QueryRow(countQuery, id).Scan(&totalItem); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	query := `SELECT id, name, description, url_image FROM "Category" WHERE user_id = $1 ORDER BY id DESC LIMIT $2 OFFSET $3`
	rows, err := h.DB.Query(query, id, limit, offset)
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

	totalPages := int(math.Ceil(float64(totalItem) / float64(limit)))

	c.JSON(http.StatusOK, gin.H{
		"pagination": gin.H{
			"current_page": page,
			"limit":        limit,
			"total_data":   totalItem,
			"total_pages":  totalPages,
		},
		"data": categories,
	})

}
