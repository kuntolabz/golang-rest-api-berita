package routes

import (
	"github.com/gin-gonic/gin"
	"github.com/kunto/golang-rest-api-berita/config"
	"github.com/kunto/golang-rest-api-berita/controllers"
	"github.com/kunto/golang-rest-api-berita/middleware"
	"github.com/kunto/golang-rest-api-berita/repositories"
	"github.com/kunto/golang-rest-api-berita/services"
)

func UserRoute(r *gin.Engine) {

	// Repository
	userRepo := repositories.NewUserRepository(config.DB)

	// Service
	userService := services.NewUserService(userRepo)
	authService := services.NewAuthService(userRepo)

	// Controller
	userCtrl := controllers.NewUserController(userService)
	authCtrl := controllers.NewAuthController(authService)

	// AUTH (Public)
	api := r.Group("/api/v1/auth")
	{
		api.POST("/login", authCtrl.Login)
	}

	// USER (Protected)
	v1 := r.Group("/api/v1")
	users := v1.Group("/users")
	users.Use(middleware.AuthMiddleware()) // 👈 Tambahkan middleware
	{
		users.POST("/create", userCtrl.CreateUser)
		users.GET("/list", userCtrl.GetListUsers)
		users.GET("/detail/:id", userCtrl.GetByID)
		users.PUT("/update/:id", userCtrl.UpdateUser)
		users.DELETE("/delete/:id", userCtrl.DeleteUser)
	}
}
