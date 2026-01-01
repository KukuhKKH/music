package service

import (
	"git.dev.siap.id/kukuhkkh/app-music/app/database/schema"
	"git.dev.siap.id/kukuhkkh/app-music/app/module/track/repository"
	"git.dev.siap.id/kukuhkkh/app-music/utils/paginator"
)

type trackService struct {
	repo repository.TrackRepository
}

type TrackService interface {
	GetPaginatedTracks(search string, p *paginator.Pagination) (tracks []schema.Track, pagination *paginator.Pagination, err error)
}

func NewTrackService(repo repository.TrackRepository) TrackService {
	return &trackService{
		repo: repo,
	}
}

func (s *trackService) GetPaginatedTracks(search string, p *paginator.Pagination) (tracks []schema.Track, pagination *paginator.Pagination, err error) {
	return s.repo.PaginateTracks(search, p)
}
