package site

import (
	"fmt"
	"html/template"
	"path/filepath"
	"strings"
	"time"
)

// The site files are essentially a blueprint for the site. What is defined here
// eventually makes its way into HTML.

// Other details are probably in the render functions.

var rootTmpl = loadTemplate(nil, "site.html")

var AllNavElem = []NavElem{
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

func nameToNav(name string) NavElem {
	short := strings.ReplaceAll(name, " ", "-")
	short = strings.ToLower(short)
	return NavElem{
		Short: short,
		Text:  name,
	}
}

type InlineMarkdown string

type PageType int

const (
	Index PageType = iota
	Nav
	BlogPost
	DigitalRestoration
)

type Page struct {
	Template *template.Template
	Type     PageType
	Short    string
	Data     Root
}

func page(
	tmpl *template.Template,
	typ PageType,
	short string,
	data Root,
) Page {
	return Page{
		Template: tmpl,
		Type:     typ,
		Short:    short,
		Data:     data,
	}
}

type Root struct {
	Title          string
	SEODescription string
	NavElem        []NavElem
	Article        Article
	Year           int
}

func root(
	title string,
	seoDesc string,
	article Article,
) Root {
	return Root{
		Title:          title,
		SEODescription: "SEODescription",
		NavElem:        AllNavElem,
		Article:        article,
		Year:           time.Now().Year(),
	}
}

type NavElem struct {
	Short string
	Text  string
}

type Article struct {
	Header   string
	Date     time.Time // Optional
	Elements []Element // In the header
	Sections []Section
}

func article(
	header string,
	elements []Element,
	sections []Section,
) Article {
	return Article{
		Header:   header,
		Elements: elements,
		Sections: sections,
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
	Elements []Element
}

func section(
	header string,
	asideFigures []Figure,
	elements []Element,
) Section {
	s := Section{
		Header:   header,
		Elements: elements,
	}
	s.Aside.Figures = asideFigures
	return s
}

type Figure struct {
	Images  []Image
	Caption string
}

func figure(
	images []Image,
	caption string,
) Figure {
	return Figure{
		Images:  images,
		Caption: caption,
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

type ElementType string

const (
	Paragraph     ElementType = "Paragraph"
	UnorderedList ElementType = "UnorderedList"
	OrderedList   ElementType = "OrderedList"
)

type Element struct {
	Type  ElementType
	Text  string
	Items []string
}

// Parses some light markdown into elements. Will log issues
// and panic if wrong.
func markdown(md string) []Element {
	//TODO:
	return []Element{
		{
			Type: Paragraph,
			Text: "TODO: markdown",
		},
	}
}
