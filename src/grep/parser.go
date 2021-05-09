package main

import (
	"fmt"
	"os"

	"github.com/akamensky/argparse"
)

// parses all the grep input arguments. 
// returns a Flags struct indicating what flags should be performed 
// and a Regex struct indicating what files and pattern should be searched.
func parseGrepArgs() (*Regex, *Flags) {

	parser := argparse.NewParser("grep", "parser used for parsing grep arguments")

	ignoreCase := parser.Flag(
		"i",
		"--ignore-case",
		&argparse.Options{
			Required: false,
			Help:     "Ignore case distinctions, so that characters that differ only in case match each other.",
		},
	)
	count := parser.Flag(
		"c",
		"count",
		&argparse.Options{
			Required: false,
			Help:     "Suppress  normal output; instead print a count of matching lines for each input file.  With the -v, --invert-match option (see below), count non-matching lines.",
		},
	)
	linesAndFiles := parser.Flag(
		"n",
		"line-number",
		&argparse.Options{
			Required: false,
			Help:     "Prefix each line of output with the 1-based line number within its input file.",
		},
	)
	invert := parser.Flag(
		"v",
		"--invert-match",
		&argparse.Options{
			Required: false,
			Help:     "Invert the sense of matching, to select non-matching lines.",
		},
	)
	filesWithMatches := parser.Flag(
		"l",
		"--files-with-matches",
		&argparse.Options{
			Required: false,
			Help:     "Suppress normal output; instead print the name of each input file from which output would normally have been  printed. The  scanning  will stop on the first match.",
		},
	)

	files := parser.StringList("f", "files", &argparse.Options{Required: true, Help: "Files that pattern should be searched in"})
	pattern := parser.String("e", "pattern", &argparse.Options{Required: true, Help: "The regex/pattern to search for"})

	err := parser.Parse(os.Args)

	if err != nil {
		fmt.Println(parser.Usage(err))
		os.Exit(1)
	}

	flags := &Flags{}
	regex := &Regex{}

	regex.files = *files
	regex.pattern = *pattern

	flags.isIgnoreCase = *ignoreCase
	flags.isInvert = *invert
	flags.isFile = *filesWithMatches
	flags.isFileAndLine = *count && !flags.isFile
	flags.isFileAndLineAndMatch = *linesAndFiles && !flags.isFile && !flags.isFileAndLine
	flags.isNormal = !flags.isFile && !flags.isFileAndLineAndMatch && !flags.isFileAndLine

	return regex, flags
}

