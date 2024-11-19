package utils

import (
	"time"

	"github.com/danielpnjt/speed-engine/internal/config"
	"github.com/golang-jwt/jwt"
)

type Claims struct {
	ID       int
	Username string
	jwt.StandardClaims
}

type JWTClaimsData struct {
	ID       int    `json:"id"`
	Username string `json:"username"`
}

func JwtSign(id int, username string) (token string, exp int64, err error) {
	secret := []byte(config.GetString("jwt"))
	exp = time.Now().Add(7 * 24 * time.Hour).Unix()
	claims := &Claims{
		id,
		username,
		jwt.StandardClaims{
			Issuer:    "USER_SERVICE",
			ExpiresAt: exp, // * 7 days
		},
	}
	token, err = jwt.NewWithClaims(jwt.SigningMethodHS256, claims).SignedString(secret)
	return
}

func JwtVerify(token string) (claims *Claims, err error) {
	secret := []byte(config.GetString("jwt"))
	claims = &Claims{}
	_, err = jwt.ParseWithClaims(token, claims, func(t *jwt.Token) (interface{}, error) {
		return secret, nil
	})

	return
}
