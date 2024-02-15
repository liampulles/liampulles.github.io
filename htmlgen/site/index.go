package site

import (
	"fmt"
	"html/template"
	"time"
)

var indexTmpl = loadTemplate(rootTmpl, "index.html")

func IndexPage() Page {
	return page(indexTmpl, Index, "index",
		root("Liam Pulles", "Homepage for Liam Pulles's blog.",
			article("Welcome!", "", mul(
				section("", nil, markdown(`
Hi there - if you're interested in my writing, read on. 
If you want to hire me (or otherwise find out more about me), then you may wish to see
my [biography](/biography.html) or my [code](/code.html).`)),
				section("Blog posts", nil, blogPostsTable()),
			))))
}

func blogPostsTable() template.HTML {
	var t [][]any
	for _, blogPost := range BlogPosts {
		t = append(t, []any{
			blogPost.Date.In(time.Local).Format("Jan 02, 2006"),
			template.HTML(fmt.Sprintf(`<a href="/%s.html">%s</a>`, blogPost.Short, blogPost.Data.Title)),
		})
	}
	return table(t)
}
