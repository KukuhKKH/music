package dashboard

import (
	"git.dev.siap.id/kukuhkkh/app-music/app/middleware"
	"git.dev.siap.id/kukuhkkh/app-music/app/module/dashboard/controller"
	"git.dev.siap.id/kukuhkkh/app-music/app/module/dashboard/repository"
	"git.dev.siap.id/kukuhkkh/app-music/app/module/dashboard/service"
	"github.com/gofiber/fiber/v2"
	"go.uber.org/fx"
)

type DashboardRouter struct {
	App        fiber.Router
	Controller *controller.Controller
}

var NewDashboardModule = fx.Options(
	fx.Provide(repository.NewDashboardRepository),
	fx.Provide(service.NewDashboardService),
	fx.Provide(controller.NewController),
	fx.Provide(NewDashboardRouter),
)

func NewDashboardRouter(fiber *fiber.App, controller *controller.Controller) *DashboardRouter {
	return &DashboardRouter{
		App:        fiber,
		Controller: controller,
	}
}

func (_i *DashboardRouter) RegisterDashboardRoutes() {
	// define controllers
	dashboardController := _i.Controller.Dashboard

	_i.App.Get("/stats/summary", middleware.RequireAuth(), dashboardController.GetSummary)
}
