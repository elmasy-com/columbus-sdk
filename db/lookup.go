package db

import (
	"context"
	"fmt"
	"strings"

	"github.com/elmasy-com/columbus-sdk/domain"
	"github.com/elmasy-com/columbus-sdk/fault"
	eldomain "github.com/elmasy-com/elnet/domain"
	"go.mongodb.org/mongo-driver/bson"
)

// Lookup query the DB and returns a list subdomains.
// If d is invalid return fault.ErrInvalidDomain.
func Lookup(d string) ([]string, error) {

	// Use Find() to find every shard of the domain

	if !eldomain.IsValid(d) {
		return nil, fault.ErrInvalidDomain
	}

	d = strings.ToLower(d)
	d = eldomain.GetDomain(d)
	if d == "" {
		return nil, fault.ErrInvalidDomain
	}

	cursor, err := Domains.Find(context.TODO(), bson.M{"domain": d})
	if err != nil {
		return nil, fmt.Errorf("failed to find: %s", err)
	}
	defer cursor.Close(context.TODO())

	var subs []string

	for cursor.Next(context.TODO()) {

		var r domain.Domain

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
