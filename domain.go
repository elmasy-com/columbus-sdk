package sdk

import (
	"strings"
)

type Domain struct {
	Domain string   `bson:"domain" json:"domain"`
	Shard  int      `bson:"shard" json:"shard"`
	Subs   []string `bson:"subs" json:"subs"`
}

// GetFull resturns the hostnames as a slice of string.
// If Subs is empty returns nil (theoretically impossible).
func (d *Domain) GetFull() []string {

	var list []string

	for i := range d.Subs {
		if d.Subs[i] == "" {
			list = append(list, d.Domain)
		} else {
			list = append(list, strings.Join([]string{d.Subs[i], d.Domain}, "."))
		}
	}

	return list
}
