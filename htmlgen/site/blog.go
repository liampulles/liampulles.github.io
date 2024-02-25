package site

import (
	"fmt"
	"html/template"
	"time"

	"cloud.google.com/go/civil"
)

type DatedPost struct {
	Page
	Date     time.Time
	Unlisted bool
}

var BlogPosts []DatedPost

func blogPost(
	short string,
	title string,
	seoDesc string,
	date civil.Date,
	opening template.HTML,
	sections ...Section,
) DatedPost {
	t := date.In(time.Local)

	// Insert an opening section
	var allSections []Section
	allSections = append(allSections, section("", opening))
	allSections = append(allSections, sections...)

	page := page(rootTmpl, short,
		root(title, seoDesc,
			article(title,
				mul(
					withHeaderContent(markdown(fmt.Sprintf(
						"*Written %s*",
						t.Format("2 January 2006"),
					))),
				),
				allSections...,
			),
			withCommentsFooter(short),
		))

	return DatedPost{
		Page: page,
		Date: t,
	}
}
