package handlers

import (
	"encoding/json"
	"jwttask/config"
	"jwttask/db"
	"jwttask/helper"
	"strings"

	"jwttask/models"
	user2 "jwttask/user"
	"net/http"
	"time"
)

var Collection = db.ConnectDB()

func Register(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	guid := r.URL.Query().Get("guid")

	ok := user2.IsValidGuid(guid)
	if !ok {
		w.Write([]byte("Bad GUID"))
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	var user models.User

	user.GUID = guid

	_, ok = db.GetUserByGUID(user.GUID)

	if ok {
		w.Write([]byte("User alredy exist"))
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	tokens := user2.GenerateTokens(user)

	ok = db.InsertUserByGUID(user, tokens)

	if !ok {
		w.Write([]byte("Error adding to database"))
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	expTime, err := time.ParseDuration(config.Conf.AccessTokenTime)
	if err != nil {
		w.Write([]byte("refresh token expired"))
		w.WriteHeader(http.StatusUnauthorized)
		return
	}

	response, err := json.Marshal(&models.TokenResponse{
		AccessToken:  tokens.AuthToken,
		TokenType:    "Bearer",
		ExpiresAt:    time.Now().Add(expTime).UTC().Unix(),
		RefreshToken: tokens.RefreshToken,
	})

	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	w.Write(response)

}

func Refresh(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	authHeader := r.Header.Get("Authorization")
	if authHeader == "" {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	headerParts := strings.Split(authHeader, " ")
	if len(headerParts) != 2 {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	if headerParts[0] != "Bearer" {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	var request models.RefreshToken

	err := json.NewDecoder(r.Body).Decode(&request)
	if err != nil {
		w.Write([]byte("Cant decode json"))
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	userFromClaims, err := user2.ParseAccessToken(headerParts[1])
	if err != nil {
		//bad access token
		w.WriteHeader(http.StatusUnauthorized)
		return
	}

	ok := user2.IsExpired(userFromClaims.Claims.ExpiresAt)
	if ok {
		//access token expire
		w.Write([]byte("access token expire"))
		w.WriteHeader(http.StatusUnauthorized)
		return
	}

	userFromDb, ok := db.GetUserByGUID(userFromClaims.Claims.Id)
	if !ok {
		w.Write([]byte("user not found"))
		w.WriteHeader(http.StatusUnauthorized)
		return
	}

	ok = helper.CompareHashAndToken(request.RefreshToken, userFromDb.RefreshToken)
	if !ok {
		w.Write([]byte("Bad refresh token"))
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	ok = user2.IsExpired(userFromDb.ExpiresAt.Unix())
	if ok {
		w.Write([]byte("refresh token expired"))
		w.WriteHeader(http.StatusUnauthorized)
		return
	}

	var user models.User
	user.GUID = userFromClaims.Claims.Id

	tokens := user2.GenerateTokens(user)

	ok = db.UpdateUserByGUID(user, tokens)
	if !ok {
		w.Write([]byte("Failed to update user"))
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	expTime, err := time.ParseDuration(config.Conf.AccessTokenTime)
	if err != nil {
		w.Write([]byte("Parse time error"))
		w.WriteHeader(http.StatusInternalServerError)
		return
		
	}

	response, err := json.Marshal(&models.TokenResponse{
		AccessToken:  tokens.AuthToken,
		TokenType:    "Bearer",
		ExpiresAt:    time.Now().Add(expTime).UTC().Unix(),
		RefreshToken: tokens.RefreshToken,
	})
	if err != nil {
		w.Write([]byte("Error when marshal json"))
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	w.Write(response)
	//w.Write([]byte(headerParts[0]))
	//w.Write([]byte(request.RefreshToken))

}
