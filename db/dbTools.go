package db

import (
	"context"
	"go.mongodb.org/mongo-driver/bson"
	"jwttask/config"
	"jwttask/helper"
	"jwttask/models"
	"log"
	"time"
)

func GetUserByGUID(guid string) (models.User, bool) {
	var user models.User
	collection := ConnectDB()

	filter := bson.M{"guid": guid}
	err := collection.FindOne(context.Background(), filter).Decode(&user)
	if err != nil {
		log.Println("User not find")
		return user, false
	}
	return user, true
}

func DeleteUserByGUID(guid string) bool {
	collection := ConnectDB()

	filter := bson.M{"guid": guid}
	_, err := collection.DeleteOne(context.Background(), filter)
	if err != nil {
		log.Println("Cant delete user")
		return false
	}
	return true
}

func InsertUserByGUID(user models.User, tokens models.AuthToken) bool {
	collection := ConnectDB()

	hashRefreshToken, err := helper.HashToken(tokens.RefreshToken)
	if err != nil {
		log.Println("Cant encrypt token")
		return false
	}
	expTime, err := time.ParseDuration(config.Conf.RefreshTokenTime)
	if err != nil {
		log.Println(err)
	}

	_, err = collection.InsertOne(context.Background(), bson.D{
		{"guid", user.GUID},
		{"refreshToken", hashRefreshToken},
		{"expiresAT", time.Now().Add(expTime)},
	})

	if err != nil {
		log.Println(err)
		return false
	}
	return true

}

func UpdateUserByGUID(user models.User, tokens models.AuthToken) bool {
	collection := ConnectDB()

	hashRefreshToken, err := helper.HashToken(tokens.RefreshToken)
	if err != nil {
		log.Println(err)
		return false
	}
	expTime, err := time.ParseDuration(config.Conf.RefreshTokenTime)
	if err != nil {
		log.Println(err)
	}

	filter := bson.M{"guid": user.GUID}

	update := bson.D{{"$set", bson.D{
		{"refreshToken", hashRefreshToken},
		{"expiresAT", time.Now().Add(expTime)},
	}}}

	_, err = collection.UpdateOne(context.Background(), filter, update)

	if err != nil {
		log.Println("failed to update user")
		return false
	}
	return true

}
