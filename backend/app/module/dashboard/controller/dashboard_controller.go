package controller

import (
	"git.dev.siap.id/kukuhkkh/app-music/app/module/dashboard/service"
	"git.dev.siap.id/kukuhkkh/app-music/utils/response"
	"github.com/gofiber/fiber/v2"
)

type DashboardController interface {
	GetSummary(c *fiber.Ctx) error
}

type dashboardController struct {
	service service.DashboardService
}

func NewDashboardController(service service.DashboardService) DashboardController {
	return &dashboardController{
		service: service,
	}
}

// GetSummary godoc
// @Summary      Get dashboard summary
// @Description  Get total songs, total size and last upload time
// @Tags         Stats
// @Accept       json
// @Produce      json
// @Success      200 {object} response.Response
// @Router       /stats/summary [get]
func (ctrl *dashboardController) GetSummary(c *fiber.Ctx) error {
	res, err := ctrl.service.GetSummary()
	if err != nil {
		return err
	}

	return response.Resp(c, response.Response{
		Messages: response.Messages{"Get dashboard summary success"},
		Data:     res,
	})
}
