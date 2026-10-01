package middleware

import (
	"catatuangbackend/config"
	"catatuangbackend/utils"
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/redis/go-redis/v9"
)

func AuthMiddleware(redisdb *redis.Client) gin.HandlerFunc {
	return func(c *gin.Context) {
		authHeader := c.GetHeader("Authorization")
		if authHeader == "" {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Token tidak ditemukan"})
			c.Abort()
			return
		}

		tokenString := strings.TrimPrefix(authHeader, "Bearer ")

		exists, err := redisdb.Exists(config.Context, "blacklist:"+tokenString).Result()
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "gagal memverifikasi token"})
			c.Abort()
			return
		}

		if exists > 0 {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "token sudah tidak berlaku, silahkan login ulang"})
			c.Abort()
			return
		}

		claims, err := utils.ValidateToken(tokenString)
		if err != nil {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Token tidak valid atau expired"})
			c.Abort()
			return
		}

		id, ok := claims["id"].(string)
		if !ok {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "token tidak valid"})
			c.Abort()
			return
		}

		c.Set("id", id)
		c.Next()
	}
}
