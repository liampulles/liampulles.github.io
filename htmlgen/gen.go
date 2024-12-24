package main

import (
	"encoding/json"
	"encoding/xml"
	"errors"
	"fmt"
	"io"
	"os"
	"path"
	"path/filepath"
	"strings"
	"sync"

	chromahtml "github.com/alecthomas/chroma/v2/formatters/html"
	"github.com/alecthomas/chroma/v2/styles"
	"github.com/liampulles/liampulles.github.io/htmlgen/site"
	"github.com/rs/zerolog/log"
)

func GenSite(outputFolder string) error {
	// Delete the folder, start fresh
	err := deleteFolder(outputFolder)
	if err != nil {
		return err
	}

	// Do some HTML templating, and stylesheet writing
	// -> HTML
	var jobs []jobFn
	indexPage := site.IndexPage()
	notFoundPage := site.NotFoundPage()
	jobs = append(jobs, genJob(outputFolder, notFoundPage.Short, page(notFoundPage)))
	jobs = append(jobs, genJob(outputFolder, indexPage.Short, page(indexPage)))
	jobs = append(jobs, genJob(outputFolder, site.BiographyPage.Short, page(site.BiographyPage)))
	reviewsPage := site.ReviewsPage()
	jobs = append(jobs, genJob(outputFolder, reviewsPage.Short, page(reviewsPage)))
	for _, blogPage := range site.BlogPosts {
		jobs = append(jobs, genJob(outputFolder, blogPage.Page.Short, page(blogPage.Page)))
	}
	for _, dr := range site.DigitalRestorations {
		jobs = append(jobs, genJob(outputFolder, dr.Page.Short, page(dr.Page)))
	}
	for _, r := range site.RedirectPages {
		jobs = append(jobs, genJob(outputFolder, r.Short, redirect(r)))
	}
	for _, s := range site.Snippets {
		jobs = append(jobs, genJob(outputFolder, "snippet/"+s.Short, snippet(s)))
	}
	// -> CSS
	jobs = append(jobs, writeStyle(outputFolder, "monokai", "dark"))
	jobs = append(jobs, writeStyle(outputFolder, "tango", "light"))
	// -> Sitemap
	jobs = append(jobs, writeSitemap(outputFolder))
	// -> Javascript
	jobs = append(jobs, writeMaybePages(outputFolder))

	return doAll(jobs...)
}

func deleteFolder(outputFolder string) error {
	err := os.RemoveAll(outputFolder)
	if err != nil {
		log.Err(err).
			Str("output_folder", outputFolder).
			Msg("could not clear and delete output folder, failing")
		return err
	}

	return nil
}

func genJob(outputFolder string, short string, with withFile) jobFn {
	return func() error {
		loc := path.Join(outputFolder, short+".html")

		// Make folder
		dir := filepath.Dir(loc)
		err := os.MkdirAll(dir, os.ModePerm)
		if err != nil {
			log.Err(err).
				Str("dir", dir).
				Msg("could not make output dir, failing")
			return err
		}

		// Open the file
		w, err := os.Create(loc)
		if err != nil {
			log.Err(err).
				Str("loc", loc).
				Msg("could not create output file, failing")
			return err
		}
		// -> Close it later
		defer func() {
			cErr := w.Close()
			if cErr != nil {
				log.Err(cErr).
					Str("loc", loc).
					Msg("could not close output file, failing")
			}
			err = errors.Join(err, cErr)
		}()

		// Use file
		err = with(w)
		if err != nil {
			return err
		}

		log.Debug().Str("file", loc).Msg("generated")
		return nil
	}
}

type withFile func(io.Writer) error

// Should only contain meaningful, "actual" pages (no redirects)
var sitemapShorts []string

func writeSitemap(outputFolder string) jobFn {
	return func() error {
		// Make folder
		err := os.MkdirAll(outputFolder, os.ModePerm)
		if err != nil {
			log.Err(err).
				Str("dir", outputFolder).
				Msg("could not make output dir, failing")
			return err
		}

		loc := path.Join(outputFolder, "sitemap.xml")
		// Open the file
		w, err := os.Create(loc)
		if err != nil {
			log.Err(err).
				Str("loc", loc).
				Msg("could not create output file, failing")
			return err
		}
		// -> Close it later
		defer func() {
			cErr := w.Close()
			if cErr != nil {
				log.Err(cErr).
					Str("loc", loc).
					Msg("could not close output file, failing")
			}
			err = errors.Join(err, cErr)
		}()

		// Template XML
		type sitemapURL struct {
			Location string `xml:"loc"`
		}
		var sitemapURLSet struct {
			XMLName xml.Name     `xml:"http://www.sitemaps.org/schemas/sitemap/0.9 urlset"`
			URLs    []sitemapURL `xml:"url"`
		}
		for _, short := range sitemapShorts {
			sitemapURLSet.URLs = append(sitemapURLSet.URLs, sitemapURL{
				Location: fmt.Sprintf("%s/%s.html", site.LiveURL, short),
			})
		}
		xmlBytes, err := xml.Marshal(sitemapURLSet)
		if err != nil {
			log.Err(err).
				Str("loc", loc).
				Msg("could not marshal sitemap xml")
			return fmt.Errorf("invalid sitemap xml: %w", err)
		}

		// Write it to file
		_, err = w.Write(xmlBytes)
		if err != nil {
			log.Err(err).
				Str("loc", loc).
				Msg("could not write sitemap xml")
			return fmt.Errorf("could not write sitemap xml: %w", err)
		}

		return nil
	}
}

