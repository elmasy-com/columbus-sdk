package db

import (
	"context"
	"fmt"
	"strings"

	sdk "github.com/elmasy-com/columbus-sdk"
	"github.com/elmasy-com/columbus-sdk/fault"
	eldomain "github.com/elmasy-com/elnet/domain"
	"github.com/elmasy-com/slices"
	"go.mongodb.org/mongo-driver/bson"
)

// Lookup query the DB and returns a list subdomains.
//
// If d has a subdomain, removes it before the query.
//
// If d is invalid return fault.ErrInvalidDomain.
func Lookup(d string) ([]string, error) {

	if !eldomain.IsValid(d) {
		return nil, fault.ErrInvalidDomain
	}

	d = eldomain.Clean(d)

	p := eldomain.GetParts(d)
	if p == nil || p.Domain == "" {
		return nil, fault.ErrInvalidDomain
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

		subs = append(subs, r.Subs...)
	}

	if err := cursor.Err(); err != nil {
		return subs, fmt.Errorf("cursor failed: %w", err)
	}

	return subs, nil
}

// TLD query the DB and returns a list of TLDs for the given domain d.
//
// If d has a subdomain, removes it before the query.
//
// If d is invalid return fault.ErrInvalidDomain.
func TLD(d string) ([]string, error) {

	if !eldomain.IsValid(d) || strings.ContainsRune(d, '.') {
		return nil, fault.ErrInvalidDomain
	}

	d = eldomain.Clean(d)

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
