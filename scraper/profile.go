package scraper

import (
	"fmt"
	"log"
	"regexp"
	"strconv"
	"strings"

	"github.com/gocolly/colly"
	"github.com/s0wcy/letterunboxd/db"
	"github.com/s0wcy/letterunboxd/models"
)

// Internal
func scrapProfilePage(url string) []models.Film {
	fmt.Printf("​⌛ Scraping Profile Page... %s\n", url)
	isWatchlist := strings.Contains(url, "watchlist")

	// Prepare films
	var films []models.Film
	var filmsIds []string
	var filmsSlugs []string
	var filmsRatings []string
	var filmsLikes []string

	// Scraper search
	cProfile := colly.NewCollector()

	// Ids & Slugs
	cProfile.OnHTML(".poster-container .poster", func(e *colly.HTMLElement) {
		id := e.Attr("data-film-id")
		filmsIds = append(filmsIds, id)

		slug := e.Attr("data-film-slug")
		filmsSlugs = append(filmsSlugs, slug)
	})

	// Ratings
	if !isWatchlist {
		cProfile.OnHTML(".poster-container .poster-viewingdata", func(e *colly.HTMLElement) {
			rate := e.DOM.Find(".rating")
			rating := "null"

			if rate.Length() > 0 {
				class, _ := rate.Attr("class")
				classes := strings.Split(class, " ")
				fullRating := regexp.MustCompile(`rated-(\d+)`)

				for _, class := range classes {
					match := fullRating.FindStringSubmatch(class)
					if len(match) > 1 {
						rating = match[1]
					}
				}
			}

			filmsRatings = append(filmsRatings, rating)

			likeStatus := "not liked"
			if e.DOM.Parent().Find(".icon-liked").Length() > 0 {
				likeStatus = "liked"
			}

			filmsLikes = append(filmsLikes, likeStatus)
		})
	}

	// Execute scraper
	err := cProfile.Visit(url)
	if err != nil {
		log.Fatal(err)
	}
	cProfile.Wait()

	// Combine results
	for i, _ := range filmsIds {
		rated := "null"
		if !isWatchlist {
			rated = filmsRatings[i]
		}
		liked := "null"
		if !isWatchlist {
			liked = filmsLikes[i]
		}

		film := models.Film{
			ID:       filmsIds[i],
			Slug:     filmsSlugs[i],
			Ratings:  rated,
			Director: liked, // Poopy code, I know... I just use director as a pass variable for like status
		}
		films = append(films, film)
	}

	fmt.Printf("​📦​ Scraped %s\n", url)
	return films
}

func scrapProfileCategoryPages(username string, category string) []models.Film {
	fmt.Printf("\n⌛ Scraping Profile Category Pages...\n")
	// Scrap films
	var films []models.Film
	url := fmt.Sprintf("https://letterboxd.com/%s/%s/", username, category)
	fmt.Printf("\n🔍​ Category: '%s'\n", category)

	// Scraper
	cCategory := colly.NewCollector()

	// Total pages
	totalPages := 1
	cCategory.OnHTML(".paginate-page:last-child a", func(e *colly.HTMLElement) {
		pageNum, _ := strconv.Atoi(strings.TrimSpace(e.Text))
		totalPages = pageNum
	})

	// Execute scraper
	err := cCategory.Visit(url)
	if err != nil {
		log.Fatal(err)
	}
	cCategory.Wait()

	// Scrap all film pages
	fmt.Printf("\n📄​ Pages found: %d\n\n", totalPages)
	for i := 1; i <= totalPages; i++ {
		pageURL := fmt.Sprintf("https://letterboxd.com/%s/%s/page/%d/", username, category, i)
		newFilms := scrapProfilePage(pageURL)
		films = append(films, newFilms...)
	}

	fmt.Printf("\n😎​ '%s' %s: \n\n", username, category)
	for _, film := range films {
		fmt.Printf("%s - %s\n", film.ID, film.Slug)
	}

	return films
}

// TODO: Refresh dynamically on new requests
func scrapProfileFilmsCount(username string) (int, int) {
	// Films count
	cFilms := colly.NewCollector()
	filmsCount := 0
	cFilms.OnHTML(".profile-stats h4:first-child a .value", func(e *colly.HTMLElement) {
		filmsCount, _ = strconv.Atoi(strings.TrimSpace(e.Text))
	})

	// Films URL
	urlFilms := fmt.Sprintf("https://letterboxd.com/%s/", username)

	// Execute scraper
	err := cFilms.Visit(urlFilms)
	if err != nil {
		log.Fatal(err)
	}
	cFilms.Wait()

	// Watchlist count
	cWatchlist := colly.NewCollector()
	watchlistCount := 0
	cWatchlist.OnHTML(".js-watchlist-content", func(e *colly.HTMLElement) {
		watchlistCount, _ = strconv.Atoi(e.Attr("data-num-entries"))
	})

	// Films URL
	urlWatchlist := fmt.Sprintf("https://letterboxd.com/%s/watchlist", username)

	// Execute scraper
	err = cWatchlist.Visit(urlWatchlist)
	if err != nil {
		log.Fatal(err)
	}
	cWatchlist.Wait()

	return filmsCount, watchlistCount
}

// Public
func ScrapProfile(username string) {
	fmt.Printf("\n\n===============➡️   Scrap '%s' Profile  ​⬅️===============\n\n", username)

	// Scrap Profile
	films := scrapProfileCategoryPages(username, "films")
	watchlist := scrapProfileCategoryPages(username, "watchlist")

	// Update User & Films DB
	user := models.User{
		ID:        username,
		Watched:   []string{},
		Rated:     []string{},
		Liked:     []string{},
		Watchlist: []string{},
	}
	for _, film := range films {
		user.Watched = append(user.Watched, film.ID)
		user.Rated = append(user.Rated, film.Ratings)
		user.Liked = append(user.Liked, film.Director)
		ScrapFilm(film.Slug)
	}
	for _, watchl := range watchlist {
		user.Watchlist = append(user.Watched, watchl.ID)
		ScrapFilm(watchl.Slug)
	}
	fmt.Printf("ScrapProfile ===============> %v\n", user)
	db.UpdateUserProfile(user)
	fmt.Printf("\n\n===============➡️  '%s' Profile Scraped ​⬅️===============\n\n", username)
}
