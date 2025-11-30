package routes

import (
	"github.com/gin-gonic/gin"
	"github.com/kunto/golang-rest-api-berita/config"
	"github.com/kunto/golang-rest-api-berita/controllers"
	"github.com/kunto/golang-rest-api-berita/repositories"
	"github.com/kunto/golang-rest-api-berita/services"
)

func UserRoute(r *gin.Engine) {

	userRepo := repositories.NewUserRepository(config.DB)
	userService := services.NewUserService(userRepo)
	userCtrl := controllers.NewUserController(userService)

	v1 := r.Group("/api/v1")
	users := v1.Group("/users")
	{
		users.POST("/create", userCtrl.CreateUser)
		users.GET("/list", userCtrl.GetListUsers)
		users.GET("/detail/:id", userCtrl.GetByID)
		users.PUT("/update/:id", userCtrl.UpdateUser)
		users.DELETE("/delete/:id", userCtrl.DeleteUser)
	}
}
