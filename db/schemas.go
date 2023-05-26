package db

// Schema used in *notFound* collection
type NotFoundSchema struct {
	Domain string `bson:"domain" json:"domain"`
}
