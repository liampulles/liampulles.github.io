package site

import (
	"html/template"
	"sort"
	"strings"
	"time"

	"github.com/liampulles/liampulles.github.io/htmlgen/letterboxd"
)

func ReviewsPage() Page {
	return page(rootTmpl, "reviews", root(
		"Reviews",
		"Large compilation of film reviews written by me, Liam Pulles.",
		article("Film Reviews", mul(withRawContent(reviewsPageContent()))),
	))
}

type ReviewsPageContent struct {
	Years []ReviewYear
}

type ReviewYear struct {
	Year    int
	Reviews []Review
}

type Review struct {
	Stars         template.HTML
	Name          string
	Year          int
	DateReviewed  time.Time
	Review        template.HTML
	LetterboxdURI string
	PosterHref    string
}

func reviewsPageContent() template.HTML {
	// Read letterboxd data
	export, err := letterboxd.ReadExport()
	if err != nil {
		panic(err)
	}

	// Sort reviews by date, latest first
	sort.Slice(export.Reviews, func(i, j int) bool {
		return export.Reviews[i].Date.After(export.Reviews[j].Date)
	})

	// Map to data format
	var data ReviewsPageContent
	var currentReviewYear ReviewYear
	for _, review := range export.Reviews {
		// -> New year
		if review.Date.Year != currentReviewYear.Year {
			if len(currentReviewYear.Reviews) > 0 {
				data.Years = append(data.Years, currentReviewYear)
			}
			currentReviewYear = ReviewYear{Year: review.Date.Year}
		}

		// Skip?
		if shouldSkipReview(review) {
			continue
		}

		// Fix review text
		reviewText := preFixReviewText(review.Review)

		// -> Add this review to this year
		r := Review{
			Stars:         starRating(review.Rating),
			Name:          review.Name,
			Year:          review.Year,
			DateReviewed:  review.Date.In(time.UTC),
			Review:        markdown(reviewText),
			LetterboxdURI: review.LetterboxdURI,
			PosterHref:    review.PosterHref,
		}
		currentReviewYear.Reviews = append(currentReviewYear.Reviews, r)
	}

	// Template
	return execTemplate(rootTmpl, "reviews", data)
}

func starRating(rating int) template.HTML {
	// Make the stars out of 8
	if rating > 8 {
		rating = 8
	}

	fullStars := rating / 2
	halfStar := rating%2 == 1

	s := strings.Repeat(`<i class="fa-solid fa-star"></i>`, fullStars)
	if halfStar {
		s += `<i class="fa-solid fa-star-half"></i>`
	}
	return template.HTML(s)
}

func shouldSkipReview(review letterboxd.Review) bool {
	// Skip crawley film series
	if strings.Contains(review.Review, "industrial-films-crawley-films-ranked") {
		return true
	}

	return false
}

func preFixReviewText(review string) string {
	// If the first line contains a bracketed section, omit it.
	// We'll style this ourselves.
	firstLine, rest, _ := strings.Cut(review, "\n")
	firstLine = strings.TrimSpace(firstLine)
	if strings.HasSuffix(firstLine, ")") || strings.HasSuffix(firstLine, "]") {
		firstLine = ""
	}
	review = firstLine + "\n" + rest

	return review
}
