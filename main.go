package main

import (
	"log"

	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
	"github.com/kunto/golang-rest-api-berita/config"
	"github.com/kunto/golang-rest-api-berita/routes"
)

func main() {
	err := godotenv.Load()
	if err != nil {
		log.Fatal(".env tidak ditemukan")
	}

	config.ConnectDB()

	r := gin.Default()
	//r.Use(middleware.AuthMiddleware())

	routes.UserRoute(r)

	r.Run(":8080")
}
