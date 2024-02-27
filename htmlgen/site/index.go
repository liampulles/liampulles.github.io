package site

import (
	"html/template"
	"sort"
)

func IndexPage() Page {
	return page(rootTmpl, "index", root(
		"Liam Pulles",
		"Homepage for Liam Pulles's blog.",
		article("Welcome!", mul(withRawContent(indexTOC()))),
	))
}

type IndexTOC struct {
	BlogPosts           []DatedPost
	DigitalRestorations []DatedPost
}

func indexTOC() template.HTML {
	// Only keep listed items
	var blogPosts []DatedPost
	for _, post := range BlogPosts {
		if post.Unlisted {
			continue
		}
		blogPosts = append(blogPosts, post)
	}
	var digRestores []DatedPost
	for _, post := range DigitalRestorations {
		if post.Unlisted {
			continue
		}
		digRestores = append(digRestores, post)
	}

	// Sort posts
	sort.Slice(blogPosts, func(i, j int) bool {
		return blogPosts[i].Date.After(blogPosts[j].Date)
	})
	sort.Slice(digRestores, func(i, j int) bool {
		return digRestores[i].Date.After(digRestores[j].Date)
	})

	data := IndexTOC{
		BlogPosts:           blogPosts,
		DigitalRestorations: digRestores,
	}
	return execTemplate(rootTmpl, "index-toc", data)
}
