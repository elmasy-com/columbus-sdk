package db

import (
	"context"
	"fmt"

	sdk "github.com/elmasy-com/columbus-sdk"
	"github.com/elmasy-com/columbus-sdk/fault"
	"github.com/elmasy-com/elnet/domain"
	"github.com/elmasy-com/slices"
	"go.mongodb.org/mongo-driver/bson"
)

// Lookup query the DB and returns a list subdomains.
//
// If d has a subdomain, removes it before the query.
//
// If d is invalid return fault.ErrInvalidDomain.
// If failed to get parts of d (eg.: d is a TLD), returns ault.ErrGetPartsFailed.
func Lookup(d string) ([]string, error) {

	if !domain.IsValid(d) {
		return nil, fault.ErrInvalidDomain
	}

	d = domain.Clean(d)

	p := domain.GetParts(d)
	if p == nil || p.Domain == "" || p.TLD == "" {
		return nil, fault.ErrGetPartsFailed
	}

	// Use Find() to find every shard of the domain
	cursor, err := Domains.Find(context.TODO(), bson.M{"domain": p.Domain, "tld": p.TLD})
	if err != nil {
		return nil, fmt.Errorf("failed to find: %s", err)
	}
	defer cursor.Close(context.TODO())

	var subs []string

	for cursor.Next(context.TODO()) {

		var r sdk.Domain

		err = cursor.Decode(&r)
		if err != nil {
			return nil, fmt.Errorf("failed to decode: %s", err)
		}

		subs = append(subs, r.Sub)
	}

	if err := cursor.Err(); err != nil {
		return subs, fmt.Errorf("cursor failed: %w", err)
	}

	return subs, nil
}

// TLD query the DB and returns a list of TLDs for the given domain d.
//
// Domain d must be a valid Second Level Domain (eg.: "example").
//
// NOTE: This function not validate adn Clean() d!
func TLD(d string) ([]string, error) {

	// Use Find() to find every shard of the domain
	cursor, err := Domains.Find(context.TODO(), bson.M{"domain": d})
	if err != nil {
		return nil, fmt.Errorf("failed to find: %s", err)
	}
	defer cursor.Close(context.TODO())

	var tlds []string

	for cursor.Next(context.TODO()) {

		var r sdk.Domain

		err = cursor.Decode(&r)
		if err != nil {
			return nil, fmt.Errorf("failed to decode: %s", err)
		}

		tlds = slices.AppendUnique(tlds, r.TLD)
	}

	if err := cursor.Err(); err != nil {
		return nil, fmt.Errorf("cursor failed: %w", err)
	}

	return tlds, nil
}
