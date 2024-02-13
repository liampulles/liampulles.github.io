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
	outputFlag := fs.String("output", "_site_test", "folder to render the outputted \"site\" to")
	if err := fs.Parse(os.Args[1:]); err != nil {
		log.Err(err).Msg("arg parse fail")
		os.Exit(1)
	}

	// Run the program
	err := run(*outputFlag)
	if err != nil {
		os.Exit(2)
	}
}

func run(
	outputFolder string,
) error {
	// Generate the site
	err := htmlgen.GenSite(outputFolder)
	if err != nil {
		return err
	}

	log.Info().
		Str("output_folder", outputFolder).
		Msg("site generated!")
	return nil
}
