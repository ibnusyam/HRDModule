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
    uploadDir := filepath.Join("uploads", "cleaning")

    if _, err := os.Stat(uploadDir); os.IsNotExist(err) {
        if err := os.MkdirAll(uploadDir, os.ModePerm); err != nil {
            return "", err
        }
    }

    src, err := fileHeader.Open()
    if err != nil {
        return "", err
    }
    defer src.Close()

    ext := filepath.Ext(fileHeader.Filename)

    if ext == "" || ext == ".blob" {
        contentType := fileHeader.Header.Get("Content-Type")
        switch contentType {
        case "image/jpeg":
            ext = ".jpg"
        case "image/png":
            ext = ".png"
        case "image/webp":
            ext = ".webp"
        default:
            ext = ".jpg" 
        }
    }

    filename := fmt.Sprintf("%s_%d%s", prefix, time.Now().Unix(), ext)
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
    loc, err := time.LoadLocation("Asia/Jakarta")
    if err != nil {
        loc = time.Local 
    }

    startTime, err := time.ParseInLocation(layout, input.StartTimeStr, loc)
    if err != nil {
        return nil, fmt.Errorf("invalid start time format: %v", err)
    }

    endTime, err := time.ParseInLocation(layout, input.EndTimeStr, loc)
    if err != nil {
        return nil, fmt.Errorf("invalid end time format: %v", err)
    }

    fmt.Println(startTime)
    fmt.Println(endTime)

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

func (s *CleaningLogService) GetFormOptions(siteID int) (map[string]interface{}, error) {
    locTypes, err := s.Repo.GetLocationTypesBySite(siteID)
    if err != nil {
        return nil, err
    }

    locations, err := s.Repo.GetLocationsBySite(siteID)
    if err != nil {
        return nil, err
    }

    return map[string]interface{}{
        "location_types": locTypes,
        "locations":      locations,
    }, nil
}