package main

import (
	"errors"
	"os"
	"path"
	"sync"

	"github.com/liampulles/liampulles.github.io/htmlgen/site"
	"github.com/rs/zerolog/log"
)

func GenSite(outputFolder string) error {
	err := recreateFolder(outputFolder)
	if err != nil {
		return err
	}

	var jobs []func() error
	jobs = append(jobs, genJob(outputFolder, site.IndexPage()))
	jobs = append(jobs, genJob(outputFolder, site.BiographyPage))
	for _, blogPage := range site.BlogPosts {
		jobs = append(jobs, genJob(outputFolder, blogPage.Page))
	}
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

func genJob(outputFolder string, page site.Page) func() error {
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

func doAll(jobs ...func() error) (err error) {
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
