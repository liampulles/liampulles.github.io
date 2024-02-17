package site

import (
	"errors"
	"fmt"
	"html/template"
	"regexp"
	"strings"
)

// Snippets are little isolated pieces of information which the site can link to
// dynamically.

var ValidSnippetShorts = map[string]struct{}{
	"digital-restorations": {},
}

var Snippets = []SnippetPage{
	snippetPage("digital-restorations", "Digital restorations", markdown(`
Digital restorations are pieces of artwork that I have cleaned and upscaled,
to make them more suitable for high-DPI printing.

I'm no expert at this - its more of a hobby really.`)),
}

type SnippetPage struct {
	Template *template.Template
	Short    string
	Header   string
	Content  template.HTML
}

func snippetPage(
	short string,
	header string,
	content template.HTML,
) SnippetPage {
	// Verify the short
	_, found := ValidSnippetShorts[short]
	if !found {
		panic(fmt.Errorf("no known snippet short for page (%s)", short))
	}

	return SnippetPage{
		Template: rootTmpl,
		Short:    short,
		Header:   header,
		Content:  content,
	}
}

var snippetLinkRegex = regexp.MustCompile(`\?:(.*)\?`)

func replaceSnippetLinks(s string) (string, error) {
	var err error
	s = snippetLinkRegex.ReplaceAllStringFunc(s, func(submatch string) string {
		// Extract
		elem := snippetLinkRegex.FindStringSubmatch(s)
		title, short, hasShort := strings.Cut(elem[1], ":")
		if !hasShort {
			short = inferShort(title)
		}

		// Verify
		if len(title) == 0 {
			err = errors.Join(err, fmt.Errorf("title must not be empty on snippet (%s)", submatch))
			return submatch
		}
		_, ok := ValidSnippetShorts[short]
		if !ok {
			err = errors.Join(err, fmt.Errorf("no known snippet short for %q (%s)", short, submatch))
			return submatch
		}

		// Format
		return fmt.Sprintf(`<a class=snippet-link hx-get=/snippet/%s.html hx-swap=afterend>%s</a>`, short, title)
	})
	return s, err
}

func inferShort(title string) string {
	// Remove html wrappings
	for strings.Contains(title, "<") {
		open := strings.Index(title, "<")
		close := strings.Index(title, ">")
		title = title[:open] + title[close+1:]
	}

	s := strings.TrimSpace(title)
	s = strings.ToLower(s)
	s = strings.ReplaceAll(s, " ", "-")
	return s
}
