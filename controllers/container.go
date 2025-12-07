package controllers

import "github.com/kunto/golang-rest-api-berita/services"

type ControllerContainer struct {
	User *UserController
	Auth *AuthController
}

func NewControllerContainer(s *services.ServiceContainer) *ControllerContainer {
	return &ControllerContainer{
		User: NewUserController(s.UserService),
		Auth: NewAuthController(s.AuthService),
	}
}
