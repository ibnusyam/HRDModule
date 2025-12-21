package handler

import (
	"HRD/internal/service"
	"HRD/model"
	"strconv"

	"github.com/labstack/echo/v4"
)

type LocationHandler struct {
	Service *service.LocationService
}

func NewLocationHandler(s *service.LocationService) *LocationHandler {
	return &LocationHandler{Service: s}
}

// ============ TYPE HANDLERS ============

func (h *LocationHandler) GetTypes(c echo.Context) error {
	data, err := h.Service.GetTypes()
	if err != nil { return c.JSON(500, map[string]string{"error": err.Error()}) }
	return c.JSON(200, map[string]interface{}{"data": data})
}

func (h *LocationHandler) CreateType(c echo.Context) error {
	var input model.LocationTypeIt
	if err := c.Bind(&input); err != nil { return c.JSON(400, map[string]string{"error": "Invalid input"}) }
	
	if err := h.Service.CreateType(input); err != nil { return c.JSON(500, map[string]string{"error": err.Error()}) }
	return c.JSON(201, map[string]string{"message": "Type Created"})
}

func (h *LocationHandler) UpdateType(c echo.Context) error {
	id, _ := strconv.Atoi(c.Param("id"))
	var input model.LocationTypeIt
	if err := c.Bind(&input); err != nil { return c.JSON(400, map[string]string{"error": "Invalid input"}) }

	if err := h.Service.UpdateType(id, input); err != nil { return c.JSON(500, map[string]string{"error": err.Error()}) }
	return c.JSON(200, map[string]string{"message": "Type Updated"})
}

func (h *LocationHandler) DeleteType(c echo.Context) error {
	id, _ := strconv.Atoi(c.Param("id"))
	if err := h.Service.DeleteType(id); err != nil { return c.JSON(500, map[string]string{"error": err.Error()}) }
	return c.JSON(200, map[string]string{"message": "Type Deleted"})
}

// ============ LOCATION HANDLERS ============

func (h *LocationHandler) GetLocations(c echo.Context) error {
	data, err := h.Service.GetLocations()
	if err != nil { return c.JSON(500, map[string]string{"error": err.Error()}) }
	return c.JSON(200, map[string]interface{}{"data": data})
}

func (h *LocationHandler) CreateLocation(c echo.Context) error {
	var input model.Location
	if err := c.Bind(&input); err != nil { return c.JSON(400, map[string]string{"error": "Invalid input"}) }

	if err := h.Service.CreateLocation(input); err != nil { return c.JSON(500, map[string]string{"error": err.Error()}) }
	return c.JSON(201, map[string]string{"message": "Location Created"})
}

func (h *LocationHandler) UpdateLocation(c echo.Context) error {
	id, _ := strconv.Atoi(c.Param("id"))
	var input model.Location
	if err := c.Bind(&input); err != nil { return c.JSON(400, map[string]string{"error": "Invalid input"}) }

	if err := h.Service.UpdateLocation(id, input); err != nil { return c.JSON(500, map[string]string{"error": err.Error()}) }
	return c.JSON(200, map[string]string{"message": "Location Updated"})
}

func (h *LocationHandler) DeleteLocation(c echo.Context) error {
	id, _ := strconv.Atoi(c.Param("id"))
	if err := h.Service.DeleteLocation(id); err != nil { return c.JSON(500, map[string]string{"error": err.Error()}) }
	return c.JSON(200, map[string]string{"message": "Location Deleted"})
}