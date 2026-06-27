package user

import (
	"errors"
	"net/http"
	"spotsync/internal/domain/user/dto"
	"spotsync/internal/httpresponse"

	"github.com/labstack/echo/v4"
)

type Handler struct {
	svc Service
}

func NewHandler(svc Service) *Handler {
	return &Handler{
		svc: svc,
	}
}

// CreateUser - Handler: Create a new user (sign-up or registration)
func (h *Handler) CreateUser(c echo.Context) error {
	var req dto.CreateRequest

	if err := c.Bind(&req); err != nil {
		return c.JSON(http.StatusBadRequest, httpresponse.Error{
			Code:    http.StatusBadRequest,
			Message: "Invalid request payload",
			Details: err.Error(),
		})
	}

	if err := c.Validate(&req); err != nil {
		return c.JSON(http.StatusBadRequest, httpresponse.Error{
			Code:    http.StatusBadRequest,
			Message: "Validation failed",
			Details: err.Error(),
		})
	}

	response, err := h.svc.CreateUser(req)
	if err != nil {
		if err.Error() == "admin already exists" {
			return c.JSON(http.StatusConflict, httpresponse.Error{
				Code:    http.StatusConflict,
				Message: "Registration failed",
				Details: "Admin already exists! Only one admin is allowed in the system.",
			})
		}

		// Check for specific error types and return appropriate HTTP status codes
		if errors.Is(err, ErrorAlreadyExist) {
			return c.JSON(http.StatusConflict, httpresponse.Error{
				Code:    http.StatusConflict,
				Message: "Failed to create User",
				Details: err.Error(),
			})
		}

		return c.JSON(http.StatusInternalServerError, httpresponse.Error{
			Code:    http.StatusInternalServerError,
			Message: "Failed to create user",
			Details: err.Error(),
		})
	}

	return c.JSON(http.StatusCreated, response)
}

// LoginUser
func (h *Handler) LoginUser(c echo.Context) error {
	var req dto.LoginRequest // আপনার service.go এর ইনপুট টাইপ অনুযায়ী LoginRequest ব্যবহার করা হয়েছে

	if err := c.Bind(&req); err != nil {
		return c.JSON(http.StatusBadRequest, httpresponse.Error{
			Code:    http.StatusBadRequest,
			Message: "Invalid request payload",
			Details: err.Error(),
		})
	}

	if err := c.Validate(&req); err != nil {
		return c.JSON(http.StatusBadRequest, httpresponse.Error{
			Code:    http.StatusBadRequest,
			Message: "Validation failed",
			Details: err.Error(),
		})
	}

	// আপনার সার্ভিস লেয়ারের সঠিক মেথড LoginUser কল করা হয়েছে
	response, err := h.svc.LoginUser(req)
	if err != nil {
		// service.go তে থাকা ErrInvalidCredentials এররের সাথে ম্যাচ করা হয়েছে
		if errors.Is(err, ErrInvalidCredentials) {
			return c.JSON(http.StatusUnauthorized, httpresponse.Error{
				Code:    http.StatusUnauthorized,
				Message: "Cannot login user",
				Details: err.Error(),
			})
		}

		return c.JSON(http.StatusInternalServerError, httpresponse.Error{
			Code:    http.StatusInternalServerError,
			Message: "Failed to login user",
			Details: err.Error(),
		})
	}

	return c.JSON(http.StatusOK, response)
}

// GetMe
func (h *Handler) GetMe(c echo.Context) error {
	userID, ok := c.Get("user_id").(uint)
	if !ok {
		return c.JSON(http.StatusUnauthorized, httpresponse.Error{
			Code:    http.StatusUnauthorized,
			Message: "Cannot get user information",
			Details: "missing user id in context",
		})
	}

	email, _ := c.Get("user_email").(string)
	name, _ := c.Get("user_name").(string)

	return c.JSON(http.StatusOK, map[string]interface{}{
		"id":    userID,
		"name":  name,
		"email": email,
	})
}
