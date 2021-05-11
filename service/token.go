package service

import (
	"time"
	"yanglu/config"
	"yanglu/def"

	"github.com/dgrijalva/jwt-go"
)

type TokenService struct {
	uid int
}

func NewTokenService(uid int) *TokenService {
	return &TokenService{
		uid: uid,
	}
}

func (ts *TokenService) BuildToken() (string, error) {
	secret := def.ApiJwtSecretDev
	if config.IsOnline() {
		secret = def.ApiJwtSecret
	}
	expireTime := 3600*24*7 + time.Now().Unix()

	claims := struct {
		Uid int
		jwt.StandardClaims
	}{
		Uid:            ts.uid,
		StandardClaims: jwt.StandardClaims{ExpiresAt: expireTime},
	}

	tokenClaims := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	token, err := tokenClaims.SignedString([]byte(secret))
	return token, err
}
