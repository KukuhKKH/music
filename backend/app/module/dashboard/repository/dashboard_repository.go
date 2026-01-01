package repository

import (
	"git.dev.siap.id/kukuhkkh/app-music/app/database/schema"
	"git.dev.siap.id/kukuhkkh/app-music/internal/bootstrap/database"
)

type DashboardRepository interface {
	GetSummary() (totalSongs int64, totalSize int64, lastUpload string, err error)
}

type dashboardRepository struct {
	DB *database.Database
}

func NewDashboardRepository(db *database.Database) DashboardRepository {
	return &dashboardRepository{
		DB: db,
	}
}

func (r *dashboardRepository) GetSummary() (totalSongs int64, totalSize int64, lastUpload string, err error) {
	// Total Songs
	err = r.DB.DB.Model(&schema.Track{}).Count(&totalSongs).Error
	if err != nil {
		return
	}

	// Total Size
	err = r.DB.DB.Model(&schema.Track{}).Select("IFNULL(SUM(file_size), 0)").Row().Scan(&totalSize)
	if err != nil {
		return
	}

	// Last Upload
	var track schema.Track
	err = r.DB.DB.Order("created_at desc").First(&track).Error
	if err == nil {
		lastUpload = track.CreatedAt.Format("2006-01-02 15:04:05")
	} else {
		lastUpload = "-"
		err = nil // Reset error if no tracks found
	}

	return
}
