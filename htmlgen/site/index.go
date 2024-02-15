package site

var indexTmpl = loadTemplate(rootTmpl, "index.html")

var IndexPage = page(indexTmpl, Index, "index",
	root("Liam Pulles", "Homepage for Liam Pulles's blog.",
		article("Welcome!", nil, mul(
			section("", nil, markdown(
				`Hi there - if you're interested in my writing, read on. 
				 If you want to hire me (or otherwise find out more about me), then you may wish to see
		         my [[biography]] or my [[code]].`,
			))))))
