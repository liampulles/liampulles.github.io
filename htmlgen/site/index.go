package site

func IndexPage() Page {
	return page(rootTmpl, "index", root(
		"Liam Pulles",
		"Homepage for Liam Pulles's blog.",
		article("Welcome!", mul(withRawContent(indexTOC()))),
	))
}
