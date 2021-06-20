package service

import (
	"errors"
	"time"
	"yanglu/config"
	"yanglu/def"

	"github.com/dgrijalva/jwt-go"
)

type Claims struct {
	Uid int
	jwt.StandardClaims
}

type TokenService struct {
	uid int
}

func NewTokenService(uid int) *TokenService {
	return &TokenService{
		uid: uid,
	}
}

func NewEmptyTokenService() *TokenService {
	return &TokenService{}
}

func (ts *TokenService) BuildToken() (string, error) {
	secret := def.ApiJwtSecretDev
	if config.IsOnline() {
		secret = def.ApiJwtSecret
	}
	expireTime := 3600*24*7 + time.Now().Unix()

	claims := Claims{
		Uid:            ts.uid,
		StandardClaims: jwt.StandardClaims{ExpiresAt: expireTime},
	}

	tokenClaims := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	token, err := tokenClaims.SignedString([]byte(secret))
	return token, err
}

func (ts *TokenService) CheckToken(token string, secret string) (*Claims, error) {
	claims, err := ts.ParseToken(token, secret)
	if err != nil {
		return claims, err
	}
	if claims.ExpiresAt < time.Now().Unix() {
		return claims, errors.New("token已经过期，请重新登陆")
	}
	return claims, nil
}

func (ts *TokenService) ParseToken(token, secret string) (*Claims, error) {
	tokenClaims, err := jwt.ParseWithClaims(token, &Claims{}, func(token *jwt.Token) (interface{}, error) {
		return []byte(secret), nil
	})

	if tokenClaims != nil {
		if claims, ok := tokenClaims.Claims.(*Claims); ok && tokenClaims.Valid {
			return claims, nil
		}
	}
	return nil, err
}
