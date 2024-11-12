package db

import (
	"database/sql"
	"linkkeeper/types"
	"log"

	_ "github.com/mattn/go-sqlite3"
)

type Database struct {
	db *sql.DB
}

func (db *Database) initDB() error {
	_, err := db.db.Exec("CREATE TABLE IF NOT EXISTS bookmarks(id INTEGER PRIMARY KEY AUTOINCREMENT, url TEXT)")
	if err != nil {
		log.Fatal(err)
		return err
	}
	return nil
}

func NewDB(dbPath string) (*Database, error) {
	db, error := sql.Open("sqlite3", dbPath)
	if error != nil {
		log.Fatal(error)
		return nil, error
	}
	newDb := &Database{db: db}
	newDb.initDB()
	return newDb, nil
}

func (db *Database) CreateBookmark(url string) error {
	_, err := db.db.Exec("INSERT INTO bookmarks(url) VALUES(?)", url)
	if err != nil {
		log.Fatal(err)
		return err
	}
	return nil
}

func (db *Database) ListBookmarks() ([]types.Bookmark, error) {
	rows, err := db.db.Query("SELECT * FROM bookmarks")
	if err != nil {
		log.Fatal(err)
		return nil, err
	}
	defer rows.Close()

	var bookmarks []types.Bookmark
	for rows.Next() {
		var bookmark types.Bookmark
		rows.Scan(&bookmark.Id, &bookmark.Url)
		bookmarks = append(bookmarks, bookmark)
	}
	return bookmarks, nil
}
