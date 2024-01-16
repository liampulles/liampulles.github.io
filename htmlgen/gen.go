package htmlgen

import (
	"os"
	"path/filepath"

	"github.com/rs/zerolog/log"
)

const (
	postsOutFolder = "blog"
)

func GenSite(project ProjectView, outputFolder string) error {
	// Delete the output folder, to start from scratch.
	err := os.RemoveAll(outputFolder)
	if err != nil {
		log.Err(err).
			Str("output_folder", outputFolder).
			Msg("could not clear and delete output folder, failing")
		return err
	}

	// Parse, render, and write posts
	err = parseRenderWritePosts(project.PostFilepaths, outputFolder)
	if err != nil {
		return err
	}

	return nil
}

func parseRenderWritePosts(posts []MarkdownFile, outputFolder string) error {
	// Can be done post-by-post
	for _, post := range posts {
		parsed, err := ParseMarkdownish(post.Contents)
		if err != nil {
			log.Err(err).
				Str("post_name", post.NamePart).
				Msg("could not parse a post, skipping")
			continue
		}
		html := RenderPost(parsed)
		postFile := filepath.Join(outputFolder, postsOutFolder, post.NamePart+".html")
		err = writeFile(postFile, html)
		if err != nil {
			return err
		}
	}
	return nil
}

func writeFile(loc string, b []byte) error {
	// Ensure path exists
	dir := filepath.Dir(loc)
	err := os.MkdirAll(dir, os.ModePerm)
	if err != nil {
		log.Err(err).
			Str("dir", dir).
			Msg("could not make output dir, failing")
		return err
	}

	// Write file
	err = os.WriteFile(loc, b, 0755)
	if err != nil {
		log.Err(err).
			Str("file", loc).
			Msg("could not write file, failing")
		return err
	}

	return nil
}
