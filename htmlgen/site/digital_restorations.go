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
	image Image,
	description template.HTML,
) DatedPost {
	t := date.In(time.Local)
	page := page(rootTmpl, short,
		root(title, seoDesc,
			article(title, mul(
				withHeaderContent(markdown(fmt.Sprintf("*Written %s*", t.Format("2 January 2006")))),
				withRawContent(restorationPage(image, description)),
			)),
			withCommentsFooter(short),
		))
	return DatedPost{
		Page: page,
		Date: t,
	}
}

type RestorationPage struct {
	Image       Image
	Description template.HTML
}

func restorationPage(image Image, description template.HTML) template.HTML {
	return execTemplate(rootTmpl, "restoration-page", RestorationPage{
		Image:       image,
		Description: description,
	})
}
