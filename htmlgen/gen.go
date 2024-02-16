package main

import (
	"errors"
	"fmt"
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
	// Delete and make the folder
	err := recreateFolder(outputFolder)
	if err != nil {
		return err
	}

	// Do some HTML templating, and stylesheet writing
	// -> HTML
	var jobs []jobFn
	jobs = append(jobs, genJob(outputFolder, site.IndexPage()))
	jobs = append(jobs, genJob(outputFolder, site.BiographyPage))
	for _, blogPage := range site.BlogPosts {
		jobs = append(jobs, genJob(outputFolder, blogPage.Page))
	}
	// -> CSS
	jobs = append(jobs, writeStyle(outputFolder, "monokai", "dark"))
	jobs = append(jobs, writeStyle(outputFolder, "tango", "light"))
	return doAll(jobs...)
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

func genJob(outputFolder string, page site.Page) jobFn {
	return func() error {
		// Open the file
		loc := path.Join(outputFolder, page.Short+".html")
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

		// Template the HTML out to the file
		err = page.Template.ExecuteTemplate(w, "root", page.Data)
		if err != nil {
			log.Err(err).
				Str("loc", loc).
				Msg("templating failed")
			return err
		}

		log.Debug().Str("file", loc).Msg("generated")
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
