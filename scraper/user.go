package scraper

import (
	"fmt"
	"log"
	"strings"

	"github.com/gocolly/colly"
	"github.com/s0wcy/letterunboxd/db"
	"github.com/s0wcy/letterunboxd/models"
)

func scrapUser(username string) models.User {
	fmt.Printf("\n\n⌛ Scraping user: %s...\n", username)

	url := fmt.Sprintf("https://letterboxd.com/%s/", username)

	// Prepare user data
	var following, followers []string

	// Following
	cFollowing := colly.NewCollector()
	followingUrl := url + "following/"
	cFollowing.OnHTML("tr .person-summary .title-3 a", func(e *colly.HTMLElement) {
		following = append(following, strings.Split(e.Attr("href"), "/")[1])
	})

	// Execute scraper
	err := cFollowing.Visit(followingUrl)
	if err != nil {
		log.Fatal(err)
	}
	cFollowing.Wait()

	// Followers
	cFollowers := colly.NewCollector()
	followersUrl := url + "followers/"
	cFollowers.OnHTML("tr .person-summary .title-3 a", func(e *colly.HTMLElement) {
		followers = append(followers, strings.Split(e.Attr("href"), "/")[1])
	})

	// Execute scraper
	err = cFollowers.Visit(followersUrl)
	if err != nil {
		log.Fatal(err)
	}
	cFollowers.Wait()

	// Combine results
	user := models.User{
		ID:        username,
		Following: following,
		Followers: followers,
	}

	fmt.Printf("\n📦 Scraped: %s\n", url)
	return user
}

func ScrapUser(username string) {
	fmt.Printf("\n\n===============🤠​    Scraping User  🤠​===============\n\n", username)
	fmt.Printf("​😎​ User: %s", username)
	user := scrapUser(username)
	fmt.Printf("==========================================================\n")
	fmt.Printf("🗃️ Storing user: %s...\n", user.ID)
	db.UpdateUserFollows(user)
	fmt.Printf("🗃️ Stored user: %s.\n", user.ID)
	fmt.Printf("\n\n===============🤠​    Scraped User   🤠​===============\n\n", username)
}
