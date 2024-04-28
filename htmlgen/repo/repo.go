package repo

import (
	"database/sql"
	"encoding/json"
	"errors"

	_ "github.com/mattn/go-sqlite3"
	"github.com/rs/zerolog/log"
)

var db *sql.DB

func init() {
	// Open
	var err error
	db, err = sql.Open("sqlite3", "./cache.sqlite")
	if err != nil {
		log.Fatal().Err(err).
			Str("location", "./cache.sqlite").
			Msg("could not open cache db")
	}

	// Test
	_, err = db.Exec("SELECT 1")
	if err != nil {
		log.Fatal().Err(err).
			Str("location", "./cache.sqlite").
			Msg("db test failed")
	}

	// Migrate
	sql := `
CREATE TABLE IF NOT EXISTS letterboxd(
	review_uri TEXT NOT NULL PRIMARY KEY,
	data JSONB NOT NULL
)`

	_, err = db.Exec(sql)
	if err != nil {
		log.Fatal().Err(err).
			Str("location", "./cache.sqlite").
			Msg("migration failed")
	}

	log.Debug().Msg("opened cache db")
}

type LetterboxdInfo struct {
	TMDBid int `json:"tmdb_id"`
}

func GetLetterboxdInfo(letterboxdURI string) (LetterboxdInfo, bool) {
	var j string
	query := `
SELECT data FROM letterboxd WHERE review_uri = $1`
	err := db.QueryRow(query, letterboxdURI).Scan(&j)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return LetterboxdInfo{}, false
		}
		log.Fatal().Err(err).
			Str("query", query).
			Msg("unexpected sqlite fail")
	}

	var info LetterboxdInfo
	err = json.Unmarshal([]byte(j), &info)
	if err != nil {
		log.Fatal().Err(err).
			Str("info", j).
			Msg("could not unmarshal letterboxd info")
	}

	return info, true
}

func InsertLetterboxdInfo(letterboxdURI string, info LetterboxdInfo) {
	j, err := json.Marshal(info)
	if err != nil {
		log.Fatal().Err(err).
			Interface("info", info).
			Msg("couldn't marshal letterboxd info")
	}

	query := `
INSERT INTO letterboxd VALUES ($1,$2)`
	_, err = db.Exec(query, letterboxdURI, string(j))
	if err != nil {
		log.Fatal().Err(err).
			Str("query", query).
			Msg("unexpected sqlite fail")
	}
}