func writeMaybePages(outputFolder string) jobFn {
	return func() error {
		// Make folder
		err := os.MkdirAll(outputFolder, os.ModePerm)
		if err != nil {
			log.Err(err).
				Str("dir", outputFolder).
				Msg("could not make output dir, failing")
			return err
		}

		loc := path.Join(outputFolder, "maybe_pages.js")
		// Open the file
		w, err := os.Create(loc)
		if err != nil {
			log.Err(err).
				Str("loc", loc).
				Msg("could not create output file, failing")
			return err
		}
		// -> Close it later
		defer func() {
			cErr := w.Close()
			if cErr != nil {
				log.Err(cErr).
					Str("loc", loc).
					Msg("could not close output file, failing")
			}
			err = errors.Join(err, cErr)
		}()

		// Build a set of suggestions
		type maybePage struct {
			Location string `json:"location"`
			Title    string `json:"title"`
		}
		var maybes []maybePage
		for _, post := range site.BlogPosts {
			maybes = append(maybes, maybePage{
				Location: fmt.Sprintf("%s/%s.html", site.LiveURL, post.Short),
				Title:    post.Page.Data.Title,
			})
		}
		for _, post := range site.DigitalRestorations {
			maybes = append(maybes, maybePage{
				Location: fmt.Sprintf("%s/%s.html", site.LiveURL, post.Short),
				Title:    post.Page.Data.Title,
			})
		}

		// Template out a javascript array
		arrBytes, err := json.Marshal(maybes)
		if err != nil {
			log.Err(err).
				Str("loc", loc).
				Msg("could not write sitemap xml")
			return fmt.Errorf("could not write json array %w", err)
		}
		array := fmt.Sprintf("var maybe_pages = %s", string(arrBytes))

		// Write it to file
		_, err = w.WriteString(array)
		if err != nil {
			log.Err(err).
				Str("loc", loc).
				Msg("could not write sitemap xml")
			return fmt.Errorf("could not write sitemap xml: %w", err)
		}

		return nil
	}
}

func page(p site.Page) withFile {
	// We consider all "pages" to be actual places users can navigate to,
	// so include them in the sitemap.
	sitemapShorts = append(sitemapShorts, p.Short)

	return func(w io.Writer) error {
		err := p.Template.ExecuteTemplate(w, "root", p.Data)
		if err != nil {
			log.Err(err).
				Msg("templating failed")
			return err
		}
		return nil
	}
}

func redirect(p site.RedirectPage) withFile {
	return func(w io.Writer) error {
		err := p.Template.ExecuteTemplate(w, "redirect", p)
		if err != nil {
			log.Err(err).
				Msg("templating failed")
			return err
		}
		return nil
	}
}

func snippet(p site.SnippetPage) withFile {
	return func(w io.Writer) error {
		err := p.Template.ExecuteTemplate(w, "snippet", p)
		if err != nil {
			log.Err(err).
				Msg("templating failed")
			return err
		}
		return nil
	}
}

type jobFn func() error

func doAll(jobs ...jobFn) (err error) {
	var wg sync.WaitGroup
	for i := range jobs {
		job := jobs[i]
		wg.Add(1)
		go func() {
			jErr := job()
			err = errors.Join(err, jErr)
			wg.Done()
		}()
	}
	wg.Wait()
	return err
}

func writeStyle(outputFolder, name, lightDark string) jobFn {
	return func() error {
		// We're going to want to prepend some stuff to the stylesheet
		// to make ti work with my dark/light toggle. So write to string first.
		var sb strings.Builder
		formatter := chromahtml.New(
			chromahtml.WithLineNumbers(true),
			chromahtml.WithClasses(true),
		)
		s := styles.Get(name)
		err := formatter.WriteCSS(&sb, s)
		if err != nil {
			log.Err(err).Str("style", name).Msg("chroma could not write style")
			return err
		}
		css := sb.String()

		// Now add the light dark mode after the comment on each line
		css = strings.ReplaceAll(css, "*/", fmt.Sprintf(`*/ :root[color-mode="%s"]`, lightDark))

		// Make folder
		err = os.MkdirAll(outputFolder, os.ModePerm)
		if err != nil {
			log.Err(err).
				Str("dir", outputFolder).
				Msg("could not make output dir, failing")
			return err
		}

		// Now we can write it out
		loc := filepath.Join(outputFolder, lightDark+".css")
		err = os.WriteFile(loc, []byte(css), 0664)
		if err != nil {
			log.Err(err).Str("loc", loc).Msg("could not write stylesheet")
			return err
		}

		log.Debug().Str("file", loc).Msg("generated")
		return nil
	}
}
