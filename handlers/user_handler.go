package handlers

import (
	inputjson "catatuangbackend/models/input_json"
	resultjson "catatuangbackend/models/result_json"
	"catatuangbackend/utils"
	"database/sql"
	"log"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/lib/pq"
	"golang.org/x/crypto/bcrypt"
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

	//retrieve file from body
	file, err := c.FormFile("url_user_image")
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "url_user_image wajib diisi"})
		return
	}

	//store file into /uploads
	savePath := "uploads/" + file.Filename
	log.Println(savePath)
	if err := c.SaveUploadedFile(file, savePath); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	//hashed password from string password
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(createUser.Password), bcrypt.DefaultCost)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	var user resultjson.User
	query := `INSERT INTO "User" (first_name, last_name, url_user_image, email, phone, password) 
			  VALUES ($1, $2, $3, $4, $5, $6)
			  RETURNING id`
	err = h.DB.QueryRow(query, createUser.FirstName, createUser.LastName, savePath, createUser.Email, createUser.Phone, hashedPassword).Scan(&user.Id)
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
	query := `SELECT id, first_name, last_name, reward_status, url_user_image, email, phone, password, "createdAt" FROM "User" WHERE id = $1`
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

func (h *UserHandler) LoginUserAccount(c *gin.Context) {
	var login inputjson.UserLogin

	if err := c.ShouldBind(&login); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	var user resultjson.User
	query := `SELECT id, password FROM "User" WHERE email = $1 OR phone = $1`
	err := h.DB.QueryRow(query, login.EmailOrPhone).Scan(&user.Id, &user.Password)
	if err == sql.ErrNoRows {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "akun atau password salah"})
		return
	} else if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	log.Println(login.Password)
	log.Println(user.Password)

	// let's check diff password from body and password from database
	if err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(login.Password)); err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "password salah"})
		return
	}

	// generate token from JWT
	token, err := utils.GenerateToken(user.Id)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"id": user.Id, "token": token})
}
