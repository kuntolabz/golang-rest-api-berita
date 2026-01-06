package middleware

import (
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/kunto/golang-rest-api-berita/utils"
)

// AuthMiddleware -> cek token dari header Authorization
func AuthMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		// Ambil token dari header
		authHeader := c.GetHeader("Authorization")

		if authHeader == "" {
			utils.ErrorResponse(c, http.StatusUnauthorized, "Authorization header required")
			c.Abort()
			return

		}

		// Format: "Bearer <token>"
		token := strings.Split(authHeader, " ")
		if len(token) != 2 {
			utils.ErrorResponse(c, http.StatusUnauthorized, "Invalid token format")
			c.Abort()
			return
		}

		// Validasi token
		userData, err := utils.ValidateToken(token[1])
		if err != nil {
			utils.ErrorResponse(c, http.StatusUnauthorized, "Invalid or expired token")
			c.Abort()
			return
		}

		// Set user context biar bisa dipakai di controller/service
		c.Set("user", userData)
		newCtx := utils.SetUserID(c.Request.Context(), userData.UserID)
		c.Request = c.Request.WithContext(newCtx)

		c.Next()
	}
}
