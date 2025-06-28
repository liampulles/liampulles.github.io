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
	"sync"

	"cloud.google.com/go/civil"
	"github.com/liampulles/liampulles.github.io/htmlgen/parallel"
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
	// As map reader
	r := headerReader(csvReader)

	// Read rows and create jobs
	var reviews []Review
	var jobs []parallel.Job
	var mu sync.Mutex
	for {
		// -> Get and check row
		var row map[string]string
		row, err := r.Read()
		if err == io.EOF {
			break
		}
		jobs = append(jobs, func() error {
			review, ok, err := readReviewCSVRow(row)
			if err != nil {
				return err
			}
			if !ok {
				return nil
			}

			mu.Lock()
			reviews = append(reviews, review)
			mu.Unlock()

			return nil
		})
	}

	// Run in parallel
	err := parallel.Concurrent(jobs, 1)
	if err != nil {
		return nil, err
	}

	return reviews, nil
}

func readReviewCSVRow(row map[string]string) (Review, bool, error) {
	if len(row) == 0 {
		return Review{}, false, nil
	}
	if len(row) < 7 {
		err := fmt.Errorf("unexpected column count: %d", len(row))
		log.Err(err).Msg("malformed reviews.csv")
		return Review{}, false, err
	}

	strDate := row["Watched Date"]
	if strDate == "" {
		strDate = row["Date"]
	}
	date, err := civil.ParseDate(strDate)
	if err != nil {
		log.Err(err).
			Str("date", strDate).
			Msg("malformed date column in reviews.csv")
		return Review{}, false, err
	}
	year, err := strconv.Atoi(row["Year"])
	if err != nil {
		log.Err(err).
			Str("year", row["Year"]).
			Msg("malformed year column in reviews.csv")
		return Review{}, false, err
	}
	var ratingF float64
	if row["Rating"] != "" {
		starRating, err := strconv.ParseFloat(row["Rating"], 64)
		if err != nil {
			log.Err(err).
				Str("rating", row["Rating"]).
				Msg("malformed rating column in reviews.csv")
			return Review{}, false, err
		}
		ratingF = starRating * 2
		if ratingF != math.Trunc(ratingF) {
			err = errors.New("not a star rating. must go up in 0.5 increments")
			log.Err(err).
				Str("rating", row["Rating"]).
				Msg("malformed rating column in reviews.csv")
			return Review{}, false, err
		}
	}
	rewatch := strings.EqualFold(row["Rewatch"], "Yes")

	// Resolve some external info
	externalInfo := FetchData(row["Letterboxd URI"])

	review := Review{
		Date:          date,
		Name:          row["Name"],
		Year:          year,
		LetterboxdURI: row["Letterboxd URI"],
		Rating:        int(ratingF),
		Rewatch:       rewatch,
		Review:        row["Review"],
		PosterHref:    externalInfo.PosterHref,
	}
	return review, true, nil
}

type csvHeaderReader struct {
	r      *csv.Reader
	header []string
}

func headerReader(r *csv.Reader) csvHeaderReader {
	// Read header and map
	header, err := r.Read()
	if err != nil {
		log.Fatal().Err(err).
			Msg("could not read header row")
	}

	return csvHeaderReader{
		r:      r,
		header: header,
	}
}

func (hr csvHeaderReader) Read() (map[string]string, error) {
	row, err := hr.r.Read()
	if err != nil {
		return nil, err
	}

	m := make(map[string]string, len(row))
	for i, s := range row {
		m[hr.header[i]] = s
	}
	return m, nil
}
