package controller

import (
	"git.dev.siap.id/kukuhkkh/app-music/app/middleware"
	"git.dev.siap.id/kukuhkkh/app-music/app/module/track/request"
	"git.dev.siap.id/kukuhkkh/app-music/app/module/track/service"
	"git.dev.siap.id/kukuhkkh/app-music/utils/paginator"
	"git.dev.siap.id/kukuhkkh/app-music/utils/response"
	"github.com/gofiber/fiber/v2"
	"github.com/golang-jwt/jwt/v4"
)

type trackController struct {
	trackService service.TrackService
}

type TrackController interface {
	GetTracks(c *fiber.Ctx) error
	GetTrackByID(c *fiber.Ctx) error
	Create(c *fiber.Ctx) error
	Update(c *fiber.Ctx) error
	Delete(c *fiber.Ctx) error
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

// GetTrackByID godoc
// @Summary      Get track by ID
// @Description  Get a single track by its ID
// @Tags         Music
// @Accept       json
// @Produce      json
// @Param        id   path uint64 true "Track ID"
// @Success      200 {object} response.Response
// @Router       /music/{id} [get]
func (_i *trackController) GetTrackByID(c *fiber.Ctx) error {
	id, err := c.ParamsInt("id")
	if err != nil {
		return err
	}

	res, err := _i.trackService.GetTrackByID(uint64(id))
	if err != nil {
		return err
	}

	return response.Resp(c, response.Response{
		Messages: response.Messages{"Get track success"},
		Data:     res,
	})
}

// Create godoc
// @Summary      Upload new track
// @Description  Upload new track with metadata and file
// @Tags         Music
// @Accept       multipart/form-data
// @Produce      json
// @Param        title    formData string true  "Track Title"
// @Param        artist   formData string false "Artist Name"
// @Param        album    formData string false "Album Name"
// @Param        duration formData int    false "Duration in seconds"
// @Param        file     formData file   true  "Audio File"
// @Success      201 {object} response.Response
// @Security     Bearer
// @Router       /music [post]
func (_i *trackController) Create(c *fiber.Ctx) error {
	userToken := c.Locals("user").(*jwt.Token)
	claims := userToken.Claims.(*middleware.JWTClaims)

	req := new(request.CreateTrackRequest)
	if err := c.BodyParser(req); err != nil {
		return err
	}

	file, err := c.FormFile("file")
	if err != nil {
		return err
	}

	// Validate file type
	allowedMimeTypes := map[string]bool{
		"audio/mpeg":   true,
		"audio/wav":    true,
		"audio/ogg":    true,
		"audio/flac":   true,
		"audio/x-m4a":  true,
		"audio/mp4":    true,
		"audio/aac":    true,
		"audio/midi":   true,
		"audio/x-midi": true,
		"audio/webm":   true,
	}

	contentType := file.Header.Get("Content-Type")
	if !allowedMimeTypes[contentType] {
		return &response.Error{
			Code:    fiber.StatusBadRequest,
			Message: "File type not allowed. Only audio files are permitted.",
		}
	}

	res, err := _i.trackService.CreateTrack(*req, claims.UserID, file)
	if err != nil {
		return err
	}

	return response.Resp(c, response.Response{
		Messages: response.Messages{"Create track success"},
		Data:     res,
		Code:     fiber.StatusCreated,
	})
}

// Update godoc
// @Summary      Update track metadata
// @Description  Update track title, artist and album
// @Tags         Music
// @Accept       json
// @Produce      json
// @Param        id   path uint64 true "Track ID"
// @Param        body body request.UpdateTrackRequest true "Track Metadata"
// @Success      200 {object} response.Response
// @Security     Bearer
// @Router       /music/{id} [put]
func (_i *trackController) Update(c *fiber.Ctx) error {
	id, err := c.ParamsInt("id")
	if err != nil {
		return err
	}

	userToken := c.Locals("user").(*jwt.Token)
	claims := userToken.Claims.(*middleware.JWTClaims)

	req := new(request.UpdateTrackRequest)
	if err := c.BodyParser(req); err != nil {
		return err
	}

	res, err := _i.trackService.UpdateTrack(uint64(id), *req, claims.UserID)
	if err != nil {
		return err
	}

	return response.Resp(c, response.Response{
		Messages: response.Messages{"Update track success"},
		Data:     res,
	})
}

// Delete godoc
// @Summary      Delete track
// @Description  Delete track metadata and file from storage
// @Tags         Music
// @Accept       json
// @Produce      json
// @Param        id   path uint64 true "Track ID"
// @Success      200 {object} response.Response
// @Security     Bearer
// @Router       /music/{id} [delete]
func (_i *trackController) Delete(c *fiber.Ctx) error {
	id, err := c.ParamsInt("id")
	if err != nil {
		return err
	}

	userToken := c.Locals("user").(*jwt.Token)
	claims := userToken.Claims.(*middleware.JWTClaims)

	err = _i.trackService.DeleteTrack(uint64(id), claims.UserID)
	if err != nil {
		return err
	}

	return response.Resp(c, response.Response{
		Messages: response.Messages{"Delete track success"},
	})
}
