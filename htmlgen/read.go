package htmlgen

import (
	"fmt"
	"os"
	"path/filepath"
	"text/tabwriter"

	"github.com/rs/zerolog/log"
)

type ProjectView struct {
	PostFilepaths []MarkdownFile
}

// For debugging
func (p ProjectView) Log() {
	fmt.Fprintln(os.Stderr, "---PROJECT VIEW---")
	logMarkdownFiles("Posts", p.PostFilepaths)
}

func logMarkdownFiles(title string, files []MarkdownFile) {
	tw := tabwriter.NewWriter(os.Stderr, 1, 1, 1, ' ', 0)
	defer tw.Flush()
	fmt.Fprintf(os.Stderr, "  %s:\n", title)
	for _, file := range files {
		fmt.Fprintf(tw, "    | %s\t| %d bytes\t|\n",
			file.NamePart, len(file.Contents),
		)
	}
}

type MarkdownFile struct {
	NamePart string
	Contents []byte
}

const (
	postsSubfolder   = "_posts"
	markdownWildcard = "*.md"
)

// Read the root folder, identify key files, make a ProjectView.
func ReadProjectView(rootFolder string) (ProjectView, error) {
	// Read files
	postFiles, err := readMarkdownFilesByGlob(filepath.Join(rootFolder, postsSubfolder, markdownWildcard))
	if err != nil {
		log.Err(err).
			Str("root", rootFolder).
			Msg("could not glob for post files")
		return ProjectView{}, err
	}

	// Build project view
	return ProjectView{
		PostFilepaths: postFiles,
	}, nil
}

func readMarkdownFilesByGlob(glob string) ([]MarkdownFile, error) {
	// Look for files
	files, err := filepath.Glob(glob)
	if err != nil {
		log.Err(err).
			Str("glob", glob).
			Msg("could not glob for markdown files")
		return nil, err
	}

	// Read them into structs
	mdFiles := make([]MarkdownFile, len(files))
	for i, file := range files {
		mdFile, err := readMarkdownFile(file)
		if err != nil {
			log.Warn().
				Str("issue", err.Error()).
				Str("file", file).
				Msg("could not read markdown file - skipping")
			continue
		}
		mdFiles[i] = mdFile
	}

	return mdFiles, nil
}

func readMarkdownFile(loc string) (MarkdownFile, error) {
	// Parse path
	base := filepath.Base(loc)
	namePart := base[:len(base)-len(filepath.Ext(base))]

	// Read bytes
	b, err := os.ReadFile(loc)
	if err != nil {
		log.Err(err).
			Str("file", loc).
			Msg("could not read markdown file")
		return MarkdownFile{}, err
	}

	// Build struct
	return MarkdownFile{
		NamePart: namePart,
		Contents: b,
	}, nil
}
