package db

import (
	"context"
	"fmt"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"jwttask/config"
	"log"
)

func ConnectDB() *mongo.Collection {

	clientOptions := options.Client().ApplyURI(config.Conf.DatabaseURI)
	client, err := mongo.Connect(context.TODO(), clientOptions)

	if err != nil {
		log.Fatal(err)
	}

	fmt.Println("Successfully connected to MongoDB!")

	collection := client.Database(config.Conf.DatabaseName).Collection("users")

	fmt.Println("Successfully create collection")
	return collection
}
