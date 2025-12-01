package controllers

import (
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

	utils.ResponseSuccess(c, token, "Login berhasil")

}
