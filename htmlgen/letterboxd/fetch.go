package letterboxd

import (
	"errors"
	"fmt"
	"io"
	"net/http"
	"os"
	"path/filepath"
	"regexp"
	"strconv"

	"github.com/liampulles/liampulles.github.io/htmlgen/repo"
	"github.com/rs/zerolog/log"
)

type letterboxdInfo struct {
	TMDBid     int
	PosterHref string
}

// Fetch data related to a film. Will try and used cached info in the db first.
func FetchData(letterboxdURI string) letterboxdInfo {
	// Try get from cache
	info, ok := fetchFromCache(letterboxdURI)
	if ok {
		return info
	}

	// Ok, we'll have to get it manually. Resolve the TMDB id and poster url first.
	log.Debug().
		Str("letterboxd_uri", letterboxdURI).
		Msg("need to resolve letterboxd info 'manually'")
	tmdbID, posterURL := resolveTMDBidAndPosterURL(letterboxdURI)

	// Check and download the poster
	posterHref := findOrDownloadImage(tmdbID, posterURL)

	// Set it in the db for next time
	repoInfo := repo.LetterboxdInfo{
		TMDBid: tmdbID,
	}
	repo.InsertLetterboxdInfo(letterboxdURI, repoInfo)

	// Map our version
	return letterboxdInfo{
		TMDBid:     tmdbID,
		PosterHref: posterHref,
	}
}

func fetchFromCache(letterboxdURI string) (letterboxdInfo, bool) {
	// Query
	repoInfo, ok := repo.GetLetterboxdInfo(letterboxdURI)
	if !ok {
		return letterboxdInfo{}, false
	}

	// Map
	return letterboxdInfo{
		TMDBid:     repoInfo.TMDBid,
		PosterHref: posterHref(repoInfo.TMDBid),
	}, true
}

var filmLinkRegex = regexp.MustCompile(`<[^>]*"film-title-wrapper"[^>]*>[\s\S]*<a href="(\/film\/[^/]+)/">`)
var tmdbIDRegex = regexp.MustCompile(`"https://www\.themoviedb\.org/\w*/(\d*)/"`)
var posterRegex = regexp.MustCompile(`{"image":"([^"]*)",`)

func resolveTMDBidAndPosterURL(letterboxdURI string) (int, string) {
	// Get the review page
	reviewBody := fetchPage(letterboxdURI)

	// Parse the film link
	elem := filmLinkRegex.FindSubmatch(reviewBody)
	if len(elem) < 2 {
		err := errors.New("couldn't extract film link from review")
		log.Fatal().Err(err).
			Str("letterboxd_uri", letterboxdURI).
			Msg("could not resolve TMDB id")
	}
	filmURL := "https://letterboxd.com" + string(elem[1])

	// Parse the poster link
	elem = posterRegex.FindSubmatch(reviewBody)
	if len(elem) < 2 {
		err := errors.New("couldn't extract poster url from review")
		log.Fatal().Err(err).
			Str("letterboxd_uri", letterboxdURI).
			Msg("could not resolve TMDB id")
	}
	posterURL := elem[1]

	// Now fetch the film page
	filmBody := fetchPage(filmURL)

	// Parse the TMDB id
	elem = tmdbIDRegex.FindSubmatch(filmBody)
	if len(elem) < 2 {
		err := errors.New("couldn't extract TMDB id from film page")
		log.Fatal().Err(err).
			Str("url", filmURL).
			Msg("could not resolve TMDB id")
	}
	tmdbID, err := strconv.Atoi(string(elem[1]))
	if err != nil {
		log.Fatal().Err(err).
			Str("url", filmURL).
			Str("tmdb_id", string(elem[1])).
			Msg("did not find correct TMDB id - not an int")
	}

	return tmdbID, string(posterURL)
}

func fetchPage(url string) []byte {
	// Make request
	res, err := http.Get(url)
	if err == nil && (res.StatusCode < 200 || res.StatusCode > 399) {
		err = fmt.Errorf("error reading url: %d", res.StatusCode)
	}
	if err != nil {
		log.Fatal().Err(err).
			Str("url", url).
			Msg("http client error")
	}

	// Read response
	b, err := io.ReadAll(res.Body)
	if err != nil {
		log.Fatal().Err(err).
			Str("url", url).
			Msg("could not read response body")
	}

	return b
}

func findOrDownloadImage(tmdbID int, posterURL string) string {
	filename := fmt.Sprintf("%d.jpg", tmdbID)
	p := filepath.Join("static", "images", "review-posters", filename)
	href := posterHref(tmdbID)

	// Is it downloaded already? Great if so.
	_, err := os.Stat(p)
	if err == nil {
		return href
	}

	// Ok, then we need to download it.
	file, err := os.Create(p)
	if err != nil {
		log.Fatal().Err(err).
			Str("path", p).
			Msg("could not create file to write image to")
	}
	defer file.Close()

	resp, err := http.Get(posterURL)
	if err != nil {
		log.Fatal().Err(err).
			Str("poster_url", posterURL).
			Msg("could not download poster")
	}
	defer resp.Body.Close()

	_, err = io.Copy(file, resp.Body)
	if err != nil {
		log.Fatal().Err(err).
			Str("poster_url", posterURL).
			Msg("could not write file")
	}

	// Ok, now we're done
	return href
}
func posterHref(tmdbID int) string {
	return fmt.Sprintf("/images/review-posters/%d.jpg", tmdbID)
}
