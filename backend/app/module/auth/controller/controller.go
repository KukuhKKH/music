package controller

import "git.dev.siap.id/kukuhkkh/app-music/app/module/auth/service"

type Controller struct {
	Auth AuthController
}

func NewController(authService service.AuthService) *Controller {
	return &Controller{
		Auth: NewAuthController(authService),
	}
}
