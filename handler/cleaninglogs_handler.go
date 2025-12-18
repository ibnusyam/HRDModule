package handlers

import (
	"HRD/internal/service"
	"HRD/model"
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

    if cleanerName == "" || locationID == 0 || locationTypeID == 0 {
        return c.JSON(http.StatusBadRequest, echo.Map{"error": "data tidak lengkap"})
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

func (h *CleaningLogHandler) GetFormOptions(c echo.Context) error {
	options, err := h.Service.GetFormOptions()
	if err != nil {
		return c.JSON(http.StatusInternalServerError, echo.Map{"error": err.Error()})
	}
	return c.JSON(http.StatusOK, echo.Map{"data": options})
}