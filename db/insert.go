package db

import (
	"context"
	"fmt"
	"strings"

	"github.com/elmasy-com/columbus-sdk/fault"
	"github.com/elmasy-com/elnet/domain"
	"github.com/elmasy-com/elnet/valid"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo/options"
)

// Insert insert the given domain d to the database.
// Firstly, checks if d is valid. Then split into sub|domain|tld parts.
// Sharding means, if the document is reached the 16MB limit increase the "shard" field by one.
//
// If domain is invalid, returns fault.ErrInvalidDomain.
// If failed to get parts of d (eg.: d is a TLD), returns ault.ErrGetPartsFailed.
func Insert(d string) error {

	if !valid.Domain(d) {
		return fault.ErrInvalidDomain
	}

	d = domain.Clean(d)

	p := domain.GetParts(d)
	if p == nil || p.Domain == "" || p.TLD == "" {
		return fault.ErrGetPartsFailed
	}

	shard := 0

	/*
	 * Always iterate over every shard, because $addToSet iterate over every shard's every subs and append it only if the subdomain not exist.
	 * If sub exist, do nothing.
	 * If sub not exist, add it to the last shard.
	 * This method is slow, but working well to handle duplications.
	 */

	for {

		filter := bson.D{{Key: "domain", Value: p.Domain}, {Key: "tld", Value: p.TLD}, {Key: "shard", Value: shard}}
		update := bson.D{{Key: "$addToSet", Value: bson.M{"subs": p.Sub}}}
		opts := options.Update().SetUpsert(true)

		_, err := Domains.UpdateOne(context.TODO(), filter, update, opts)
		if err == nil {
			return nil
		}

		switch {
		case strings.Contains(err.Error(), "Resulting document after update is larger than 16777216"):
			// Increase shard number by one.
			// So, if document with (domain == example.com && shard == 0) is full, update the (document == example.com && shard == 1).
			shard++
		default:
			return fmt.Errorf("failed to update %s: %s", d, err)
		}
	}
}
