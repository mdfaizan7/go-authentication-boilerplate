package util

import (
	db "go-authentication-boilerplate/database"
	"go-authentication-boilerplate/models"
	"time"

	"github.com/dgrijalva/jwt-go"
)

var jwtKey = []byte(db.PRIVKEY)

// GenerateClaims generates jwt token
func GenerateClaims(u *models.User) (*models.Claims, string) {

	t := time.Now()
	claim := &models.Claims{
		StandardClaims: jwt.StandardClaims{
			Issuer:    u.UUID.String(),
			ExpiresAt: t.Add(30 * time.Minute).Unix(),
			Subject:   "access_token",
			IssuedAt:  t.Unix(),
		},
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claim)
	tokenString, err := token.SignedString(jwtKey)
	if err != nil {
		panic(err)
	}

	return claim, tokenString
}

// GenerateRefreshClaims generates refresh tokens
func GenerateRefreshClaims(cl *models.Claims) string {
	result := db.DB.Where(&models.Claims{StandardClaims: jwt.StandardClaims{Issuer: cl.Issuer}}).Find(&models.Claims{})
	// checking the number of refresh tokens stored.
	// If the number is higher than 3, remove all the refresh tokens and leave only new one.
	if result.RowsAffected > 3 {
		db.DB.Where(&models.Claims{StandardClaims: jwt.StandardClaims{Issuer: cl.Issuer}}).Delete(&models.Claims{})
	}

	t := time.Now()
	refreshClaim := &models.Claims{
		StandardClaims: jwt.StandardClaims{
			Issuer:    cl.Issuer,
			ExpiresAt: t.Add(30 * 24 * time.Hour).Unix(),
			Subject:   "refresh_token",
			IssuedAt:  t.Unix(),
		},
	}

	// create a claim on DB
	db.DB.Create(&refreshClaim)

	refreshToken := jwt.NewWithClaims(jwt.SigningMethodHS256, refreshClaim)
	refreshTokenString, err := refreshToken.SignedString(jwtKey)
	if err != nil {
		panic(err)
	}

	return refreshTokenString
}
