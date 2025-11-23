package utils

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

func SuccessResponse(c *gin.Context, data interface{}, total int64, offset int, limit int) {
	c.JSON(http.StatusOK, gin.H{
		"responseCode": 200,
		"dataList":     data,
		"pagination": gin.H{
			"totalRow": total,
			"offset":   offset,
			"limit":    limit,
		},
	})
}

func ErrorResponse(c *gin.Context, code int, message string) {
	c.JSON(code, gin.H{
		"responseCode": code,
		"message":      message,
	})
}
