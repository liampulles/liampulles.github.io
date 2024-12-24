package site

import "html/template"

// Make a custom not found page, which gives suggestions on what the page might be.
// Since this is dynamic functionality, the meat of the work is done by javascript.
// We just create the HTML shell here.
func NotFoundPage() Page {
	target := template.HTML(`<div id="maybePages"></div>`)

	return page(rootTmpl, "404", root(
		"Page Not Found",
		"The page given in the URL does not exist on liampulles.com",
		article("Page Not Found", mul(withRawContent(target))),
	))

}
