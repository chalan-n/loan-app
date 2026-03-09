package handlers

import (
	"loan-app/config"

	"github.com/golang-jwt/jwt/v5"
)

// parseJWTUsername extracts the "username" claim from a JWT token string.
// Returns an empty string if the token is invalid, expired, or empty.
func parseJWTUsername(tokenStr string) string {
	if tokenStr == "" {
		return ""
	}
	token, err := jwt.Parse(tokenStr, func(t *jwt.Token) (interface{}, error) {
		return []byte(config.GetConfig().JWTSecret), nil
	})
	if err != nil || !token.Valid {
		return ""
	}
	if claims, ok := token.Claims.(jwt.MapClaims); ok {
		u, _ := claims["username"].(string)
		return u
	}
	return ""
}

// parseJWTSessionID extracts the "session_id" claim from a JWT token string.
func parseJWTSessionID(tokenStr string) string {
	if tokenStr == "" {
		return ""
	}
	token, err := jwt.Parse(tokenStr, func(t *jwt.Token) (interface{}, error) {
		return []byte(config.GetConfig().JWTSecret), nil
	})
	if err != nil || !token.Valid {
		return ""
	}
	if claims, ok := token.Claims.(jwt.MapClaims); ok {
		s, _ := claims["session_id"].(string)
		return s
	}
	return ""
}
