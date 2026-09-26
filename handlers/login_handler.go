package handlers

import (
	inputjson "catatuangbackend/models/input_json"
	resultjson "catatuangbackend/models/result_json"
	"database/sql"
	"log"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/lib/pq"
)

type UserHandler struct {
	DB *sql.DB
}

func UserHandlerFunc(db *sql.DB) *UserHandler {
	return &UserHandler{DB: db}
}

func (h *UserHandler) CreateUserNew(c *gin.Context) {
	var createUser inputjson.User

	if err := c.ShouldBind(&createUser); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	//store file
	file, err := c.FormFile("url_user_image")
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "url_user_image wajib diisi"})
		return
	}

	savePath := "uploads/" + file.Filename
	log.Println(savePath)
	if err := c.SaveUploadedFile(file, savePath); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	var user resultjson.User
	query := `INSERT INTO "User" (first_name, last_name, url_user_image, email, phone, password) 
			  VALUES ($1, $2, $3, $4, $5, $6)
			  RETURNING id`
	err = h.DB.QueryRow(query, createUser.FirstName, createUser.LastName, savePath, createUser.Email, createUser.Phone, createUser.Password).Scan(&user.Id)
	if err != nil {
		if pqErr, ok := err.(*pq.Error); ok && pqErr.Code == "23505" {
			c.JSON(http.StatusConflict, gin.H{"error": "email sudah terdaftar"})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusCreated, gin.H{"id": user.Id})
}

func (h *UserHandler) GetUserById(c *gin.Context) {
	id := c.Param("id")

	var user resultjson.User
	query := `SELECT id, first_name, last_name, reward_status, url_user_image, email, phone, password, created_at WHERE id = $1`
	err := h.DB.QueryRow(query, id).Scan(&user.Id, &user.FirstName, &user.LastName, &user.RewardStatus, &user.UrlUserImage, &user.Email, &user.Phone, &user.Password, &user.CreatedAt)
	if err == sql.ErrNoRows {
		c.JSON(http.StatusNotFound, gin.H{"error": "user tidak ditemukan"})
		return
	} else if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, user)
}
