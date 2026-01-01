package repository

import (
	"git.dev.siap.id/kukuhkkh/app-music/app/database/schema"
	"git.dev.siap.id/kukuhkkh/app-music/internal/bootstrap/database"
	"git.dev.siap.id/kukuhkkh/app-music/utils/paginator"
)

type trackRepository struct {
	DB *database.Database
}

//go:generate mockgen -destination=track_repository_mock.go -package=repository . TrackRepository
type TrackRepository interface {
	FindTrackByID(id uint64) (track *schema.Track, err error)
	ListTracks() (tracks []schema.Track, err error)
	PaginateTracks(search string, p *paginator.Pagination) (tracks []schema.Track, pagination *paginator.Pagination, err error)
	CreateTrack(track *schema.Track) (res *schema.Track, err error)
	UpdateTrack(id uint64, track *schema.Track) (res *schema.Track, err error)
	DeleteTrack(id uint64) (err error)
}

func NewTrackRepository(db *database.Database) TrackRepository {
	return &trackRepository{
		DB: db,
	}
}

func (_i *trackRepository) PaginateTracks(search string, p *paginator.Pagination) (tracks []schema.Track, pagination *paginator.Pagination, err error) {
	query := _i.DB.DB.Model(&schema.Track{}).Preload("User")

	if search != "" {
		s := "%" + search + "%"
		query = query.Where("title LIKE ? OR artist LIKE ?", s, s)
	}

	if err = query.Count(&p.Count).Error; err != nil {
		return
	}

	err = query.Offset(p.Offset).Limit(p.Limit).Order("created_at DESC").Find(&tracks).Error

	return tracks, p, err
}

func (_i *trackRepository) FindTrackByID(id uint64) (track *schema.Track, err error) {
	if err := _i.DB.DB.Preload("User").First(&track, id).Error; err != nil {
		return nil, err
	}

	return
}

func (_i *trackRepository) ListTracks() (tracks []schema.Track, err error) {
	if err := _i.DB.DB.Preload("User").Find(&tracks).Error; err != nil {
		return nil, err
	}

	return
}

func (_i *trackRepository) CreateTrack(track *schema.Track) (res *schema.Track, err error) {
	if err := _i.DB.DB.Create(&track).Error; err != nil {
		return nil, err
	}

	return track, nil
}

func (_i *trackRepository) UpdateTrack(id uint64, track *schema.Track) (res *schema.Track, err error) {
	if err := _i.DB.DB.Model(&schema.Track{}).Where("id = ?", id).Updates(track).Error; err != nil {
		return nil, err
	}

	return track, nil
}

func (_i *trackRepository) DeleteTrack(id uint64) (err error) {
	if err := _i.DB.DB.Delete(&schema.Track{}, id).Error; err != nil {
		return err
	}

	return nil
}
