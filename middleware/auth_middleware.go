package middleware

import (
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
)

func AuthMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {

		// ambil header Authorization
		authHeader := c.GetHeader("Authorization")

		if authHeader == "" {
			c.JSON(http.StatusUnauthorized, gin.H{
				"status":  false,
				"message": "Unauthorized - Token kosong",
			})
			c.Abort()
			return
		}

		// format header seharusnya: Bearer token123
		token := strings.TrimPrefix(authHeader, "Bearer ")

		// --- validasi token (contoh sederhana) ---
		if token != "secret123" {
			c.JSON(http.StatusUnauthorized, gin.H{
				"status":  false,
				"message": "Token tidak valid",
			})
			c.Abort()
			return
		}

		// jika valid → lanjut ke controller
		c.Next()
	}
}
