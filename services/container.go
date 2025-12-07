package services

import (
	"github.com/kunto/golang-rest-api-berita/config"
	"github.com/kunto/golang-rest-api-berita/repositories"
)

type ServiceContainer struct {
	UserService UserService
	AuthService AuthService
	// Tambah service lain di sini...
}

func NewServiceContainer() *ServiceContainer {

	// Repository
	userRepo := repositories.NewUserRepository(config.DB)

	return &ServiceContainer{
		UserService: NewUserService(userRepo),
		AuthService: NewAuthService(userRepo),
	}
}
