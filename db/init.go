package db

import (
	"context"
	"fmt"
	"log"

	"github.com/sirupsen/logrus"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type Resource struct {
	DB *mongo.Database
}

// Close use this method to close database connection
func (r *Resource) Close() {
	logrus.Warning("Closing all db connections")
}

func InitResource() (*Resource, error) {
	// Set client options

	// for local:
	// clientOptions := options.Client().ApplyURI("mongodb://root:root@localhost:27017")

	// for docker compose:
	clientOptions := options.Client().ApplyURI("mongodb://root:root@mongodb-go-gin:27017")

	// Connect to MongoDB
	client, err := mongo.Connect(context.TODO(), clientOptions)
	fmt.Println("bbbb")
	if err != nil {
		log.Fatal(err)
	}

	// Check the connection
	err = client.Ping(context.TODO(), nil)

	if err != nil {
		log.Fatal(err)
	}

	fmt.Println("Connected to MongoDB!")

	return &Resource{DB: client.Database("demo")}, nil
}
