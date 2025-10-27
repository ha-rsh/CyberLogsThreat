package auth

import (
	"errors"
	"time"
	"github.com/golang-jwt/jwt/v5"
)

var jwtSecret = []byte("secret")

type Claims struct {
	UserID   					string 		`json:"userId"`
	Username 					string 		`json:"username"`
	Role     					string 		`json:"role"`
	jwt.RegisteredClaims
}

func GenerateToken(userID, username, role string) (string, error) {
	expirationTime := time.Now().Add(2 * time.Hour)
	
	claims := &Claims{
		UserID:   userID,
		Username: username,
		Role:     role,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(expirationTime),
			IssuedAt:  jwt.NewNumericDate(time.Now()),
			Issuer:    "cyber-threat-detection",
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString(jwtSecret)
}

func ValidateToken(tokenString string) (*Claims, error) {
	claims := &Claims{}

	token, err := jwt.ParseWithClaims(tokenString, claims, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, errors.New("invalid signing method")
		}
		return jwtSecret, nil
	})

	if err != nil {
		return nil, err
	}

	if !token.Valid {
		return nil, errors.New("invalid token")
	}

	return claims, nil
}

func RefreshToken(tokenString string) (string, error) {
	claims, err := ValidateToken(tokenString)
	if err != nil {
		return "", err
	}

	return GenerateToken(claims.UserID, claims.Username, claims.Role)
}

func SetJWTSecret(secret string) {
	jwtSecret = []byte(secret)
}