package db

import (
	"context"
	"fmt"

	sdk "github.com/elmasy-com/columbus-sdk"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo/options"
)

// UpdateUniques updates uniqueTlds (every unique TLDs), uniqueDomains (every unique domains) and uniqueSubs (every unique subdomains) collection.
// This function iterates over the domains collection.
func UpdateUniques() error {

	cursor, err := Domains.Find(context.TODO(), bson.M{})
	if err != nil {
		return fmt.Errorf("find() failed: %w", err)
	}
	defer cursor.Close(context.TODO())

	for cursor.Next(context.TODO()) {

		var r sdk.Domain

		err = cursor.Decode(&r)
		if err != nil {
			return fmt.Errorf("failed to decode: %w", err)
		}

		_, err = UniqueTlds.UpdateOne(context.TODO(), bson.M{"tld": r.TLD}, bson.M{"$setOnInsert": bson.M{"tld": r.TLD}}, options.Update().SetUpsert(true))
		if err != nil {
			return fmt.Errorf("failed to update TLD %s: %w", r.TLD, err)
		}

		_, err = UniqueDomains.UpdateOne(context.TODO(), bson.M{"domain": r.Domain}, bson.M{"$setOnInsert": bson.M{"domain": r.Domain}}, options.Update().SetUpsert(true))
		if err != nil {
			return fmt.Errorf("failed to update domain %s: %w", r.Domain, err)
		}

		_, err = UniqueFullDomains.UpdateOne(context.TODO(), bson.M{"domain": r.FullDomain()}, bson.M{"$setOnInsert": bson.M{"domain": r.FullDomain()}}, options.Update().SetUpsert(true))
		if err != nil {
			return fmt.Errorf("failed to update full domain %s: %w", r.Domain, err)
		}

		if r.Sub == "" {
			// Do not insert empty subdomain into uniqueSubs
			continue
		}

		_, err = UniqueSubs.UpdateOne(context.TODO(), bson.M{"sub": r.Sub}, bson.M{"$setOnInsert": bson.M{"sub": r.Sub}}, options.Update().SetUpsert(true))
		if err != nil {
			return fmt.Errorf("failed to update subdomain %s: %w", r.TLD, err)
		}
	}

	if err := cursor.Err(); err != nil {
		return fmt.Errorf("cursor failed: %w", err)
	}

	return nil

}
