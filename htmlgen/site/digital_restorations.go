package site

import (
	"fmt"
	"html/template"
	"time"

	"cloud.google.com/go/civil"
)

var DigitalRestorations []DatedPost

func digitalRestoration(
	short string,
	title string,
	seoDesc string,
	date civil.Date,
	linkImage LinkImage,
	description template.HTML,
) DatedPost {
	t := date.In(time.Local)
	page := page(rootTmpl, short,
		root(title, seoDesc,
			article(title, mul(
				withHeaderContent(markdown(fmt.Sprintf("*Written %s*", t.Format("2 January 2006")))),
				withRawContent(restorationPage(linkImage, description)),
			)),
			withCommentsFooter(short),
			withJSONld(JSONldBlogPosting(title, string(linkImage.Link), date)),
		))
	return DatedPost{
		Page: page,
		Date: t,
	}
}

type RestorationPage struct {
	LinkImage   LinkImage
	Description template.HTML
}

func restorationPage(linkImage LinkImage, description template.HTML) template.HTML {
	return execTemplate(rootTmpl, "restoration-page", RestorationPage{
		LinkImage:   linkImage,
		Description: description,
	})
}
