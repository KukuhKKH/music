package controller

import (
	"git.dev.siap.id/kukuhkkh/app-music/app/module/track/service"
	"git.dev.siap.id/kukuhkkh/app-music/utils/paginator"
	"git.dev.siap.id/kukuhkkh/app-music/utils/response"
	"github.com/gofiber/fiber/v2"
)

type trackController struct {
	trackService service.TrackService
}

type TrackController interface {
	GetTracks(c *fiber.Ctx) error
}

func NewTrackController(trackService service.TrackService) TrackController {
	return &trackController{
		trackService: trackService,
	}
}

// GetTracks godoc
// @Summary      Get paginated tracks
// @Description  Get list of tracks with search and pagination
// @Tags         Music
// @Accept       json
// @Produce      json
// @Param        search query string false "Search by title or artist"
// @Param        page   query int    false "Page number"
// @Param        limit  query int    false "Items per page"
// @Success      200 {object} response.Response
// @Router       /music [get]
func (_i *trackController) GetTracks(c *fiber.Ctx) error {
	p, _ := paginator.Paginate(c)
	search := c.Query("search")

	tracks, p, err := _i.trackService.GetPaginatedTracks(search, p)
	if err != nil {
		return err
	}

	return response.Resp(c, response.Response{
		Messages: response.Messages{"Get tracks success"},
		Data:     tracks,
		Meta:     paginator.Paging(p),
	})
}
