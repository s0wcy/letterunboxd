package db

import (
	"database/sql"
	"fmt"
	"log"
	"strings"

	"github.com/s0wcy/letterunboxd/models"
	_ "modernc.org/sqlite"
)

/**
 * Setup
 */

var db *sql.DB

func InitDB() {
	var err error
	db, err = sql.Open("sqlite", "letterunboxd.db")
	if err != nil {
		log.Fatal("⚠️ InitDB:", err)
	}

	createTables()
	fmt.Println("⚡​ Connected to DB")
}

func createTables() {
	queryCreateFilmsTable := `CREATE TABLE IF NOT EXISTS films (
		id TEXT PRIMARY KEY,
		slug TEXT,
		title TEXT,
		image TEXT,
		release INTEGER,
		genres TEXT,
		ratings TEXT,
		director TEXT,
		cast TEXT
	);`
	_, err := db.Exec(queryCreateFilmsTable)
	if err != nil {
		log.Fatal("⚠️ queryCreateFilmsTable:", err)
	}

	queryCreateUsersTable := `CREATE TABLE IF NOT EXISTS users (
		id TEXT PRIMARY KEY,
		watched TEXT,
		rated TEXT,
		liked TEXT,
		watchlist TEXT,
		following TEXT,
		followers TEXT
	);`
	_, err = db.Exec(queryCreateUsersTable)
	if err != nil {
		log.Fatal("⚠️ queryCreateUsersTable:", err)
	}

	fmt.Println("✅ Setup tables")
}

func CloseDB() {
	db.Close()
	fmt.Println("🔌 Disconnected from DB")
}

/**
 * INSERT
 */
func AddFilm(film models.Film) {
	_, err := db.Exec(`INSERT OR IGNORE INTO films (id, slug, title, image, release, genres, ratings, director, cast)
	VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?)`, film.ID, film.Slug, film.Title, film.Image, film.Release, strings.Join(film.Genres, ","), film.Ratings, film.Director, strings.Join(film.Cast, ","))
	if err != nil {
		log.Println("⚠️ AddFilm:", err)
	}
}

func AddUser(user models.User) {
	_, err := db.Exec(`INSERT OR REPLACE INTO users (id, watched, rated, liked, watchlist, following, followers)
	VALUES (?, ?, ?, ?, ?, ?, ?)`, user.ID, strings.Join(user.Watched, ","), strings.Join(user.Rated, ","), strings.Join(user.Liked, ","), strings.Join(user.Watchlist, ","), strings.Join(user.Following, ","), strings.Join(user.Followers, ","))
	if err != nil {
		log.Println("⚠️ AddUser:", err)
	}
}

func UpdateUserFollows(user models.User) error {
	// Insert the user if not exists
	_, err := db.Exec(`INSERT OR IGNORE INTO users (id) VALUES (?)`, user.ID)
	if err != nil {
		return err
	}

	// Update following and followers
	_, err = db.Exec(`UPDATE users SET following = ?, followers = ? WHERE id = ?`,
		strings.Join(user.Following, ","),
		strings.Join(user.Followers, ","),
		user.ID)
	if err != nil {
		return err
	}

	return nil
}

func UpdateUserProfile(user models.User) error {
	// Insert the user if not exists
	_, err := db.Exec(`INSERT OR IGNORE INTO users (id) VALUES (?)`, user.ID)
	if err != nil {
		return err
	}

	// Update profile information
	_, err = db.Exec(`UPDATE users SET watched = ?, rated = ?, liked = ?, watchlist = ? WHERE id = ?`,
		strings.Join(user.Watched, ","),
		strings.Join(user.Rated, ","),
		strings.Join(user.Liked, ","),
		strings.Join(user.Watchlist, ","),
		user.ID)
	if err != nil {
		return err
	}

	return nil
}

/**
 * GET
 */
func GetFilm(id string) models.Film {
	row := db.QueryRow("SELECT * FROM films WHERE id = ?", id)
	var film models.Film
	var genres, cast string
	err := row.Scan(&film.ID, &film.Slug, &film.Title, &film.Image, &film.Release, &genres, &film.Ratings, &film.Director, &cast)
	if err != nil {
		log.Println("⚠️ GetFilm:", err)
	}
	film.Genres = strings.Split(genres, ",")
	film.Cast = strings.Split(cast, ",")
	return film
}

func GetFilms(ids []string) []models.Film {
	placeholders := strings.Repeat("?,", len(ids)-1) + "?"
	query := "SELECT * FROM films WHERE id IN (" + placeholders + ")"

	args := make([]any, len(ids))
	for i, id := range ids {
		args[i] = id
	}

	rows, err := db.Query(query, args...)
	if err != nil {
		log.Println("⚠️ GetFilms:", err)
		return nil
	}
	defer rows.Close()

	var films []models.Film
	for rows.Next() {
		var film models.Film
		var genres, cast string
		err := rows.Scan(&film.ID, &film.Slug, &film.Title, &film.Image, &film.Release, &genres, &film.Ratings, &film.Director, &cast)
		if err != nil {
			log.Println("⚠️ GetFilms:", err)
			continue
		}
		film.Genres = strings.Split(genres, ",")
		film.Cast = strings.Split(cast, ",")
		films = append(films, film)
	}
	return films
}

func GetUser(id string) models.User {
	row := db.QueryRow("SELECT * FROM users WHERE id = ?", id)
	var user models.User
	var watched, rated, liked, watchlist, following, followers string
	err := row.Scan(&user.ID, &watched, &rated, &liked, &watchlist, &following, &followers)
	if err != nil {
		log.Println("⚠️ GetUser:", err)
	}
	user.Watched = strings.Split(watched, ",")
	user.Rated = strings.Split(rated, ",")
	user.Liked = strings.Split(liked, ",")
	user.Watchlist = strings.Split(watchlist, ",")
	user.Following = strings.Split(following, ",")
	user.Followers = strings.Split(followers, ",")
	return user
}
