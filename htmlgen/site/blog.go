package site

import (
	"fmt"
	"html/template"
	"time"

	"cloud.google.com/go/civil"
)

type BlogPost struct {
	Page
	Date time.Time
}

var BlogPosts []BlogPost

func blogPost(
	short string,
	title string,
	seoDesc string,
	date civil.Date,
	opening template.HTML,
	sections ...Section,
) BlogPost {
	t := date.In(time.Local)
	var allSections []Section
	allSections = append(allSections, section("", opening))
	allSections = append(allSections, sections...)
	page := page(rootTmpl, short,
		root(title, seoDesc,
			article(title,
				mul(withHeaderContent(markdown(fmt.Sprintf("*Written %s*", t.Format("2 January 2006"))))),
				allSections...,
			)))
	return BlogPost{
		Page: page,
		Date: t,
	}
}
