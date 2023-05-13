package db

import (
	"context"

	"github.com/elmasy-com/columbus-sdk/fault"
	"github.com/elmasy-com/elnet/domain"
	"github.com/elmasy-com/elnet/valid"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo/options"
)

// Insert insert the given domain d to the database.
// Firstly, checks if d is valid. Then split into sub|domain|tld parts.
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

	doc := bson.D{{Key: "domain", Value: p.Domain}, {Key: "tld", Value: p.TLD}, {Key: "sub", Value: p.Sub}}

	_, err := Domains.UpdateOne(context.TODO(), doc, bson.M{"$setOnInsert": doc}, options.Update().SetUpsert(true))
	if err != nil {
		return err
	}

	return nil
}
