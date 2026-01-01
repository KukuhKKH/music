package track

import (
	"git.dev.siap.id/kukuhkkh/app-music/app/middleware"
	"git.dev.siap.id/kukuhkkh/app-music/app/module/track/controller"
	"git.dev.siap.id/kukuhkkh/app-music/app/module/track/repository"
	"git.dev.siap.id/kukuhkkh/app-music/app/module/track/service"
	"github.com/gofiber/fiber/v2"
	"go.uber.org/fx"
)

type TrackRouter struct {
	App        fiber.Router
	Controller *controller.Controller
}

var NewTrackModule = fx.Options(
	// register repository of track module
	fx.Provide(repository.NewTrackRepository),

	// register service of track module
	fx.Provide(service.NewTrackService),

	// register controller of track module
	fx.Provide(controller.NewController),

	// register router of track module
	fx.Provide(NewTrackRouter),
)

func NewTrackRouter(fiber *fiber.App, controller *controller.Controller) *TrackRouter {
	return &TrackRouter{
		App:        fiber,
		Controller: controller,
	}
}

func (_i *TrackRouter) RegisterTrackRoutes() {
	// define controllers
	trackController := _i.Controller.Track

	// define routes
	_i.App.Route("/music", func(router fiber.Router) {
		router.Get("", middleware.Protected(), trackController.GetTracks)
	})
}
