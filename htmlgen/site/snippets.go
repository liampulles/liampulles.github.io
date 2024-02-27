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
	"rating-system":        {},
}

var Snippets = []SnippetPage{
	snippetPage("digital-restorations", "Digital restorations", markdown(`
Digital restorations are pieces of artwork that I have cleaned and upscaled,
to make them more suitable for high-DPI printing.

I'm no expert at this - its more of a hobby really.`)),

	snippetPage("rating-system", "Rating system", template.HTML(`
<table>
	<tr>
		<td class="stars">(Zero)</td>
		<td>=</td>
		<td>Did Not Finish</td>
	</tr>
	<tr>
		<td class="stars"><i class="fa-solid fa-star-half"></i></td>
		<td>=</td>
		<td>Practically Unwatchable</td>
	</tr>
	<tr>
		<td class="stars"><i class="fa-solid fa-star"></i></td>
		<td>=</td>
		<td>Horrible</td>
	</tr>
	<tr>
		<td class="stars"><i class="fa-solid fa-star"></i><i class="fa-solid fa-star-half"></i></td>
		<td>=</td>
		<td>Bad</td>
	</tr>
	<tr>
		<td class="stars"><i class="fa-solid fa-star"></i><i class="fa-solid fa-star"></i></td>
		<td>=</td>
		<td>Moderately Bad</td>
	</tr>
	<tr>
		<td class="stars"><i class="fa-solid fa-star"></i><i class="fa-solid fa-star"></i><i class="fa-solid fa-star-half"></i></td>
		<td>=</td>
		<td>Ok</td>
	</tr>
	<tr>
		<td class="stars"><i class="fa-solid fa-star"></i><i class="fa-solid fa-star"></i><i class="fa-solid fa-star"></i></td>
		<td>=</td>
		<td>Good</td>
	</tr>
	<tr>
		<td class="stars"><i class="fa-solid fa-star"></i><i class="fa-solid fa-star"></i><i class="fa-solid fa-star"></i><i class="fa-solid fa-star-half"></i></td>
		<td>=</td>
		<td>Very Good</td>
	</tr>
	<tr>
		<td class="stars"><i class="fa-solid fa-star"></i><i class="fa-solid fa-star"></i><i class="fa-solid fa-star"></i><i class="fa-solid fa-star"></i></td>
		<td>=</td>
		<td>Great</td>
	</tr>
</table>
`)),
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

var snippetLinkRegex = regexp.MustCompile(`\?:([^?]*)\?`)

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
