package helper

import (
	"encoding/base64"
	"errors"
	"github.com/dgrijalva/jwt-go"
	"github.com/google/uuid"
	"jwttask/config"
	"jwttask/models"
	"log"
	"time"
)

type userClaim struct {
	Token  string
	Claims jwt.StandardClaims
}

func GenerateTokens(user models.User) models.AuthToken {
	expTime, err := time.ParseDuration(config.Conf.AccessTokenTime)
	if err != nil {
		log.Println(err)
	}

	claims := jwt.StandardClaims{
		ExpiresAt: time.Now().Add(expTime).UTC().Unix(),
		Id:        user.GUID,
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS512, claims)
	accessToken, err := token.SignedString([]byte(config.Conf.SecretKey))
	if err != nil {
		log.Fatalln(err)
	}

	refreshToken := uuid.New()

	refreshTokenBase64 := base64.StdEncoding.EncodeToString([]byte(refreshToken.String()))

	tokens := models.AuthToken{
		RefreshToken: refreshTokenBase64,
		AuthToken:    accessToken,
	}
	return tokens

}

func ParseAccessToken(tokenString string) (*userClaim, error) {
	type claims struct {
		jwt.StandardClaims
	}

	token, err := jwt.ParseWithClaims(
		tokenString,
		&claims{},
		func(token *jwt.Token) (interface{}, error) {
			return []byte(config.Conf.SecretKey), nil
		},
	)
	if err != nil {
		v, _ := err.(*jwt.ValidationError)

		if v.Errors != jwt.ValidationErrorExpired {
			return &userClaim{}, err
		}
	}

	tokenData, ok := token.Claims.(*claims)
	if !ok {
		return &userClaim{}, errors.New(("Token is invalid"))
	}

	return &userClaim{
		Token:  tokenString,
		Claims: tokenData.StandardClaims,
	}, nil

}

func IsExpired(expiresAt int64) bool {
	now := time.Now().UTC()
	exp := time.Unix(expiresAt, 0)
	return now.After(exp)
}

func IsValidGuid(guid string) bool {
	if _, err := uuid.Parse(guid); err != nil {
		return false
	}

	return true
}
