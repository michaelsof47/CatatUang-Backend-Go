package handlers

import (
	"catatuangbackend/models"
	"database/sql"
	"log"
	"net/http"
	"os"
	"path/filepath"

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
