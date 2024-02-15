package main

import (
	"errors"
	"os"
	"path"

	"github.com/liampulles/liampulles.github.io/htmlgen/site"
	"github.com/rs/zerolog/log"
)

func GenSite(outputFolder string) error {
	err := recreateFolder(outputFolder)
	if err != nil {
		return err
	}

	err = errors.Join(err, gen(outputFolder, site.IndexPage()))
	for _, blogPage := range site.BlogPosts {
		err = errors.Join(err, gen(outputFolder, blogPage.Page))
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

func gen(outputFolder string, page site.Page) error {
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
