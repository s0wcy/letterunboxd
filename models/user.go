package models

type User struct {
	ID        string   `bson:"_id"`
	Watched   []string `bson:"watched"`
	Rated     []string `bson:"rated"`
	Liked     []string `bson:"liked"`
	Watchlist []string `bson:"watchlist"`
	Following []string `bson:"following"`
	Followers []string `bson:"followers"`
}
