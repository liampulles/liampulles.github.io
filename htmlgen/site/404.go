package site

func NotFoundPage() Page {
	return page(rootTmpl, "404", root(
		"Page Not Found",
		"The page given in the URL does not exist on liampulles.com",
		article("Page Not Found", nil),
	))
}
