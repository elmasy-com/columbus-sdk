package sdk

var (
	ApiKey = ""
)

type User struct {
	Key   string `bson:"key" json:"key"`
	Name  string `bson:"name" json:"name"`
	Admin bool   `bson:"admin" json:"admin"`
}
