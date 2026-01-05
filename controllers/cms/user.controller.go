package controllers

import (
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	dto "github.com/kunto/golang-rest-api-berita/dto/cms"
	"github.com/kunto/golang-rest-api-berita/services"
	"github.com/kunto/golang-rest-api-berita/utils"
)

type UserController struct {
	service services.UserService
}

func NewUserController(service services.UserService) *UserController {
	return &UserController{service}
}

func (c *UserController) GetListUsers(ctx *gin.Context) {
	// Ambil query param
	search := ctx.DefaultQuery("search", "")
	limitQuery := ctx.DefaultQuery("limit", "10")
	pageQuery := ctx.DefaultQuery("page", "1")

	// Convert ke int
	limit, _ := strconv.Atoi(limitQuery)
	page, _ := strconv.Atoi(pageQuery)

	// hitung offset
	offset := (page - 1) * limit

	// panggil service
	users, total, err := c.service.GetList(search, limit, offset)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"error": err.Error(),
		})
		return
	}

	// response
	ctx.JSON(http.StatusOK, gin.H{
		"data":  users,
		"total": total,
		"page":  page,
		"limit": limit,
	})
}

func (c *UserController) CreateUser(ctx *gin.Context) {
	var input dto.InsertUserDTO
	if err := ctx.ShouldBindJSON(&input); err != nil {
		utils.ErrorResponse(ctx, 400, "invalid payload")
		return
	}

	user, err := c.service.Create(input)
	if err != nil {
		utils.ErrorResponse(ctx, 400, err.Error())
		return
	}

	// 3️⃣ Response ke client
	utils.ResponseSuccess(ctx, user, "Data user berhasil dibuat")
}

func (c *UserController) GetByID(ctx *gin.Context) {
	id := ctx.Param("id")
	user, err := c.service.GetDetail(id)
	if err != nil {
		utils.ErrorResponse(ctx, 404, err.Error())
		return
	}
	utils.ResponseSuccess(ctx, user, "Data user berhasil diambil")

}

func (c *UserController) UpdateUser(ctx *gin.Context) {
	id := ctx.Param("id")

	var input dto.InsertUserDTO
	if err := ctx.ShouldBindJSON(&input); err != nil {
		utils.ErrorResponse(ctx, 400, "payload tidak valid")
		return
	}

	result, err := c.service.Update(id, input)
	if err != nil {
		utils.ErrorResponse(ctx, 400, err.Error())
		return
	}

	utils.ResponseSuccess(ctx, result, "Data user berhasil diupdate")
}

func (c *UserController) DeleteUser(ctx *gin.Context) {
	id := ctx.Param("id")
	_, err := c.service.GetDetail(id)
	if err != nil {
		utils.ErrorResponse(ctx, 404, err.Error())
		return
	}
	if err := c.service.Delete(id); err != nil {
		utils.ErrorResponse(ctx, 400, "Data User gagal dihapus")
		return
	}
	utils.ResponseSuccess(ctx, nil, "Data user berhasil dihapus")

}
