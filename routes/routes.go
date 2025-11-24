package routes

import (
	"github.com/gin-gonic/gin"
	"github.com/kunto/golang-rest-api-berita/controllers"
)

func UserRoute(r *gin.Engine) {
	// Group versi API
	v1 := r.Group("/api/v1")

	// Semua route users pakai middleware
	users := v1.Group("/users")
	//users.Use(middleware.AuthMiddleware())
	{
		users.POST("/", controllers.CreateUser)
		users.GET("/get-list-users", controllers.GetListUsers)
		users.GET("/:id", controllers.GetUserDetail)
		users.PUT("/:id", controllers.UpdateUser)
		users.DELETE("/:id", controllers.DeleteUser)
	}
}
