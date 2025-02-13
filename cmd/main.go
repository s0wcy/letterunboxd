package main

import (
	"github.com/s0wcy/letterunboxd/db"
	"github.com/s0wcy/letterunboxd/scraper"
)

func main() {
	// DB
	db.InitDB()
	defer db.CloseDB()

	// User
	user := "s0wcy"

	// Scrap
	scraper.ScrapUser(user)
	scraper.ScrapProfile(user)
}

// TODO: Compare two users films
// TODO: Compare two users watchlist
// TODO: Know more about your profile: Data visualization on genres, ratings, etc.
// TODO: Sharing is carying: Top rated films of a user never watched by another user.
// TODO: Manual refresh user / films by scraping and compare to db.
