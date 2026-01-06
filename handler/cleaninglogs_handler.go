package handler

import (
	"HRD/internal/service"
	"HRD/model"
	"fmt"
	"net/http"
	"strconv"

	"github.com/labstack/echo/v4"
)

type CleaningLogHandler struct {
	Service *service.CleaningLogService
}

func NewCleaningLogHandler(service *service.CleaningLogService) *CleaningLogHandler {
	return &CleaningLogHandler{Service: service}
}

func (h *CleaningLogHandler) GetAllLogs(c echo.Context) error {
	siteID, _ := strconv.Atoi(c.QueryParam("site_id"))
	locationID, _ := strconv.Atoi(c.QueryParam("location_id"))
	typeID, _ := strconv.Atoi(c.QueryParam("type_id"))
	page, _ := strconv.Atoi(c.QueryParam("page"))
	limit, _ := strconv.Atoi(c.QueryParam("limit"))
	cleanerName := c.QueryParam("cleaner_name")
	dateStr := c.QueryParam("date")

	if siteID == 0 {
		return c.JSON(http.StatusBadRequest, echo.Map{"error": "site_id is required"})
	}

	response, err := h.Service.GetAllCleaningsLogs(siteID, locationID, typeID, page, limit, cleanerName, dateStr)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, echo.Map{"error": err.Error()})
	}

	return c.JSON(http.StatusOK, response)
}

func (h *CleaningLogHandler) CreateFullLog(c echo.Context) error {
    // Tangkap data berdasarkan kunci yang terlihat di log konsol browser Anda
    cleanerName := c.FormValue("cleaner_name")
    locationID, _ := strconv.Atoi(c.FormValue("location_name"))      // Menangkap "105"
    locationTypeID, _ := strconv.Atoi(c.FormValue("location_type_name")) // Menangkap "5"
    siteID, _ := strconv.Atoi(c.FormValue("site_id"))
    
    notes := c.FormValue("notes")
    startTimeStr := c.FormValue("start_time")
    endTimeStr := c.FormValue("end_time")
    fmt.Println(cleanerName)
    fmt.Println(locationID)
    fmt.Println(locationTypeID)
    fmt.Println(startTimeStr)
    fmt.Println(endTimeStr)

    if cleanerName == "" || locationID == 0 || locationTypeID == 0 {
        return c.JSON(http.StatusBadRequest, echo.Map{"error": "Data Kosong"})
    }

    fileBefore, errBefore := c.FormFile("image_before")
    fileAfter, errAfter := c.FormFile("image_after")
    if errBefore != nil || errAfter != nil {
        return c.JSON(http.StatusBadRequest, echo.Map{"error": "foto wajib diunggah"})
    }

    input := model.CreateFullLogInput{
        CleanerName:    cleanerName,
        LocationID:     locationID,
        LocationTypeID: locationTypeID,
        Notes:          notes,
        StartTimeStr:   startTimeStr,
        EndTimeStr:     endTimeStr,
        SiteID:         siteID,
    }

    createdLog, err := h.Service.CreateFullLog(input, fileBefore, fileAfter)
    if err != nil {
        return c.JSON(http.StatusBadRequest, echo.Map{"error": err.Error()})
    }

    return c.JSON(http.StatusCreated, echo.Map{"data": createdLog})
}

func (h *CleaningLogHandler) GetFormOptionsHandler(c echo.Context) error {
    // 1. Ambil site_id dari URL (?site_id=1)
    siteIDParam := c.QueryParam("site_id")
    siteID, err := strconv.Atoi(siteIDParam)
    
    // Validasi sederhana jika site_id tidak valid atau kosong
    if err != nil || siteID == 0 {
        return c.JSON(http.StatusBadRequest, echo.Map{"error": "Invalid or missing site_id"})
    }

    // 2. Panggil Service dengan siteID
    options, err := h.Service.GetFormOptions(siteID)
    if err != nil {
        return c.JSON(http.StatusInternalServerError, echo.Map{"error": err.Error()})
    }

    return c.JSON(http.StatusOK, options)
}