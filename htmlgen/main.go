package main

import (
	"flag"
	"os"
	"time"

	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
)

// Program flow:
// - Read the root folder to build a project view
// - Parse markdownish files into AST
// - Render the site according to internal rules and any instructions given
//   in markdownish.

func main() {
	// Measure time
	start := time.Now()
	defer func() {
		log.Debug().Msgf("%v start to finish", time.Since(start))
		log.Debug().Msg("--------------------------------------")
	}()

	// Setup logging
	zerolog.TimeFieldFormat = zerolog.TimeFormatUnix
	zerolog.SetGlobalLevel(zerolog.DebugLevel)
	log.Logger = log.Output(zerolog.ConsoleWriter{Out: os.Stderr})
	log.Debug().Msg("--------------------------------------")

	// Parse flags
	fs := flag.NewFlagSet("htmlgen", flag.ContinueOnError)
	outputFlag := fs.String("output", "_site_gen", "folder to render the outputted \"site\" to")
	if err := fs.Parse(os.Args[1:]); err != nil {
		log.Err(err).Msg("arg parse fail")
		os.Exit(1)
	}

	// Run the program
	// Generate the site
	err := GenSite(*outputFlag)
	if err != nil {
		os.Exit(2)
	}

	log.Info().
		Str("output_folder", *outputFlag).
		Msg("site generated!")
}
