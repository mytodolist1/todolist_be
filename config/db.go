package config

import (
	"context"
	"fmt"
	"os"

	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

// connection to mongoDB
func MongoConnect(mongoString, dbname string) *mongo.Database {
	client, err := mongo.Connect(context.TODO(), options.Client().ApplyURI(os.Getenv(mongoString)))
	if err != nil {
		fmt.Printf("MongoConnect: %v\n", err)
	}

	return client.Database(dbname)
}
