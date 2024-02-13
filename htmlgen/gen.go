package htmlgen

import (
	"errors"
	"html/template"
	"io"
	"os"
	"path"
	"path/filepath"

	"github.com/rs/zerolog/log"
)

// ---
// --- Templates
// ---

var indexTmpl = loadTemplate("index.html")

func loadTemplate(file string) *template.Template {
	return template.Must(template.ParseFiles(
		filepath.Join("htmlgen", "templates", "_elem.html"),
		filepath.Join("htmlgen", "templates", file),
	))
}

// ---
// --- Data funcs
// ---

func indexData() any {
	return map[string]any{
		"PageTitle":       "test page",
		"PageDescription": "test description",
		"NavElem": []map[string]any{
			{"Link": "/some-link", "Text": "Some Link"},
			{"Link": "/some-link2", "Text": "Some Link 2"},
		},
		"Year": "Some year",
	}
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
		gen(outputFolder, "index.html", indexTmpl, indexData()),
	)
}

func gen(outputFolder, name string, tmpl *template.Template, data any) error {
	p := path.Join(outputFolder, name)
	return withFile(p, func(w io.Writer) error {
		return tmpl.ExecuteTemplate(w, "root", data)
	})
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
