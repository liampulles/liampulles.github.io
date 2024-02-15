package site

import (
	"time"

	"cloud.google.com/go/civil"
)

var blogTmpl = loadTemplate(rootTmpl, "blog.html")

var BlogPosts []PageDefinition

func blogPost(
	short string,
	title string,
	seoDesc string,
	date civil.Date,
) PageDefinition {
	return PageDefinition{
		Template:       blogTmpl,
		Type:           BlogPost,
		Short:          short,
		Title:          InlineMarkdown(title),
		SEODescription: InlineMarkdown(seoDesc),
		ExtraData: map[string]any{
			"Date": date.In(time.Local),
		},
	}
}
