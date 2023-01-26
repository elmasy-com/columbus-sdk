package db

import (
	"context"
	"fmt"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

// Creates the indexes in domains collection.
func SetIndex() error {

	// Create a unique compound index for domain+shard in domains.
	// MongoDB will ignore this block if the index already exist.
	_, err := Domains.Indexes().CreateOne(context.TODO(), mongo.IndexModel{
		Keys:    bson.D{{Key: "domain", Value: 1}, {Key: "shard", Value: 1}},
		Options: options.Index().SetUnique(true),
	})
	if err != nil {
		return fmt.Errorf("failed to create unique domain+shard compound index: %s", err)
	}

	return nil
}
