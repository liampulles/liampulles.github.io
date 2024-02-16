package site

import (
	"fmt"
	"html/template"
	"path/filepath"
	"strings"
	"time"

	chromahtml "github.com/alecthomas/chroma/v2/formatters/html"
	"github.com/rs/zerolog/log"
	"github.com/yuin/goldmark"
	highlighting "github.com/yuin/goldmark-highlighting/v2"
	goldhtml "github.com/yuin/goldmark/renderer/html"
	"go.abhg.dev/goldmark/wikilink"
)

// The site files are essentially a blueprint for the site. What is defined here
// eventually makes its way into HTML.

// Other details are probably in the render functions.

const liveURL = "https://liampulles.com"

var rootTmpl = loadTemplate(nil, "_tmpl.html")

var allNavElem = []NavElem{
	nameToNav("Biography"),
	nameToNav("Proverbs"),
	nameToNav("Code"),
}

// ---
// --- Helpers
// ---

func loadTemplate(root *template.Template, file string) *template.Template {
	if root == nil {
		return template.Must(template.ParseFiles(filepath.Join("htmlgen", "site", file)))
	}
	t := template.Must(root.Clone())
	return template.Must(t.ParseFiles(filepath.Join("htmlgen", "site", file)))
}

func execTemplate(t *template.Template, name string, data any) template.HTML {
	var sb strings.Builder
	err := t.ExecuteTemplate(&sb, name, data)
	if err != nil {
		log.Err(err).Msgf("could not exec template %s", name)
		panic(fmt.Errorf("could not exec template %s: %w", name, err))
	}
	return template.HTML(sb.String())
}

func nameToNav(name string) NavElem {
	short := strings.ReplaceAll(name, " ", "-")
	short = strings.ToLower(short)
	return NavElem{
		Short: short,
		Text:  name,
	}
}

type Page struct {
	Template *template.Template
	Short    string
	Data     Root
}

func page(
	tmpl *template.Template,
	short string,
	data Root,
) Page {
	return Page{
		Template: tmpl,
		Short:    short,
		Data:     data,
	}
}

type Root struct {
	Title          string
	SEODescription string
	NavElem        []NavElem
	Article        Article
	Footer         Footer
}

func root(
	title string,
	seoDesc string,
	article Article,
	opts ...func(*Root),
) Root {
	r := Root{
		Title:          title,
		SEODescription: "SEODescription",
		NavElem:        allNavElem,
		Article:        article,
		Footer: Footer{
			Year: time.Now().Year(),
		},
	}

	for _, opt := range opts {
		if opt == nil {
			continue
		}
		opt(&r)
	}

	return r
}

func withConnectWithMe(r *Root) {
	r.Footer.ConnectWithMe = true
}

func withComments(short string) func(r *Root) {
	return func(r *Root) {
		c := Comments{
			FullURL: fmt.Sprintf("%s/%s.html", liveURL, short),
		}
		r.Footer.Comments = &c
	}
}

type Footer struct {
	ConnectWithMe bool
	Year          int
	Comments      *Comments
}

type Comments struct {
	FullURL string
}

type NavElem struct {
	Short string
	Text  string
}

type Article struct {
	Header        string
	HeaderContent template.HTML
	RawContent    template.HTML
	Sections      []Section
}

func article(
	header string,
	opts []func(*Article),
	sections ...Section,
) Article {
	a := Article{
		Header:   header,
		Sections: sections,
	}

	for _, opt := range opts {
		if opt == nil {
			continue
		}
		opt(&a)
	}

	return a
}

func withHeaderContent(headerContent template.HTML) func(a *Article) {
	return func(a *Article) {
		a.HeaderContent = headerContent
	}
}

func withRawContent(content template.HTML) func(a *Article) {
	return func(a *Article) {
		a.RawContent = content
	}
}

func mul[T any](many ...T) []T {
	return many
}

type Section struct {
	Header string
	Aside  struct {
		Figures []Figure
	}
	Content template.HTML
}

func section(
	header string,
	content template.HTML,
	asideFigures ...Figure,
) Section {
	s := Section{
		Header:  header,
		Content: content,
	}
	s.Aside.Figures = asideFigures
	return s
}

type Figure struct {
	Images  []Image
	Caption template.HTML
}

func figure(
	caption string,
	images ...Image,
) Figure {
	return Figure{
		Images:  images,
		Caption: template.HTML(caption),
	}
}

type Image struct {
	Src    string
	Width  string
	Height string
	Alt    string
}

func image(
	file string,
	width int,
	height int,
	alt string,
) Image {
	return Image{
		Src:    fmt.Sprintf("/images/%s", file),
		Width:  fmt.Sprintf("%dpx", width),
		Height: fmt.Sprintf("%dpx", height),
		Alt:    alt,
	}
}

var md = goldmark.New(
	goldmark.WithExtensions(
		// E.g. [[biography]] type links
		&wikilink.Extender{},
		// Applies syntax highlighting via https://github.com/alecthomas/chroma
		highlighting.NewHighlighting(
			highlighting.WithFormatOptions(
				chromahtml.WithClasses(true),
			),
		),
	),
	goldmark.WithRendererOptions(
		goldhtml.WithUnsafe(),
	),
)

// Parses markdown into HTML. Will log issues and panic if wrong.
//
// Don't go crazy - just use light elements. Should ideally use above DSL
// for structuring.
func markdown(s string) template.HTML {
	var sb strings.Builder
	err := md.Convert([]byte(s), &sb)
	if err != nil {
		log.Err(err).Msgf("could not parse markdown (%s...)", head(s, 20))
		panic(fmt.Errorf("could not parse markdown: %w", err))
	}
	return template.HTML(sb.String())
}

func head(s string, count int) string {
	runes := []rune(s)
	if count > len(runes) {
		count = len(runes)
	}
	return string(runes[:count])
}

type IndexTOC struct {
	BlogPosts           []BlogPost
	DigitalRestorations []BlogPost // TODO
}

func indexTOC() template.HTML {
	data := IndexTOC{
		BlogPosts: BlogPosts,
	}
	return execTemplate(rootTmpl, "index-toc", data)
}

// Emits markdown with code wrapping, taking care of certain requirements
// needed for this site
func codeFigureMarkdown(lang string, code string) string {
	// Replace any backtick-hacks
	code = strings.ReplaceAll(code, "\\'", "`")
	return fmt.Sprintf(`
<figure class="highlight">

~~~%s
%s
~~~

</figure>
`, lang, code)
}
