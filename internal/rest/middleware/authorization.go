package middleware

import (
	"fmt"
	"net/http"
	"santapan_payment_service/pkg/json"
	"strings"

	"github.com/golang-jwt/jwt/v5"
	"github.com/labstack/echo/v4"
)

var jwtSecret = []byte("SANTAPANSECRET") // Replace with your secret key

// JWTMiddleware is a middleware for validating JWT tokens
func AuthMiddleware(next echo.HandlerFunc) echo.HandlerFunc {
	return func(c echo.Context) error {
		// Get the token from the Authorization header
		tokenString := c.Request().Header.Get("Authorization")
		if tokenString == "" {
			return json.Response(c, http.StatusUnauthorized, false, "Missing Authorization header", nil)
		}

		// Remove "Bearer " prefix if present
		if strings.HasPrefix(tokenString, "Bearer ") {
			tokenString = strings.TrimPrefix(tokenString, "Bearer ")
		}

		// Parse and validate the token
		token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
			if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
				return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
			}
			return jwtSecret, nil
		})

		if err != nil || !token.Valid {
			return json.Response(c, http.StatusUnauthorized, false, "Invalid or expired token", nil)
		}

		// Extract claims
		claims, ok := token.Claims.(jwt.MapClaims)
		if !ok {
			return json.Response(c, http.StatusUnauthorized, false, "Invalid token claims", nil)
		}

		// Extract the user ID (sub) from the claims
		userID, ok := claims["sub"].(float64) // sub is usually a float64 in JWT claims
		if !ok {
			return json.Response(c, http.StatusUnauthorized, false, "Invalid token claims", nil)
		}

		// Store user ID in the context
		c.Set("userID", int64(userID))

		return next(c) // Call the next handler
	}
}
