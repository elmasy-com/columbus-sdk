package db

import (
	"context"
	"fmt"

	"go.mongodb.org/mongo-driver/bson"
)

// TODO: This is not the real domains nummber
func getDomains() (int64, error) {

	domains, err := Domains.Distinct(context.TODO(), "domain", bson.M{})

	return int64(len(domains)), err
}

func getTotal() (int64, error) {
	return Domains.CountDocuments(context.TODO(), bson.M{})
}

// GetStat resturns the total number of domains (d), the total number of subdomains (s) and the error (if any).
func GetStat() (d int64, s int64, err error) {

	d, err = getDomains()
	if err != nil {
		return 0, 0, fmt.Errorf("failed to get domains(): %w", err)
	}

	s, err = getTotal()

	return d, s, err
}
