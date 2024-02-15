package site

import (
	"fmt"
	"html/template"
	"time"

	"cloud.google.com/go/civil"
)

var blogTmpl = loadTemplate(rootTmpl, "blog.html")

type BlogPost struct {
	Page
	Date civil.Date
}

var BlogPosts []BlogPost

func blogPost(
	short string,
	title string,
	seoDesc string,
	date civil.Date,
	opening template.HTML,
	sections []Section,
) BlogPost {
	t := date.In(time.Local)
	var allSections []Section
	allSections = append(allSections, section("", nil, opening))
	allSections = append(allSections, sections...)
	page := page(blogTmpl, BlogPostType, short,
		root(title, seoDesc,
			article(title,
				markdown(fmt.Sprintf("*Written %s*", t.Format("2 January 2006"))),
				allSections,
			)))
	return BlogPost{
		Page: page,
		Date: date,
	}
}
