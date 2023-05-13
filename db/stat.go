package db

import (
	"context"
	"fmt"

	"go.mongodb.org/mongo-driver/bson"
)

// Count the documents in the uniqueTlds collection
func getUniqueTldsNum() (int64, error) {

	n, err := UniqueTlds.CountDocuments(context.TODO(), bson.M{})

	return int64(n), err
}

func getUniqueDomainsNum() (int64, error) {

	n, err := UniqueDomains.CountDocuments(context.TODO(), bson.M{})

	return int64(n), err
}

func getUniqueFullDomainsNum() (int64, error) {

	n, err := UniqueFullDomains.CountDocuments(context.TODO(), bson.M{})

	return int64(n), err
}

func getUniqueSubsNum() (int64, error) {

	n, err := UniqueSubs.CountDocuments(context.TODO(), bson.M{})

	return int64(n), err
}

func getTotal() (int64, error) {
	return Domains.CountDocuments(context.TODO(), bson.M{})
}

// GetStat returns the total number of domains, the total number of unique TLDs, the total number of unique domains,
// the total number of unique full domains and the total number of subdomains and the error (if any).
func GetStat() (total, tlds, domains, fullDomain, subs int64, err error) {

	total, err = getTotal()
	if err != nil {
		return 0, 0, 0, 0, 0, fmt.Errorf("failed to get total: %w", err)
	}

	tlds, err = getUniqueTldsNum()
	if err != nil {
		return 0, 0, 0, 0, 0, fmt.Errorf("failed to get unique TLDs: %w", err)
	}

	domains, err = getUniqueDomainsNum()
	if err != nil {
		return 0, 0, 0, 0, 0, fmt.Errorf("failed to get unique domains: %w", err)
	}

	fullDomain, err = getUniqueFullDomainsNum()
	if err != nil {
		return 0, 0, 0, 0, 0, fmt.Errorf("failed to get unique full domains: %w", err)
	}

	subs, err = getUniqueSubsNum()
	if err != nil {
		return 0, 0, 0, 0, 0, fmt.Errorf("failed to get unique subdomains: %w", err)
	}

	return total, tlds, domains, fullDomain, subs, err
}
