package db

import (
	"context"
	"fmt"

	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

var (
	Client *mongo.Client

	Domains *mongo.Collection

	UniqueTlds        *mongo.Collection
	UniqueDomains     *mongo.Collection
	UniqueFullDomains *mongo.Collection
	UniqueSubs        *mongo.Collection
)

// Connect connects to the database using the standard Connection URI.
func Connect(uri string) error {

	var err error

	Client, err = mongo.Connect(context.TODO(), options.Client().ApplyURI(uri))
	if err != nil {
		return fmt.Errorf("connect: %w", err)
	}

	err = Client.Ping(context.Background(), nil)
	if err != nil {
		return fmt.Errorf("ping: %w", err)
	}

	Domains = Client.Database("columbus").Collection("domains")
	UniqueTlds = Client.Database("columbus").Collection("uniqueTlds")
	UniqueDomains = Client.Database("columbus").Collection("uniqueDomains")
	UniqueFullDomains = Client.Database("columbus").Collection("uniqueFullDomains")
	UniqueSubs = Client.Database("columbus").Collection("uniqueSubs")

	return nil
}

// Disconnect gracefully disconnect from the database.
func Disconnect() error {
	return Client.Disconnect(context.Background())
}
