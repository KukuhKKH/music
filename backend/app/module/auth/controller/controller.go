package controller

import (
	"git.dev.siap.id/kukuhkkh/app-music/app/module/auth/service"
	"git.dev.siap.id/kukuhkkh/app-music/utils/config"
)

type Controller struct {
	Auth AuthController
}

func NewController(authService service.AuthService, cfg *config.Config) *Controller {
	return &Controller{
		Auth: NewAuthController(authService, cfg),
	}
}
