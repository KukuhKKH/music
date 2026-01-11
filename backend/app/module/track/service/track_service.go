package service

import (
	"context"
	"fmt"
	"log"
	"mime/multipart"
	"path/filepath"
	"time"

	"git.dev.siap.id/kukuhkkh/app-music/app/database/schema"
	"git.dev.siap.id/kukuhkkh/app-music/app/module/track/repository"
	"git.dev.siap.id/kukuhkkh/app-music/app/module/track/request"
	"git.dev.siap.id/kukuhkkh/app-music/app/module/track/response"
	"git.dev.siap.id/kukuhkkh/app-music/utils/helpers"
	"git.dev.siap.id/kukuhkkh/app-music/utils/paginator"
	"git.dev.siap.id/kukuhkkh/app-music/utils/storage"
)

type trackService struct {
	repo    repository.TrackRepository
	storage storage.Storage
}

type TrackService interface {
	GetPaginatedTracks(search string, p *paginator.Pagination) (tracks []response.TrackResponse, pagination *paginator.Pagination, err error)
	GetTrackByID(id uint64) (track *response.TrackResponse, err error)
	CreateTrack(ctx context.Context, req request.CreateTrackRequest, userID uint64, fileHeader *multipart.FileHeader) (track *response.TrackResponse, err error)
	UpdateTrack(id uint64, req request.UpdateTrackRequest, userID uint64) (track *response.TrackResponse, err error)
	DeleteTrack(id uint64, userID uint64) (err error)
}

func NewTrackService(repo repository.TrackRepository, storage storage.Storage) TrackService {
	return &trackService{
		repo:    repo,
		storage: storage,
	}
}

func (s *trackService) GetPaginatedTracks(search string, p *paginator.Pagination) (tracks []response.TrackResponse, pagination *paginator.Pagination, err error) {
	schemaTracks, p, err := s.repo.PaginateTracks(search, p)
	if err != nil {
		return nil, p, err
	}

	return response.FromTrackListSchema(schemaTracks, s.storage), p, nil
}

func (s *trackService) GetTrackByID(id uint64) (track *response.TrackResponse, err error) {
	schemaTrack, err := s.repo.FindTrackByID(id)
	if err != nil {
		return nil, err
	}

	res := response.FromTrackSchema(*schemaTrack, s.storage.GetURL(schemaTrack.StorageFilename))
	return &res, nil
}

func (s *trackService) CreateTrack(ctx context.Context, req request.CreateTrackRequest, userID uint64, fileHeader *multipart.FileHeader) (track *response.TrackResponse, err error) {
	start := time.Now()
	log.Printf("[track] create start user=%d title=%q size=%d ct=%q",
		userID, req.Title, fileHeader.Size, fileHeader.Header.Get("Content-Type"))

	file, err := fileHeader.Open()
	if err != nil {
		return nil, err
	}

	defer file.Close()

	ext := filepath.Ext(fileHeader.Filename)
	storageFilename := fmt.Sprintf("%d_%s%s", time.Now().UnixNano(), helpers.Slug(req.Title), ext)

	// Hard timeout agar tidak menggantung sampai Traefik timeout
	uploadCtx, cancel := context.WithTimeout(ctx, 2*time.Minute)
	defer cancel()

	log.Printf("[track] upload to storage start name=%s", storageFilename)
	if _, err := s.storage.Upload(uploadCtx, storageFilename, file); err != nil {
		log.Printf("[track] upload failed err=%v dur=%s", err, time.Since(start))
		return nil, err
	}
	log.Printf("[track] upload to storage done dur=%s", time.Since(start))

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

	res, err := s.repo.CreateTrack(newTrack)
	if err != nil {
		return nil, err
	}

	trackRes := response.FromTrackSchema(*res, s.storage.GetURL(res.StorageFilename))
	log.Printf("[track] create success id=%d total_dur=%s", res.ID, time.Since(start))

	return &trackRes, nil
}

func (s *trackService) UpdateTrack(id uint64, req request.UpdateTrackRequest, userID uint64) (track *response.TrackResponse, err error) {
	// Check if track exists and user is owner
	existingTrack, err := s.repo.FindTrackByID(id)
	if err != nil {
		return nil, err
	}

	if existingTrack.UserID != userID {
		return nil, fmt.Errorf("you don't have permission to update this track")
	}

	// Update fields
	existingTrack.Title = req.Title
	existingTrack.Artist = req.Artist
	existingTrack.Album = &req.Album

	res, err := s.repo.UpdateTrack(id, existingTrack)
	if err != nil {
		return nil, err
	}

	trackRes := response.FromTrackSchema(*res, s.storage.GetURL(res.StorageFilename))
	return &trackRes, nil
}

func (s *trackService) DeleteTrack(id uint64, userID uint64) (err error) {
	// Check if track exists and user is owner
	existingTrack, err := s.repo.FindTrackByID(id)
	if err != nil {
		return err
	}

	if existingTrack.UserID != userID {
		return fmt.Errorf("you don't have permission to delete this track")
	}

	// Delete file from storage
	err = s.storage.Delete(existingTrack.StorageFilename)
	if err != nil {
		// Log the error but continue to delete the DB record if the file is already gone
		fmt.Printf("Warning: failed to delete file from storage: %v\n", err)
	}

	// Delete DB Record
	return s.repo.DeleteTrack(id)
}
