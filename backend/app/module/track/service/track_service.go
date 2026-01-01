package service

import (
	"fmt"
	"mime/multipart"
	"path/filepath"
	"time"

	"git.dev.siap.id/kukuhkkh/app-music/app/database/schema"
	"git.dev.siap.id/kukuhkkh/app-music/app/module/track/repository"
	"git.dev.siap.id/kukuhkkh/app-music/app/module/track/request"
	"git.dev.siap.id/kukuhkkh/app-music/utils/helpers"
	"git.dev.siap.id/kukuhkkh/app-music/utils/paginator"
	"git.dev.siap.id/kukuhkkh/app-music/utils/storage"
)

type trackService struct {
	repo    repository.TrackRepository
	storage storage.Storage
}

type TrackService interface {
	GetPaginatedTracks(search string, p *paginator.Pagination) (tracks []schema.Track, pagination *paginator.Pagination, err error)
	CreateTrack(req request.CreateTrackRequest, userID uint64, fileHeader *multipart.FileHeader) (track *schema.Track, err error)
}

func NewTrackService(repo repository.TrackRepository, storage storage.Storage) TrackService {
	return &trackService{
		repo:    repo,
		storage: storage,
	}
}

func (s *trackService) GetPaginatedTracks(search string, p *paginator.Pagination) (tracks []schema.Track, pagination *paginator.Pagination, err error) {
	return s.repo.PaginateTracks(search, p)
}

func (s *trackService) CreateTrack(req request.CreateTrackRequest, userID uint64, fileHeader *multipart.FileHeader) (track *schema.Track, err error) {
	file, err := fileHeader.Open()
	if err != nil {
		return nil, err
	}

	defer file.Close()

	// Generate a unique filename
	ext := filepath.Ext(fileHeader.Filename)
	storageFilename := fmt.Sprintf("%d_%s%s", time.Now().UnixNano(), helpers.Slug(req.Title), ext)

	// Stream to Storage
	_, err = s.storage.Upload(storageFilename, file)
	if err != nil {
		return nil, err
	}

	// Create DB Record
	newTrack := &schema.Track{
		UserID:           userID,
		Title:            req.Title,
		Artist:           req.Artist,
		Album:            &req.Album,
		Duration:         req.Duration,
		StorageFilename:  storageFilename,
		OriginalFilename: fileHeader.Filename,
		FileSize:         fileHeader.Size,
		MimeType:         fileHeader.Header.Get("Content-Type"),
	}

	return s.repo.CreateTrack(newTrack)
}
