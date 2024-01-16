package htmlgen

import (
	"path/filepath"

	"github.com/rs/zerolog/log"
)

type ProjectView struct {
	PostFilepaths []string
}

const (
	postsSubfolder   = "_posts"
	markdownWildcard = "*.md"
)

// Read the root folder, identify key files, make a ProjectView.
func ReadProjectView(rootFolder string) (ProjectView, error) {
	// Look for files
	postFiles, err := filepath.Glob(filepath.Join(rootFolder, postsSubfolder, markdownWildcard))
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
