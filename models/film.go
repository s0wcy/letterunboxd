package models

type Film struct {
	ID       string   `bson:"_id"`
	Slug     string   `bson:"slug"`
	Title    string   `bson:"title"`
	Image    string   `bson:"image"`
	Release  string   `bson:"release"`
	Genres   []string `bson:"genres"`
	Ratings  string   `bson:"ratings"`
	Director string   `bson:"director"`
	Cast     []string `bson:"cast"`
}
