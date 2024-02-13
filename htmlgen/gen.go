package main

import (
	"errors"
	"html/template"
	"io"
	"os"
	"path"
	"path/filepath"
	"time"

	"github.com/rs/zerolog/log"
)

// ---
// --- Templates
// ---

var indexTmpl = loadTemplate("index.html")

func loadTemplate(file string) *template.Template {
	return template.Must(template.ParseFiles(
		filepath.Join("templates", "_elem.html"),
		filepath.Join("templates", file),
	))
}

// ---
// --- Data funcs
// ---

func rootData(
	pageTitle string,
	pageDescription string,
) map[string]any {
	return map[string]any{
		"PageTitle":       pageTitle,
		"PageDescription": pageDescription,
		"NavElem":         NavElems,
		"Year":            time.Now().Year(),
	}
}

func mergeData(data map[string]any, from ...map[string]any) map[string]any {
	if data == nil {
		data = make(map[string]any)
	}
	for _, m := range from {
		for k, v := range m {
			data[k] = v
		}
	}
	return data
}

// ---
// --- Generating
// ---

func GenSite(outputFolder string) error {
	// Delete the output folder, to start from scratch.
	err := os.RemoveAll(outputFolder)
	if err != nil {
		log.Err(err).
			Str("output_folder", outputFolder).
			Msg("could not clear and delete output folder, failing")
		return err
	}

	// Gen files
	return errors.Join(
		gen(indexTmpl, outputFolder, "index.html", "Liam Pulles"),
	)
}

func gen(tmpl *template.Template, outputFolder, file, title string, extraData ...map[string]any) error {
	p := path.Join(outputFolder, file)
	data := mergeData(nil, rootData(title, title))
	data = mergeData(data, extraData...)
	err := withFile(p, func(w io.Writer) error {
		return tmpl.ExecuteTemplate(w, "root", data)
	})
	if err != nil {
		return err
	}

	log.Debug().Str("file", p).Msg("generated")
	return nil
}

func withFile(loc string, do func(w io.Writer) error) (err error) {
	// Ensure path exists
	dir := filepath.Dir(loc)
	err = os.MkdirAll(dir, os.ModePerm)
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

	// Use the file
	return do(w)
}
