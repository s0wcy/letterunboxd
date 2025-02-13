package scraper

import (
	"fmt"
	"log"
	"strings"

	"github.com/gocolly/colly"
	"github.com/s0wcy/letterunboxd/db"
	"github.com/s0wcy/letterunboxd/models"
)

func scrapFilm(slug string) models.Film {
	fmt.Printf("​​⌛ Scraping movie: %s...\n", slug)

	url := fmt.Sprintf("https://letterboxd.com/film/%s/", slug)

	// Prepare film data
	var filmId, filmTitle, filmImage, filmRelease, filmRatings, filmDirector string
	var filmGenres, filmCast []string

	// Scraper search
	cFilms := colly.NewCollector()

	// ID
	cFilms.OnHTML("#backdrop", func(e *colly.HTMLElement) {
		filmId = e.Attr("data-film-id")
	})

	// Title
	cFilms.OnHTML("meta[property='og:title']", func(e *colly.HTMLElement) {
		filmTitle = e.Attr("content")
	})

	// Release
	cFilms.OnHTML(".releaseyear a", func(e *colly.HTMLElement) {
		filmRelease = e.Text
	})

	// Genres
	cFilms.OnHTML("#tab-genres .text-sluglist:nth-of-type(1) p a", func(e *colly.HTMLElement) {
		filmGenres = append(filmGenres, e.Text)
	})

	// Ratings
	cFilms.OnHTML("meta[name='twitter:data2']", func(e *colly.HTMLElement) {
		filmRatings = strings.Fields(e.Attr("content"))[0]
	})

	// Director
	cFilms.OnHTML(".directorlist a span", func(e *colly.HTMLElement) {
		filmDirector = e.Text
	})

	// Director
	cFilms.OnHTML(".directorlist a span", func(e *colly.HTMLElement) {
		filmDirector = e.Text
	})

	// Cast
	cFilms.OnHTML("#tab-cast .cast-list p a", func(e *colly.HTMLElement) {
		filmCast = append(filmCast, e.Text)
	})

	// Execute scraper
	err := cFilms.Visit(url)
	if err != nil {
		log.Fatal(err)
	}
	cFilms.Wait()

	// Image
	segmentedId := strings.Join(strings.Split(filmId, ""), "/")
	image := fmt.Sprintf("https://a.ltrbxd.com/resized/film-poster/%s/%s-%s-0-230-0-345-crop.jpg", segmentedId, filmId, slug)
	// TODO: List others URLs formats such as Dune film:
	// https://a.ltrbxd.com/resized/sm/upload/nx/8b/vs/gc/cDbNAY0KM84cxXhmj8f0dLWza3t-0-70-0-105-crop.jpg?v=49eed12751
	filmImage = image

	// Combine results
	film := models.Film{
		ID:       filmId,
		Slug:     slug,
		Title:    filmTitle,
		Image:    filmImage,
		Release:  filmRelease,
		Genres:   filmGenres,
		Ratings:  filmRatings,
		Director: filmDirector,
		Cast:     filmCast,
	}

	fmt.Printf("📦​ Scraped: %s\n", url)
	return film
}

func ScrapFilm(slug string) {
	fmt.Printf("\n\n===============🍿    Scraping Film   ​🍿===============\n\n", slug)
	fmt.Printf("​🎬​ Movie: %s", slug)
	film := scrapFilm(slug)
	fmt.Printf("==========================================================\n")
	fmt.Printf("🗃️ Storing film: %s...\n", film.Slug)
	db.AddFilm(film)
	fmt.Printf("🗃️ Stored film: %s.\n", film.Slug)
	fmt.Printf("\n\n===============🍿    Scraped Film    ​🍿===============\n\n", slug)
}
