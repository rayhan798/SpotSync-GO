package reservation

import (
	"errors"
	"net/http"
	"spotsync/internal/domain/reservation/dto"
	"strconv"

	"github.com/labstack/echo/v4"
)

type Handler struct {
	svc Service
}

func NewHandler(svc Service) *Handler {
	return &Handler{svc}
}

// ReserveSpot - POST /api/v1/reservations
func (h *Handler) ReserveSpot(c echo.Context) error {
	userID := c.Get("userId").(uint)

	var req dto.ReserveRequest
	if err := c.Bind(&req); err != nil {
		return c.JSON(http.StatusBadRequest, dto.APIResponse{Success: false, Message: "Invalid payload"})
	}

	if err := c.Validate(&req); err != nil {
		return c.JSON(http.StatusBadRequest, dto.APIResponse{Success: false, Message: err.Error()})
	}

	res, err := h.svc.MakeReservation(userID, req)
	if err != nil {
		if errors.Is(err, ErrZoneFull) {
			return c.JSON(http.StatusConflict, dto.APIResponse{
				Success: false,
				Message: "Reservation failed",
				Data:    "Selected parking zone is completely full!",
			})
		}
		return c.JSON(http.StatusInternalServerError, dto.APIResponse{Success: false, Message: err.Error()})
	}

	return c.JSON(http.StatusCreated, dto.APIResponse{
		Success: true,
		Message: "Reservation confirmed successfully",
		Data:    res,
	})
}

// GetMyReservations - GET /api/v1/reservations/my-reservations
func (h *Handler) GetMyReservations(c echo.Context) error {
	userID := c.Get("userId").(uint)

	res, err := h.svc.GetMyReservations(userID)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, dto.APIResponse{Success: false, Message: "Failed to retrieve reservations"})
	}

	return c.JSON(http.StatusOK, dto.APIResponse{
		Success: true,
		Message: "My reservations retrieved successfully",
		Data:    res,
	})
}

// CancelReservation - DELETE /api/v1/reservations/:id
func (h *Handler) CancelReservation(c echo.Context) error {
	userID := c.Get("userId").(uint)
	userRole := c.Get("role").(string)

	resID, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		return c.JSON(http.StatusBadRequest, dto.APIResponse{Success: false, Message: "Invalid Reservation ID"})
	}

	err = h.svc.CancelReservation(userID, userRole, uint(resID))
	if err != nil {
		if err.Error() == "unauthorized to cancel this reservation" {
			return c.JSON(http.StatusForbidden, dto.APIResponse{Success: false, Message: err.Error()})
		}
		return c.JSON(http.StatusBadRequest, dto.APIResponse{Success: false, Message: err.Error()})
	}

	return c.JSON(http.StatusOK, dto.APIResponse{
		Success: true,
		Message: "Reservation cancelled successfully",
	})
}

// GetAllReservations - GET /api/v1/reservations (Admin Only)
func (h *Handler) GetAllReservations(c echo.Context) error {
	userRole := c.Get("role").(string)
	if userRole != "admin" {
		return c.JSON(http.StatusForbidden, dto.APIResponse{Success: false, Message: "Access forbidden: admin privileges required"})
	}

	res, err := h.svc.GetAllReservations()
	if err != nil {
		return c.JSON(http.StatusInternalServerError, dto.APIResponse{Success: false, Message: "Failed to fetch all reservations"})
	}

	return c.JSON(http.StatusOK, dto.APIResponse{
		Success: true,
		Message: "All reservations retrieved successfully",
		Data:    res,
	})
}
