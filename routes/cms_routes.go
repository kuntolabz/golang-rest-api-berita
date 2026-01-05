package routes

import (
	"github.com/gin-gonic/gin"
	controllers "github.com/kunto/golang-rest-api-berita/controllers/cms"
	"github.com/kunto/golang-rest-api-berita/middleware"
	"github.com/kunto/golang-rest-api-berita/services"
)

func InitRoutes(r *gin.Engine) {

	// INIT SERVICES
	servicesContainer := services.NewServiceContainer()

	// INIT CONTROLLERS
	c := controllers.NewControllerContainer(servicesContainer)

	// AUTH (Public)
	api := r.Group("/api/v1/auth")
	{
		api.POST("/login", c.Auth.Login)
	}

	// USER (Protected)
	v1 := r.Group("/api/v1")
	users := v1.Group("/users")
	users.Use(middleware.AuthMiddleware()) // 👈 Tambahkan middleware
	{
		users.POST("/create", c.User.CreateUser)
		users.GET("/list", c.User.GetListUsers)
		users.GET("/detail/:id", c.User.GetByID)
		users.PUT("/update/:id", c.User.UpdateUser)
		users.DELETE("/delete/:id", c.User.DeleteUser)
	}
}
