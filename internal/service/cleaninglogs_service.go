package service

import (
	"HRD/internal/repository"
	"HRD/model"
	"database/sql"
	"fmt"
	"io"
	"mime/multipart"
	"os"
	"path/filepath"
	"time"
)

type CleaningLogService struct {
	Repo *repository.CleaningLogsRepository
}

func NewCleaningLogService(repo *repository.CleaningLogsRepository) *CleaningLogService {
	return &CleaningLogService{Repo: repo}
}

func (s *CleaningLogService) GetAllCleaningsLogs(siteID, locationID, typeID, page, limit int, cleanerName, dateStr string) (*model.PaginatedResponse, error) {
	if page < 1 {
		page = 1
	}
	if limit < 1 {
		limit = 10
	}
	offset := (page - 1) * limit

	totalRecords, err := s.Repo.CountLogs(siteID, locationID, typeID, cleanerName, dateStr)
	if err != nil {
		return nil, fmt.Errorf("failed to count logs: %w", err)
	}

	logs, err := s.Repo.GetAllCleaningsLogs(siteID, locationID, typeID, limit, offset, cleanerName, dateStr)
	if err != nil {
		return nil, fmt.Errorf("service failed to fetch logs: %w", err)
	}

	totalPages := (totalRecords + limit - 1) / limit

	return &model.PaginatedResponse{
		Data: logs,
		Meta: model.PaginationMeta{
			CurrentPage:  page,
			TotalPages:   totalPages,
			TotalRecords: totalRecords,
			Limit:        limit,
		},
	}, nil
}

func saveFile(fileHeader *multipart.FileHeader, prefix string) (string, error) {
	uploadDir := "uploads"
	if _, err := os.Stat(uploadDir); os.IsNotExist(err) {
		os.Mkdir(uploadDir, os.ModePerm)
	}

	src, err := fileHeader.Open()
	if err != nil {
		return "", err
	}
	defer src.Close()

	filename := fmt.Sprintf("%s_%d_%s", prefix, time.Now().Unix(), fileHeader.Filename)
	path := filepath.Join(uploadDir, filename)

	dst, err := os.Create(path)
	if err != nil {
		return "", err
	}
	defer dst.Close()

	if _, err = io.Copy(dst, src); err != nil {
		return "", err
	}

	return path, nil
}

func (s *CleaningLogService) CreateFullLog(input model.CreateFullLogInput, fileBefore, fileAfter *multipart.FileHeader) (*model.CleaningLog, error) {
    layout := "2006-01-02 15:04:05"
    startTime, _ := time.Parse(layout, input.StartTimeStr)
    endTime, _ := time.Parse(layout, input.EndTimeStr)

    pathBefore, err := saveFile(fileBefore, "before")
    if err != nil { return nil, err }

    pathAfter, err := saveFile(fileAfter, "after")
    if err != nil {
        os.Remove(pathBefore)
        return nil, err
    }

    log := &model.CleaningLog{
        CleanerName:    input.CleanerName,
        LocationID:     input.LocationID,     // Menggunakan ID Integer
        LocationTypeID: input.LocationTypeID, // Menggunakan ID Integer
        StartTime:      startTime,
        EndTime:        sql.NullTime{Time: endTime, Valid: true},
        ImageBeforeURL: pathBefore,
        ImageAfterURL:  sql.NullString{String: pathAfter, Valid: true},
        Notes:          sql.NullString{String: input.Notes, Valid: input.Notes != ""},
        SiteID:         input.SiteID,
    }

    if err := s.Repo.CreateFullLog(log); err != nil {
        os.Remove(pathBefore)
        os.Remove(pathAfter)
        return nil, err
    }
    return log, nil
}

func (s *CleaningLogService) GetFormOptions() (map[string]interface{}, error) {
	locTypes, err := s.Repo.GetAllLocationTypes()
	if err != nil {
		return nil, err
	}

	locations, err := s.Repo.GetAllLocations()
	if err != nil {
		return nil, err
	}

	return map[string]interface{}{
		"location_types": locTypes,
		"locations":      locations,
	}, nil
}