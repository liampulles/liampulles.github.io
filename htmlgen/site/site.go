package site

import (
	"html/template"
	"path/filepath"
)

// The site files are essentially a blueprint for the site. What is defined here
// eventually makes its way into HTML.

// Other details are probably in the render functions.

var rootTmpl = loadTemplate(nil, "site.html")

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

type InlineMarkdown string

type PageType int

const (
	Index PageType = iota
	Nav
	BlogPost
	DigitalRestoration
)

type PageDefinition struct {
	Template       *template.Template
	Type           PageType
	Short          string
	Title          InlineMarkdown
	SEODescription InlineMarkdown
	ExtraData      map[string]any
}
