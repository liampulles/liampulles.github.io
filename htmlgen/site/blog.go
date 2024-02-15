package site

import (
	"fmt"
	"time"

	"cloud.google.com/go/civil"
)

var blogTmpl = loadTemplate(rootTmpl, "blog.html")

var BlogPosts []Page

func blogPost(
	short string,
	title string,
	seoDesc string,
	date civil.Date,
	opening []Element,
	sections []Section,
) Page {
	t := date.In(time.Local)
	var allSections []Section
	allSections = append(allSections, section("", nil, opening))
	allSections = append(allSections, sections...)
	return page(blogTmpl, BlogPost, short,
		root(title, seoDesc,
			article(title,
				markdown(fmt.Sprintf("*Written %s*", t.Format("2 January 2006"))),
				allSections,
			)))
}
