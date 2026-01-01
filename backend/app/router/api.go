package router

import (
	"git.dev.siap.id/kukuhkkh/app-music/app/module/auth"
	"git.dev.siap.id/kukuhkkh/app-music/app/module/track"
	"git.dev.siap.id/kukuhkkh/app-music/utils/config"
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/swagger"
)

type Router struct {
	App         fiber.Router
	Cfg         *config.Config
	AuthRouter  *auth.AuthRouter
	TrackRouter *track.TrackRouter
}

func NewRouter(
	fiber *fiber.App,
	cfg *config.Config,
	authRouter *auth.AuthRouter,
	trackRouter *track.TrackRouter,
) *Router {
	return &Router{
		App:         fiber,
		Cfg:         cfg,
		AuthRouter:  authRouter,
		TrackRouter: trackRouter,
	}
}

func (r *Router) Register() {
	// Test Routes
	r.App.Get("/ping", func(c *fiber.Ctx) error {
		return c.SendString("Pong! 👋")
	})

	// Swagger Documentation
	r.App.Get("/swagger/*", swagger.New(swagger.Config{
		DefaultModelsExpandDepth: -1,
	}))

	// routes of modules
	r.AuthRouter.RegisterAuthRoutes()
	r.TrackRouter.RegisterTrackRoutes()
}
