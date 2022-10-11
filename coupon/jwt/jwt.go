package jwt

import (
	"account/custom_error"
	"account/internal"
	"account/log"
	"errors"
	"github.com/dgrijalva/jwt-go"
	"time"
)

type CustomClaims struct {
	jwt.StandardClaims
	ID     int32
	Mobile string
}

type JWT struct {
	SigningKey []byte
}

func NewJWT() *JWT {
	return &JWT{SigningKey: []byte(internal.AppConf.JWTKey.SigningKey)}
}

func (j *JWT) GenerateJWT(claims CustomClaims) (string, error) {
	claims.StandardClaims.ExpiresAt = time.Now().Add(7 * 24 * time.Hour).Unix()
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	signedString, err := token.SignedString(j.SigningKey)
	if err != nil {
		log.Logger.Error(err.Error())
		return "", errors.New(custom_error.TokenGenerateFailed)
	}
	return signedString, nil
}

func (j *JWT) ParseJWT(tokenStr string) (*CustomClaims, error) {
	token, err := jwt.ParseWithClaims(tokenStr, &CustomClaims{}, func(token *jwt.Token) (interface{}, error) {
		return j.SigningKey, nil
	})
	if err != nil {
		log.Logger.Error(err.Error())
	}
	if token != nil {
		claims, ok := token.Claims.(*CustomClaims)
		if token.Valid && ok {
			return claims, nil
		}
		return nil, errors.New(custom_error.TokenInvalid)
	} else {
		return nil, errors.New(custom_error.TokenInvalid)
	}
}

func (j *JWT) RefreshJWT(tokenStr string) (string, error) {
	claims, err := j.ParseJWT(tokenStr)
	if err != nil {
		log.Logger.Error(err.Error())
		return "", err
	}
	claims.StandardClaims.ExpiresAt = time.Now().Add(7 * 24 * time.Hour).Unix()
	customClaims := CustomClaims{
		StandardClaims: claims.StandardClaims,
		ID:             claims.ID,
		Mobile:         claims.Mobile,
	}
	tokenStr, err = j.GenerateJWT(customClaims)
	if err != nil {
		log.Logger.Error(err.Error())
		return "", err
	}
	return tokenStr, nil
}
