package middlewares

import (
	"net/http"
	"spotsync/internal/auth"
	"strings"

	"github.com/labstack/echo/v4"
)

func AuthMiddleware(jwtService auth.JWTService) echo.MiddlewareFunc {
	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {

			// extract token from authorization header
			authHeader := c.Request().Header.Get("Authorization")
			if authHeader == "" {
				return c.JSON(http.StatusUnauthorized, map[string]string{
					"error": "Missing authorization header",
				})
			}

			// check bearer scheme
			parts := strings.Split(authHeader, " ")

			if len(parts) != 2 || parts[0] != "Bearer" {
				return c.JSON(http.StatusUnauthorized, map[string]string{
					"error": "invalid authorization header format",
				})
			}

			tokenString := parts[1]

			// validate token
			claims, err := jwtService.ValidateToken(tokenString)
			if err != nil {
				return c.JSON(http.StatusUnauthorized, map[string]string{
					"error": "invalid or expired token",
				})
			}

			// 🔒 store user info in context for handlers
			c.Set("userId", claims.UserID)
			c.Set("email", claims.Email)
			c.Set("name", claims.Name)
			c.Set("role", claims.Role)
			return next(c)
		}
	}
}
