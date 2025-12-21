package service

import (
	"HRD/internal/repository"
	"HRD/model"
)

type LocationService struct {
	Repo *repository.LocationRepository
}

func NewLocationService(repo *repository.LocationRepository) *LocationService {
	return &LocationService{Repo: repo}
}

// --- TYPE ---
func (s *LocationService) GetTypes() ([]model.LocationTypeIt, error) {
	return s.Repo.GetLocationTypes()
}
func (s *LocationService) CreateType(input model.LocationTypeIt) error {
	return s.Repo.CreateLocationType(input)
}
func (s *LocationService) UpdateType(id int, input model.LocationTypeIt) error {
	return s.Repo.UpdateLocationType(id, input)
}
func (s *LocationService) DeleteType(id int) error {
	return s.Repo.DeleteLocationType(id)
}

// --- LOCATION ---
func (s *LocationService) GetLocations() ([]model.LocationIt, error) {
	return s.Repo.GetLocations()
}
func (s *LocationService) CreateLocation(input model.Location) error {
	return s.Repo.CreateLocation(input)
}
func (s *LocationService) UpdateLocation(id int, input model.Location) error {
	return s.Repo.UpdateLocation(id, input)
}
func (s *LocationService) DeleteLocation(id int) error {
	return s.Repo.DeleteLocation(id)
}