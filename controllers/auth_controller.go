package controllers

import (
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/kunto/golang-rest-api-berita/dto"
	"github.com/kunto/golang-rest-api-berita/services"
	"github.com/kunto/golang-rest-api-berita/utils"
)

type AuthController struct {
	AuthService services.AuthService
}

func NewAuthController(service services.AuthService) *AuthController {
	return &AuthController{AuthService: service}
}

func (ctrl *AuthController) Login(c *gin.Context) {
	var req dto.LoginRequest

	// Validasi JSON input
	if err := c.ShouldBindJSON(&req); err != nil {
		utils.ErrorResponse(c, 404, err.Error())
		return
	}

	token, err := ctrl.AuthService.Login(req.Email, req.Password)
	if err != nil {
		utils.ErrorResponse(c, 404, err.Error())
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"Status":       "success",
		"ResponseCode": 200,
		"ResponseDesc": "Login Successful",
		"Timestamp":    time.Now(),
		"Token":        token,
	})

}
