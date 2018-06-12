package auth

import (
	"fmt"
	"time"

	jwt "github.com/dgrijalva/jwt-go"
)

// GenerateToken generates a token with the user id stored.
func GenerateToken(userID int, secret string) (string, error) {
	token := jwt.New(jwt.SigningMethodHS256)
	// token.Claims["user_id"] = u
	c := jwt.MapClaims{
		"user_id": userID,
		"exp":     time.Now().Add(time.Hour * 24 * 10),
	}
	token.Claims = c
	r, err := token.SignedString([]byte(secret))
	if err != nil {
		return "", err
	}
	return r, nil
}

// ValidateToken checks whether the token is valid and returns the user id.
func ValidateToken(tokenString, secret string) (int, error) {
	token, err := jwt.Parse(tokenString, func(t *jwt.Token) (interface{}, error) {
		// Check whether algorith is the same
		_, ok := t.Method.(*jwt.SigningMethodHMAC)
		if !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", t.Header["alg"])
		}
		// If everything ok, return secret
		return []byte(secret), nil
	})

	if err != nil {
		return 0, err
	}

	claims, ok := token.Claims.(jwt.MapClaims)
	if !ok || !token.Valid {
		return 0, fmt.Errorf("token is invalid")
	}

	// Temporary fix for compatility, since originally this was a string
	id, ok := claims["user_id"].(float64)
	if !ok {
		return 0, fmt.Errorf("resetting browser since user_id is now an int")
	}

	return int(id), nil
}
