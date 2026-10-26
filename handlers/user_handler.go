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

	"catatuangbackend/repo"

	firebase "firebase.google.com/go/v4"
	"github.com/gin-gonic/gin"
	"github.com/redis/go-redis/v9"
	"golang.org/x/crypto/bcrypt"
)

type UserHandler struct {
	userRepo repo.UserInterface
	Redis    *redis.Client
	Firebase *firebase.App
}

func UserHandlerFunc(userRepo repo.UserInterface, redis *redis.Client, firebase *firebase.App) *UserHandler {
	return &UserHandler{userRepo: userRepo, Redis: redis, Firebase: firebase}
}

func (h *UserHandler) CreateUserNew(c *gin.Context) {
	var createUser model.User

	if err := c.ShouldBind(&createUser); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Check if account have available
	exists, err := h.userRepo.CheckAvailabilityAccount(createUser)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	if exists {
		c.JSON(http.StatusConflict, gin.H{"error": "Akun sudah tersedia. Silahkan gunakan akun lain"})
		return
	}

	//retrieve file from body
	file, err := c.FormFile("url_user_image")
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Wajib isi foto profil"})
		return
	}

	//create a directory
	folderPath := filepath.Join("uploads", "image_profile")
	if _, err := os.Stat(folderPath); os.IsNotExist(err) {
		os.MkdirAll(folderPath, 0755)
	}

	//store file into /uploads
	savePath := filepath.Join(folderPath, file.Filename)
	log.Println("Directory image stored at " + savePath)

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

	outputUser, err := h.userRepo.InsertNewAccount(createUser, savePath, hashedPassword)
	if err != nil {
		if err.Error() == "Email sudah terdaftar" {
			c.JSON(http.StatusConflict, gin.H{"error": "email sudah terdaftar"})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusCreated, gin.H{"id": outputUser.Id})
}

func (h *UserHandler) GetUserById(c *gin.Context) {
	id := c.Param("id")

	user, err := h.userRepo.GetUserById(id)
	if err.Error() == "User tidak ditemukan" {
		c.JSON(http.StatusNotFound, gin.H{"error": "User tidak ditemukan"})
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

	user, err := h.userRepo.GetUserByPhoneAndEmail(login.EmailOrPhone)

	loginFailed := err == sql.ErrNoRows

	if err != nil {
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
		go func(token string) {
			claims, err := utils.ValidateToken(token)
			if err == nil {
				exp := int64(claims["exp"].(float64))
				ttl := time.Until(time.Unix(exp, 0))
				if ttl > 0 {
					h.Redis.Set(config.Context, "blacklist:"+token, "true", ttl)
				}
			}
		}(oldToken)
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

	go func(token, userId string, duration time.Duration) {
		h.Redis.Set(config.Context, "blacklist:"+tokenString, "true", duration)
		h.Redis.Del(config.Context, "active_token"+userId)
	}(tokenString, id, ttl)

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
	user, err := h.userRepo.GetUserPasswordById(id)
	if err.Error() == "Akun Password tidak ditemukan" {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Akun Password tidak ditemukan"})
		return
	} else if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
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

	//Update Password based on id and hashedpassword
	err = h.userRepo.UpdatePasswordById(hashedPassword, id)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Password berhasil diubah"})
}

func (h *UserHandler) ForgotPassword(c *gin.Context) {
	var input models.UserLogin

	if err := c.ShouldBind(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	output, err := h.userRepo.GetUserByPhoneAndEmail(input.EmailOrPhone)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Email atau No. Telp Tidak Ditemukan"})
		return
	}

	if err := bcrypt.CompareHashAndPassword([]byte(output.Password), []byte(input.Password)); err == nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "password yang diubah sama seperti sebelumnya. Silahkan gunakan password lain"})
		return
	}

	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(input.Password), bcrypt.DefaultCost)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	err = h.userRepo.UpdatePasswordBasedOnEmailAndPhone(hashedPassword, input.EmailOrPhone)
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

	err := h.userRepo.UpdateAccountUserInfo(input.FirstName, input.LastName, input.Email, input.Phone, id)
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

	oldImagePath, err := h.userRepo.GetUrlImagePhoto(id)
	if err.Error() == "Foto Akun tidak ditemukan" {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Foto Akun tidak ditemukan"})
		return
	} else if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	folderPath := filepath.Join("uploads", "image_profile")
	if _, err := os.Stat(folderPath); os.IsNotExist(err) {
		os.MkdirAll(folderPath, 0755)
	}

	savePath := filepath.Join(folderPath, file.Filename)
	if err := c.SaveUploadedFile(file, savePath); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "gagal menyimpan gambar"})
		return
	}

	err = h.userRepo.UpdatePhotoProfile(savePath, id)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	log.Print("image old: " + oldImagePath)

	if oldImagePath != "" {
		go func(path string) {
			os.Remove(path)
		}(oldImagePath)
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

func (h *UserHandler) GetImageProfile(c *gin.Context) {
	id := c.GetString("id")

	photoPath, err := h.userRepo.GetUrlImagePhoto(id)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
	}

	path := filepath.Join("uploads", "image_profile", filepath.Base(photoPath))
	if _, err := os.Stat(path); os.IsNotExist(err) {
		c.JSON(http.StatusNotFound, gin.H{"error": "gambar tidak ditemukan"})
		return
	}

	c.File(path)
}
