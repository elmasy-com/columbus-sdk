package domain

import "strings"

type Domain struct {
	Domain string   `bson:"domain" json:"domain"`
	TLD    string   `bson:"tld" json:"tld"`
	Shard  int      `bson:"shard" json:"shard"`
	Subs   []string `bson:"subs" json:"subs"`
}

// GetFull resturns the hostnames as a slice of string.
// If Subs, Domain or TLD is empty returns nil.
func (d *Domain) GetFull() []string {

	var list []string = nil

	if d.Domain == "" || d.TLD == "" {
		return nil
	}

	for i := range d.Subs {
		if d.Subs[i] == "" {
			list = append(list, strings.Join([]string{d.Domain, d.TLD}, "."))
		} else {
			list = append(list, strings.Join([]string{d.Subs[i], d.Domain, d.TLD}, "."))
		}
	}

	return list
}
