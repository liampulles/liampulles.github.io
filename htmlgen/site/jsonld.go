package site

import (
	"encoding/json"
	"fmt"
	"time"

	"cloud.google.com/go/civil"
)

type JSONld []byte

// See https://developers.google.com/search/docs/appearance/structured-data/article#json-ld
func JSONldBlogPosting(
	title string,
	image string,
	datePublished civil.Date,
) JSONld {
	// Pretend we wrote it at noon in SA
	timePublished := datePublished.
		In(time.FixedZone("SAST", int(+2*60*60))).
		Add(12 * time.Hour)
	m := map[string]any{
		"@context":      "https://schema.org",
		"@type":         "BlogPosting",
		"headline":      title,
		"image":         []string{fmt.Sprintf("https://liampulles.com/images/%s", image)},
		"datePublished": timePublished,
		"author": []map[string]any{
			{
				"@type":    "Person",
				"name":     "Liam Pulles",
				"url":      "https://liampulles.com/biography.html",
				"jobTitle": "Senior Software Engineer",
				"image":    "https://liampulles.com/images/profile.jpg",
			},
		},
	}

	bytes, err := json.Marshal(m)
	if err != nil {
		panic(fmt.Errorf("invalid json-ld data: %w", err))
	}
	return bytes
}
