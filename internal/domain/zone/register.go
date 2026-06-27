package zone

import (
	"github.com/labstack/echo/v4"
	"gorm.io/gorm"
)

// RegisterModule initializes all layer dependencies and registers routes
func RegisterModule(e *echo.Echo, db *gorm.DB, authMiddleware echo.MiddlewareFunc) {
	repo := NewRepository(db)
	svc := NewService(repo)
	handler := NewHandler(svc)

	// API Route Grouping V1
	v1 := e.Group("/api/v1")

	// Public Routes
	v1.GET("/zones", handler.GetAllZones)
	v1.GET("/zones/:id", handler.GetZoneByID)

	// Protected Routes (Requires JWT authMiddleware)
	v1.POST("/zones", handler.CreateZone, authMiddleware)
}
