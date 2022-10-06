package sdk

var (
	ApiKey = ""
)

type User struct {
	Key   string `bson:"key"`
	Name  string `bson:"name"`
	Admin bool   `bson:"admin"`
}
