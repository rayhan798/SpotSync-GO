package reservation

import (
	"github.com/labstack/echo/v4"
	"gorm.io/gorm"
)

func RegisterModule(e *echo.Echo, db *gorm.DB, authMiddleware echo.MiddlewareFunc) {
	repo := NewRepository(db)
	svc := NewService(repo)
	handler := NewHandler(svc)

	rGroup := e.Group("/api/v1/reservations", authMiddleware)

	rGroup.POST("", handler.ReserveSpot)
	rGroup.GET("/my-reservations", handler.GetMyReservations)
	rGroup.DELETE("/:id", handler.CancelReservation)
	rGroup.GET("", handler.GetAllReservations)
}
