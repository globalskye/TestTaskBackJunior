package models

import (
	"time"
)

type User struct {
	GUID         string    `json:"guid,omitempty" bson:"guid,omitempty"`
	RefreshToken string    `json:"refreshToken,omitempty" bson:"refreshToken,omitempty"`
	ExpiresAt    time.Time `json:"expiresAt,omitempty" bson:"expiresAT,omitempty"`
}

type AuthToken struct {
	RefreshToken string `json:"refresh_token"`
	AuthToken    string `json:"access_token"`
}

type TokenResponse struct {
	AccessToken  string `json:"access_token"`
	TokenType    string `json:"token_type"`
	ExpiresAt    int64  `json:"expires_at"`
	RefreshToken string `json:"refresh_token"`
}

type RefreshToken struct {
	RefreshToken string `json:"refresh_token"`
}
