package server

import (
	"fmt"
	"net/http"
	"spotsync/internal/auth"
	"spotsync/internal/config"

	"spotsync/internal/domain/reservation"
	"spotsync/internal/domain/user"
	"spotsync/internal/domain/zone"
	"spotsync/internal/middlewares"

	"github.com/go-playground/validator/v10"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	"gorm.io/gorm"
)

type CustomValidator struct {
	validator *validator.Validate
}

func (cv *CustomValidator) Validate(i any) error {
	if err := cv.validator.Struct(i); err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, err.Error())
	}
	return nil
}

func Start(db *gorm.DB, cfg *config.Config) {
	db.AutoMigrate(
		&user.User{},
		&zone.ParkingZone{},
		&reservation.Reservation{},
	)

	e := echo.New()
	e.Validator = &CustomValidator{validator: validator.New()}
	e.Use(middleware.Logger())

	e.GET("/", func(c echo.Context) error {
		return c.String(http.StatusOK, "Welcome to SpotSync API Server!")
	})

	// routes
	user.RegisterRoutes(e, db, cfg)

	// create jwt service and pass auth middleware to protected modules
	jwtSvc := auth.NewJWTService(cfg.JwtSecret)

	zone.RegisterModule(e, db, middlewares.AuthMiddleware(jwtSvc))
	reservation.RegisterModule(e, db, middlewares.AuthMiddleware(jwtSvc))

	port := fmt.Sprintf(":%s", cfg.Port)
	if err := e.Start(port); err != nil {
		e.Logger.Error("failed to start server", "error", err)
	}
}
