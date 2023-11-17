package database

import (
	"context"
	"fmt"
	"log"

	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

var ctx = context.TODO()
var Conn *Client

type Client struct {
	*mongo.Database
}

func Connect(user, password, uri, db string, poolSize int) {
	var err error
	var client *mongo.Client
	fullUri := fmt.Sprintf("mongodb://%s:%s@%s", user, password, uri)
	clientOptions := options.Client().ApplyURI(fullUri).SetMaxPoolSize(uint64(poolSize))
	client, err = mongo.Connect(ctx, clientOptions)
	if err != nil {
		log.Fatal(err)
	}

	err = client.Ping(ctx, nil)
	if err != nil {
		log.Fatal(err)
	}

	Conn = &Client{client.Database(db)}
	fmt.Println("Connected to MongoDB!")
}
