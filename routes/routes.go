package routes

import (
	"github.com/gin-gonic/gin"
	"github.com/kunto/golang-rest-api-berita/controllers"
)

func UserRoute(r *gin.Engine) {
	r.GET("/users", controllers.GetUsers)
	r.POST("/users", controllers.CreateUser)
}
