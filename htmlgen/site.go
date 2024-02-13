package htmlgen

// This file is essentially a blueprint for the site. What is defined here
// eventually makes its way into HTML.

// Other details are probably in the render functions.

type InlineMarkdown string

// Main page config
const (
	MainPageHeader InlineMarkdown = "Welcome!"
	MainPageText   InlineMarkdown = `
	Hi there - if you're interested in my writing, read on. If you want to hire me (or otherwise
	find out more about me), then you may wish to see my [[biography]] or my [[code]].`
	BlogPostsHeader           InlineMarkdown = "Blog posts"
	DigitalRestorationsHeader InlineMarkdown = "Digital restorations"
	PostTimeFormat                           = "2006-01-02"
)
