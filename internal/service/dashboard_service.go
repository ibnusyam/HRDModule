package service

import (
	"HRD/internal/repository"
	"HRD/model"
)

type DashboardService struct {
    Repo *repository.DashboardRepository
}

func NewDashboardService(repo *repository.DashboardRepository) *DashboardService {
    return &DashboardService{Repo: repo}
}

func (s *DashboardService) GetCleanerStats(siteID, month, year int) ([]model.CleanerStat, error) {
    return s.Repo.GetCleanerStatsByMonth(siteID, month, year)
}