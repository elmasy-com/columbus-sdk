package sdk

import "strings"

type Domain struct {
	Domain string `bson:"domain" json:"domain"`
	TLD    string `bson:"tld" json:"tld"`
	Sub    string `bson:"sub" json:"sub"`
}

func (d *Domain) String() string {

	if d.Sub == "" {
		return strings.Join([]string{d.Domain, d.TLD}, ".")
	} else {
		return strings.Join([]string{d.Sub, d.Domain, d.TLD}, ".")
	}
}
