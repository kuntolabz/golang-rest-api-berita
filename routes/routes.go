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
		users.GET("/get-list-users", controllers.GetUsers)

		// users.GET("/:id", controllers.GetUserByID)   // GET /api/v1/users/:id
		// users.PUT("/:id", controllers.UpdateUser)    // PUT /api/v1/users/:id
		// users.DELETE("/:id", controllers.DeleteUser) // DELETE /api/v1/users/:id
	}
}
