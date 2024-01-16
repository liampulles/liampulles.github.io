package main

import (
	"flag"
	"os"

	"github.com/liampulles/liampulles.github.io/htmlgen"
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
)

// Program flow:
// - Read the root folder to build a project view
// - Parse markdownish files into AST
// - Render the site according to internal rules and any instructions given
//   in markdownish.

func main() {
	// Setup logging
	zerolog.TimeFieldFormat = zerolog.TimeFormatUnix
	zerolog.SetGlobalLevel(zerolog.DebugLevel)
	log.Logger = log.Output(zerolog.ConsoleWriter{Out: os.Stderr})

	// Parse flags
	fs := flag.NewFlagSet("htmlgen", flag.ContinueOnError)
	rootFlag := fs.String("root", ".", "root folder of the site tree")
	outputFlag := fs.String("output", "_site_test", "folder to render the outputted \"site\" to")
	if err := fs.Parse(os.Args[1:]); err != nil {
		log.Err(err).Msg("arg parse fail")
		os.Exit(1)
	}

	// Run the program
	err := run(*rootFlag, *outputFlag)
	if err != nil {
		os.Exit(2)
	}
}

func run(
	rootFolder string,
	outputFolder string,
) error {
	// Read the project
	projectView, err := htmlgen.ReadProjectView(rootFolder)
	if err != nil {
		return err
	}
	log.Debug().Interface("view", projectView).Msg("read project view")

	return nil
}
