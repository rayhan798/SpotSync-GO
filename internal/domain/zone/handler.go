package zone

import (
	"net/http"
	"spotsync/internal/domain/zone/dto"
	"strconv"

	"github.com/labstack/echo/v4"
)

type Handler struct {
	svc Service
}

func NewHandler(svc Service) *Handler {
	return &Handler{svc}
}

// CreateZone - POST /api/v1/zones (Admin Only)
func (h *Handler) CreateZone(c echo.Context) error {
	//  Role Verification (Handler Level Check)
	userRole := c.Get("role").(string)
	if userRole != "admin" {
		return c.JSON(http.StatusForbidden, dto.APIResponse{
			Success: false,
			Message: "Access forbidden: admin privileges required",
		})
	}

	var req dto.CreateZoneRequest
	if err := c.Bind(&req); err != nil {
		return c.JSON(http.StatusBadRequest, dto.APIResponse{Success: false, Message: "Invalid payload"})
	}

	if err := c.Validate(&req); err != nil {
		return c.JSON(http.StatusBadRequest, dto.APIResponse{Success: false, Message: err.Error()})
	}

	res, err := h.svc.CreateZone(req)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, dto.APIResponse{Success: false, Message: "Failed to create zone"})
	}

	return c.JSON(http.StatusCreated, dto.APIResponse{
		Success: true,
		Message: "Parking zone created successfully",
		Data:    res,
	})
}

// GetAllZones - GET /api/v1/zones (Public)
func (h *Handler) GetAllZones(c echo.Context) error {
	res, err := h.svc.GetAllZones()
	if err != nil {
		return c.JSON(http.StatusInternalServerError, dto.APIResponse{Success: false, Message: "Failed to fetch zones"})
	}

	return c.JSON(http.StatusOK, dto.APIResponse{
		Success: true,
		Message: "Parking zones retrieved successfully",
		Data:    res,
	})
}

// GetZoneByID - GET /api/v1/zones/:id (Public)
func (h *Handler) GetZoneByID(c echo.Context) error {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		return c.JSON(http.StatusBadRequest, dto.APIResponse{Success: false, Message: "Invalid Zone ID"})
	}

	res, err := h.svc.GetZoneByID(uint(id))
	if err != nil {
		return c.JSON(http.StatusNotFound, dto.APIResponse{Success: false, Message: "Parking zone not found"})
	}

	return c.JSON(http.StatusOK, dto.APIResponse{
		Success: true,
		Message: "Parking zone retrieved successfully",
		Data:    res,
	})
}
