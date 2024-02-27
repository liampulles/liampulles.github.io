package letterboxd

import (
	"archive/zip"
	"encoding/csv"
	"errors"
	"fmt"
	"io"
	"math"
	"os"
	"path/filepath"
	"regexp"
	"slices"
	"strconv"
	"strings"

	"cloud.google.com/go/civil"
	"github.com/rs/zerolog/log"
)

type Review struct {
	Date          civil.Date
	Name          string
	Year          int
	LetterboxdURI string
	Rating        int // Out of 10, divide by 2 to get star rating. 0 means no rating.
	Rewatch       bool
	Review        string // Multiline. Potentially partial HTML.
	PosterHref    string
}

type UserData struct {
	Reviews []Review
}

const exportFolder = "_letterboxd_exports"

func ReadExport() (UserData, error) {
	// Choose the "best" zip in the dir to read
	zipPath, err := pickZip()
	if err != nil {
		return UserData{}, err
	}

	// Prepare a CSV reader
	csvReader, deferFn, err := openReviewsCSV(zipPath)
	if err != nil {
		return UserData{}, err
	}
	defer deferFn()

	// Get the reviews
	reviews, err := readReviewsCSV(csvReader)
	if err != nil {
		return UserData{}, err
	}

	log.Debug().Str("zip", zipPath).Msg("read letterboxd export")
	return UserData{
		Reviews: reviews,
	}, nil
}

var exportZipRegex = regexp.MustCompile(`^letterboxd-.*\.zip$`)

func pickZip() (string, error) {
	// What files are there
	entries, err := os.ReadDir(exportFolder)
	if err != nil {
		log.Err(err).
			Str("dir", exportFolder).
			Msg("could not read letterboxd folder")
		return "", err
	}

	// Want export zips only
	var zipPaths []string
	for _, entry := range entries {
		if entry.IsDir() {
			continue
		}
		if ok := exportZipRegex.MatchString(entry.Name()); !ok {
			continue
		}

		zipPaths = append(zipPaths, entry.Name())
	}

	// Must be at least one
	if len(zipPaths) == 0 {
		err = errors.New("no letterboxd zip files")
		log.Err(err).
			Str("dir", exportFolder).
			Msg("no letterboxd zip files to read")
		return "", err
	}

	// Get the latest one
	return filepath.Join(exportFolder, slices.Max(zipPaths)), nil
}

func openReviewsCSV(zipPath string) (r *csv.Reader, deferFn func(), err error) {
	// Define defer func now, in case we have to use it here.
	var closeFuncs []func() error
	deferFn = func() {
		for i := len(closeFuncs) - 1; i >= 0; i-- {
			closeFuncs[i]()
		}
	}
	defer func() {
		if err != nil {
			deferFn()
		}
	}()

	// Open the archive
	archive, err := zip.OpenReader(zipPath)
	if err != nil {
		log.Err(err).
			Str("zip", zipPath).
			Msg("could not read letterboxd zip")
		return nil, nil, err
	}
	closeFuncs = append(closeFuncs, archive.Close)

	// Get the reviews.csv file
	var file *zip.File
	for _, f := range archive.File {
		if f.Name == "reviews.csv" {
			file = f
			break
		}
	}
	if file == nil {
		err := errors.New("cannot find reviews.csv in archive")
		log.Err(err).
			Str("zip", zipPath).
			Msg("cannot read reviews in archive")
		return nil, nil, err
	}

	// Pipe through to CSV
	reviewsFile, err := file.Open()
	if err != nil {
		log.Err(err).
			Str("zip", zipPath).
			Msg("could not open reviews.csv within zip")
		return nil, nil, err
	}
	closeFuncs = append(closeFuncs, reviewsFile.Close)

	// And now open the csv reader
	r = csv.NewReader(reviewsFile)
	return
}

func readReviewsCSV(csvReader *csv.Reader) ([]Review, error) {
	// Skip the header row
	csvReader.Read()

	// Read all rows, mapping along the way
	var reviews []Review
	for {
		// -> Get and check row
		row, err := csvReader.Read()
		if err == io.EOF {
			break
		}
		if len(row) == 0 {
			continue
		}
		if len(row) < 7 {
			err = fmt.Errorf("unexpected column count: %d", len(row))
			log.Err(err).Msg("malformed reviews.csv")
			return nil, err
		}

		date, err := civil.ParseDate(row[0])
		if err != nil {
			log.Err(err).
				Str("date", row[0]).
				Msg("malformed date column in reviews.csv")
			return nil, err
		}
		year, err := strconv.Atoi(row[2])
		if err != nil {
			log.Err(err).
				Str("year", row[2]).
				Msg("malformed year column in reviews.csv")
			return nil, err
		}
		var ratingF float64
		if row[4] != "" {
			starRating, err := strconv.ParseFloat(row[4], 64)
			if err != nil {
				log.Err(err).
					Str("rating", row[4]).
					Msg("malformed rating column in reviews.csv")
				return nil, err
			}
			ratingF = starRating * 2
			if ratingF != math.Trunc(ratingF) {
				err = errors.New("not a star rating. must go up in 0.5 increments")
				log.Err(err).
					Str("rating", row[4]).
					Msg("malformed rating column in reviews.csv")
				return nil, err
			}
		}
		rewatch := strings.EqualFold(row[5], "Yes")

		// Resolve some external info for it
		externalInfo := FetchData(row[3])

		review := Review{
			Date:          date,
			Name:          row[1],
			Year:          year,
			LetterboxdURI: row[3],
			Rating:        int(ratingF),
			Rewatch:       rewatch,
			Review:        row[6],
			PosterHref:    externalInfo.PosterHref,
		}
		reviews = append(reviews, review)
	}

	return reviews, nil
}
