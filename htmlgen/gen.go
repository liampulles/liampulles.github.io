package main

import (
	"errors"
	"html/template"
	"os"
	"path/filepath"
	"time"

	"github.com/rs/zerolog/log"
)

// ---
// --- Templates
// ---

var (
	indexTmpl = loadTemplate("index.html")
	blogTmpl  = loadTemplate("blog.html")
)

func loadTemplate(file string) *template.Template {
	return template.Must(template.ParseFiles(
		filepath.Join("templates", "_elem.html"),
		filepath.Join("templates", file),
	))
}

func GenSite(outputFolder string) error {
	err := recreateFolder(outputFolder)
	if err != nil {
		return err
	}

	pages := wirePages(outputFolder)

	for _, page := range pages {
		gErr := gen(page)
		err = errors.Join(err, gErr)
	}
	return err
}

func recreateFolder(outputFolder string) error {
	err := os.RemoveAll(outputFolder)
	if err != nil {
		log.Err(err).
			Str("output_folder", outputFolder).
			Msg("could not clear and delete output folder, failing")
		return err
	}

	err = os.MkdirAll(outputFolder, os.ModePerm)
	if err != nil {
		log.Err(err).
			Str("dir", outputFolder).
			Msg("could not make output dir, failing")
		return err
	}

	return nil
}

func wirePages(outputFolder string) []wiredPage {
	// Wire site context
	site := siteContext{}
	for _, blog := range BlogPosts {
		site.blogs = append(site.blogs, blogElem{
			Date:  blog.ExtraData["Date"].(time.Time),
			Short: blog.Short,
			Title: string(blog.Title),
		})
	}

	// Ok, now we can create wired pages
	var wired []wiredPage
	wired = append(wired, wirePage(outputFolder, indexTmpl, IndexPage, site))
	for _, blog := range BlogPosts {
		wired = append(wired, wirePage(outputFolder, blogTmpl, blog, site))
	}

	return wired
}

type wiredPage struct {
	tmpl *template.Template
	loc  string
	data map[string]any
}

func wirePage(
	outputFolder string,
	tmpl *template.Template,
	page PageDefinition,
	site siteContext,
) wiredPage {
	// Construct data
	// -> Root
	data := make(map[string]any)
	data["Title"] = page.Title
	data["SEODescription"] = page.SEODescription
	data["NavElem"] = site.nav
	data["BlogPosts"] = site.blogs
	data["Year"] = time.Now().Year()
	// -> Extra
	for k, v := range page.ExtraData {
		data[k] = v
	}

	return wiredPage{
		tmpl: tmpl,
		loc:  filepath.Join(outputFolder, page.Short+".html"),
		data: data,
	}
}

type siteContext struct {
	nav   []navElem
	blogs []blogElem
}

type navElem struct {
	Short string
	Text  string
}

type blogElem struct {
	Date  time.Time
	Short string
	Title string
}

func gen(page wiredPage) error {
	// Open the file
	w, err := os.Create(page.loc)
	if err != nil {
		log.Err(err).
			Str("loc", page.loc).
			Msg("could not create output file, failing")
		return err
	}
	// -> Close it later
	defer func() {
		cErr := w.Close()
		if cErr != nil {
			log.Err(cErr).
				Str("loc", page.loc).
				Msg("could not close output file, failing")
		}
		err = errors.Join(err, cErr)
	}()

	// Template the HTML out to the file
	err = page.tmpl.ExecuteTemplate(w, "root", page.data)
	if err != nil {
		log.Err(err).
			Str("loc", page.loc).
			Msg("templating failed")
		return err
	}

	log.Debug().Str("file", page.loc).Msg("generated")
	return nil
}
