package controllers

import (
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"

	"github.com/kunto/golang-rest-api-berita/config"
	"github.com/kunto/golang-rest-api-berita/dto"
	"github.com/kunto/golang-rest-api-berita/models"
	"github.com/kunto/golang-rest-api-berita/utils"
)

func GetUsers(c *gin.Context) {
	search := c.Query("search")
	limitStr := c.Query("limit")
	offsetStr := c.Query("offset")

	limit := 10
	offset := 0

	fmt.Sscanf(limitStr, "%d", &limit)
	fmt.Sscanf(offsetStr, "%d", &offset)

	where := "WHERE 1=1"
	params := []interface{}{}

	if search != "" {
		where += " AND (name LIKE ? OR email LIKE ? OR username LIKE ?)"
		like := "%" + search + "%"
		params = append(params, like, like, like)
	}

	// Hitung total row
	var total int64
	totalQuery := "SELECT COUNT(*) FROM users " + where
	config.DB.Raw(totalQuery, params...).Scan(&total)

	// Data list
	var users []dto.UserDTO

	dataQuery := `
        SELECT id_user, name, email, username 
        FROM users
    ` + where + ` ORDER BY name ASC LIMIT ? OFFSET ?`

	params = append(params, limit, offset)
	config.DB.Raw(dataQuery, params...).Scan(&users)

	// 🔥 Gunakan fungsi reusable
	utils.SuccessResponse(c, users, total, offset, limit)
}

func CreateUser(c *gin.Context) {
	var user models.User
	if err := c.ShouldBindJSON(&user); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	config.DB.Create(&user)

	c.JSON(http.StatusOK, user)
}
