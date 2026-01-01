package controller

import "git.dev.siap.id/kukuhkkh/app-music/app/module/dashboard/service"

type Controller struct {
	Dashboard DashboardController
}

func NewController(dashboardService service.DashboardService) *Controller {
	return &Controller{
		Dashboard: NewDashboardController(dashboardService),
	}
}
