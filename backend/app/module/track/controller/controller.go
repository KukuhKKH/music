package controller

import "git.dev.siap.id/kukuhkkh/app-music/app/module/track/service"

type Controller struct {
	Track TrackController
}

func NewController(trackService service.TrackService) *Controller {
	return &Controller{
		Track: NewTrackController(trackService),
	}
}
