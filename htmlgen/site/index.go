package site

var indexTmpl = loadTemplate(rootTmpl, "index.html")

var IndexPage = PageDefinition{
	Template:       indexTmpl,
	Type:           Index,
	Short:          "index",
	Title:          "Liam Pulles",
	SEODescription: "Homepage for Liam Pulles's blog.",
	ExtraData: map[string]any{
		"MainHeader": "Welcome!",
		"MainDescription": InlineMarkdown(`Hi there - if you're interested in my writing, read on. 
		    If you want to hire me (or otherwise find out more about me), then you may wish to see
			my [[biography]] or my [[code]].`),
	},
}
