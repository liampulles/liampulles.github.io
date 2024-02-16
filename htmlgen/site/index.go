package site

import (
	"fmt"
	"html/template"
	"time"
)

func IndexPage() Page {
	return page(rootTmpl, "index", root(
		"Liam Pulles",
		"Homepage for Liam Pulles's blog.",
		article("Welcome!", "",
			section("", markdown(`
Hi there - if you're interested in my writing, read on. 
If you want to hire me (or otherwise find out more about me), then you may wish to see
my [[biography]] or my [[code]].`)),
			section("Blog posts", blogPostsTable()),
		),
	),
	)
}

func blogPostsTable() template.HTML {
	var t [][]any
	for _, blogPost := range BlogPosts {
		t = append(t, []any{
			template.HTML(fmt.Sprintf("<b>%s</b>", blogPost.Date.In(time.Local).Format("Jan 02, 2006"))),
			template.HTML(fmt.Sprintf(`<b><a href="/%s.html">%s</a></b>`, blogPost.Short, blogPost.Data.Title)),
		})
	}
	return table(t)
}
