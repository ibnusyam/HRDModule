package handlers

import (
	"HRD/internal/service"
	"HRD/model"
	"net/http"
	"strconv"
	"time"

	"github.com/labstack/echo/v4"
)

type DashboardHandler struct {
    Service *service.DashboardService
}

func NewDashboardHandler(svc *service.DashboardService) *DashboardHandler {
    return &DashboardHandler{Service: svc}
}

func (h *DashboardHandler) GetCleanerStats(c echo.Context) error {
    siteID, _ := strconv.Atoi(c.QueryParam("site_id"))
    
    // Default ke waktu sekarang jika parameter tidak dikirim
    currentDate := time.Now()
    month, _ := strconv.Atoi(c.QueryParam("month"))
    year, _ := strconv.Atoi(c.QueryParam("year"))

    if month == 0 { month = int(currentDate.Month()) }
    if year == 0 { year = currentDate.Year() }

    stats, err := h.Service.GetCleanerStats(siteID, month, year)
    if err != nil {
        return c.JSON(http.StatusInternalServerError, echo.Map{"error": err.Error()})
    }

    // Pastikan return data selalu array, meskipun kosong
    if stats == nil {
        stats = []model.CleanerStat{} // Mencegah return null ke frontend
    }

    return c.JSON(http.StatusOK, echo.Map{
        "month": month,
        "year":  year,
        "data":  stats,
    })
}