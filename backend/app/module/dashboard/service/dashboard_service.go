package service

import (
	"fmt"

	"git.dev.siap.id/kukuhkkh/app-music/app/module/dashboard/repository"
	"git.dev.siap.id/kukuhkkh/app-music/app/module/dashboard/response"
)

type DashboardService interface {
	GetSummary() (*response.DashboardSummaryResponse, error)
}

type dashboardService struct {
	repo repository.DashboardRepository
}

func NewDashboardService(repo repository.DashboardRepository) DashboardService {
	return &dashboardService{
		repo: repo,
	}
}

func (s *dashboardService) GetSummary() (*response.DashboardSummaryResponse, error) {
	totalSongs, totalSize, lastUpload, err := s.repo.GetSummary()
	if err != nil {
		return nil, err
	}

	return &response.DashboardSummaryResponse{
		TotalSongs: totalSongs,
		TotalSize:  formatSize(totalSize),
		LastUpload: lastUpload,
	}, nil
}

func formatSize(size int64) string {
	const unit = 1024
	if size < unit {
		return fmt.Sprintf("%d B", size)
	}

	div, exp := int64(unit), 0
	for n := size / unit; n >= unit; n /= unit {
		div *= unit
		exp++
	}

	return fmt.Sprintf("%.2f %cB", float64(size)/float64(div), "KMGTPE"[exp])
}
