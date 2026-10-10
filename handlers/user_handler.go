package handlers

import (
	"catatuangbackend/config"
	"catatuangbackend/models"
	model "catatuangbackend/models"
	"catatuangbackend/utils"
	"database/sql"
	"log"
	"net/http"
	"os"
	"path/filepath"
	"strings"
	"time"

	firebase "firebase.google.com/go/v4"
	"github.com/gin-gonic/gin"
	"github.com/lib/pq"
	"github.com/redis/go-redis/v9"
	"golang.org/x/crypto/bcrypt"
)

type UserHandler struct {
	DB       *sql.DB
	Redis    *redis.Client
	Firebase *firebase.App
}

func UserHandlerFunc(db *sql.DB, redis *redis.Client, firebase *firebase.App) *UserHandler {
	return &UserHandler{DB: db, Redis: redis, Firebase: firebase}
}

func (h *UserHandler) CreateUserNew(c *gin.Context) {
	var createUser model.User

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

	var user model.UserResult
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

	var user model.UserResult
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
	var login model.UserLogin

	if err := c.ShouldBind(&login); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	key := "Login Attempt: " + login.EmailOrPhone

	// Cek jumlah percobaan gagal sejauh ini
	attempts, _ := h.Redis.Get(config.Context, key).Int()
	if attempts > 5 {
		c.JSON(http.StatusTooManyRequests, gin.H{"error": "terlalu banyak percobaan, tunggu 5 menit"})
		return
	}

	var user model.UserResult
	query := `SELECT id, password FROM "User" WHERE email = $1 OR phone = $1`
	err := h.DB.QueryRow(query, login.EmailOrPhone).Scan(&user.Id, &user.Password)

	loginFailed := err == sql.ErrNoRows

	if err != nil && err != sql.ErrNoRows {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	if !loginFailed {
		// let's check diff password from body and password from database
		if err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(login.Password)); err != nil {
			loginFailed = true
		}
	}

	if loginFailed {
		h.Redis.Incr(config.Context, key)
		h.Redis.Expire(config.Context, key, 5*time.Minute)
		c.JSON(http.StatusUnauthorized, gin.H{"error": "akun dan password salah"})
		return
	}

	// Login success, reset counter
	h.Redis.Del(config.Context, key)

	// Stored token into redis
	activeTokenKey := "active_token:" + user.Id

	// Check if user have active token from before session
	oldToken, err := h.Redis.Get(config.Context, activeTokenKey).Result()
	if err == nil && oldToken != "" {
		claims, err := utils.ValidateToken(oldToken)
		if err == nil {
			exp := int64(claims["exp"].(float64))
			ttl := time.Until(time.Unix(exp, 0))
			if ttl > 0 {
				h.Redis.Set(config.Context, "blacklist:"+oldToken, "true", ttl)
			}
		}
	}

	// generate token from JWT
	token, err := utils.GenerateToken(user.Id)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	h.Redis.Set(config.Context, activeTokenKey, token, 7*24*time.Hour)

	c.JSON(http.StatusOK, gin.H{"id": user.Id, "token": token})
}

func (h *UserHandler) LogoutUserAccount(c *gin.Context) {
	id := c.GetString("id")
	authHeader := c.GetHeader("Authorization")
	tokenString := strings.TrimPrefix(authHeader, "Bearer ")

	claims, err := utils.ValidateToken(tokenString)
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "token tidak valid"})
		return
	}

	exp := int64(claims["exp"].(float64))
	ttl := time.Until(time.Unix(exp, 0))

	h.Redis.Set(config.Context, "blacklist:"+tokenString, "true", ttl)
	h.Redis.Del(config.Context, "active_token"+id)

	c.JSON(http.StatusOK, gin.H{"message": "Logout Berhasil"})
}

func (h *UserHandler) UpdatePassword(c *gin.Context) {
	var input model.PasswordCredentials
	id := c.GetString("id")

	if err := c.ShouldBind(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Get an User Data
	var user model.User
	query := `SELECT password FROM "User" WHERE id = $1`
	err := h.DB.QueryRow(query, id).Scan(&user.Password)
	if err != sql.ErrNoRows && err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
	}

	// Check if password still same as before. it will be received error message
	if err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(input.Password)); err == nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "password yang diubah sama seperti sebelumnya. Silahkan gunakan password lain"})
		return
	}

	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(input.Password), bcrypt.DefaultCost)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	query = `UPDATE "User" SET password = $1 WHERE id = $2`
	_, err = h.DB.Exec(query, hashedPassword, id)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Password berhasil diubah"})
}

func (h *UserHandler) UpdateUserAccount(c *gin.Context) {
	var input models.UserNoPassword
	id := c.GetString("id")

	if err := c.ShouldBind(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	query := `UPDATE "User" SET first_name = $1, last_name = $2, email = $3, phone = $4 WHERE id = $5`
	_, err := h.DB.Exec(query, &input.FirstName, &input.LastName, &input.Email, &input.Phone, id)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Data telah berhasil diubah"})
}

func (h *UserHandler) UpdateProfileImage(c *gin.Context) {
	id := c.GetString("id")

	file, err := c.FormFile("url_user_image")
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "file gambar wajib diunggah"})
		return
	}

	var oldImagePath string
	query := `SELECT url_user_image FROM "User" WHERE id = $1`
	err = h.DB.QueryRow(query, id).Scan(&oldImagePath)
	if err != nil && err != sql.ErrNoRows {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	savePath := "uploads/" + file.Filename
	if err := c.SaveUploadedFile(file, savePath); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "gagal menyimpan gambar"})
		return
	}

	query = `UPDATE "User" SET url_user_image = $1 WHERE id = $2`
	_, err = h.DB.Exec(query, savePath, id)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	log.Print("image old: " + oldImagePath)

	if oldImagePath != "" {
		os.Remove(oldImagePath)
	}

	c.JSON(http.StatusOK, gin.H{"message": "foto profil berhasil diperbaharui", "image_url": savePath})
}

func (h *UserHandler) VerifyFirebaseToken(c *gin.Context) {
	var input model.VerifyFirebaseToken

	if err := c.ShouldBind(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	authClient, err := h.Firebase.Auth(c)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Gagal Inisialisasi Auth"})
		return
	}

	token, err := authClient.VerifyIDToken(c, input.FirebaseToken)
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "token tidak valid " + err.Error()})
		return
	}

	phoneNumber := token.Claims["phone_number"].(string)

	c.JSON(http.StatusOK, gin.H{"message": "Akun Berhasil di Verifikasi", "phone": phoneNumber})
}

func (h *UserHandler) RetrieveImageProfile(c *gin.Context) {
	id := c.GetString("id")

	var output models.UserResult
	query := `SELECT url_user_image FROM "User" WHERE id = $1`
	err := h.DB.QueryRow(query, id).Scan(&output.UrlUserImage)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
	}

	path := filepath.Join("uploads", filepath.Base(output.UrlUserImage))
	if _, err := os.Stat(path); os.IsNotExist(err) {
		c.JSON(http.StatusNotFound, gin.H{"error": "gambar tidak ditemukan"})
		return
	}

	c.File(path)
}
