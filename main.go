package main

import (
	"context"
	"fmt"

	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

var (
	StatsDB *mongo.Client
	TestDB  *mongo.Client
)

func main() {
	ctx := context.Background()
	opts := options.Client().ApplyURI(
		"mongodb+srv://admin:<password>@skunkworks20210412.zgc06.mongodb.net/myFirstDatabase?retryWrites=true&w=majority",
	).SetAuth(options.Credential{Username: "admin", Password: "ddfda"}).SetServerAPIOptions(&options.ServerAPIOptions{ServerAPIVersion: options.ServerAPIVersion1})
	client, err := mongo.Connect(ctx, opts)
	fmt.Println(err)
	StatsDB = client

	opts = options.Client().ApplyURI(
		"mongodb+srv://admin:<password>@testdb1.zgc06.mongodb.net/myFirstDatabase?retryWrites=true&w=majority",
	).SetAuth(options.Credential{Username: "admin", Password: "ddfda"}).SetServerAPIOptions(&options.ServerAPIOptions{ServerAPIVersion: options.ServerAPIVersion1})
	client, err = mongo.Connect(ctx, opts)
	fmt.Println(err)
	TestDB = client
}
