package tmdb

import (
	"fmt"
	"os"
	"slices"

	gotmdb "github.com/cyruzin/golang-tmdb"
	"github.com/rs/zerolog/log"
)

var client *gotmdb.Client

func init() {
	// Init client
	var err error
	client, err = gotmdb.Init(os.Getenv("TMDB_API_KEY"))
	if err != nil {
		log.Fatal().Err(err).Msg("could not create tmdb client")
	}
	client.SetClientAutoRetry()
	client.SetAlternateBaseURL()
}

func GetPoster(movieID int) {
	// Fetch
	res, err := client.GetMovieImages(movieID, nil)
	if err != nil {
		log.Fatal().Err(err).
			Int("movie_id", movieID).
			Msg("could not fetch images")
	}

	// Need at least one poster
	if len(res.Posters) == 0 {
		log.Fatal().Err(err).
			Int("movie_id", movieID).
			Msg("no posters available")
	}

	// Sort posters by vote count, then size
	slices.SortFunc(res.Posters, func(a, b gotmdb.MovieImage) int {
		if a.VoteCount != b.VoteCount {
			// Reverse order, biggest first
			return int(b.VoteCount - a.VoteCount)
		}

		// Reverse order, biggest first
		return b.Height - a.Height
	})
	poster := res.Posters[0]

	fmt.Println(poster)
}
