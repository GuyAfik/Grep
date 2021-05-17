package main

import (
	"bufio"
	"fmt"
	"os"
	"regexp"
	"strings"

	"github.com/mitchellh/colorstring"
)

func openFile(filePath string) *os.File {
	file, err := os.Open(filePath)

	if err != nil {
		panic("unable to open file " + filePath)
	}

	return file
}

type Regex struct {
	files   []string
	pattern string
}

type Flags struct {
	isInvert                    bool
	isIgnoreCase                bool
	isFile                      bool
	isFileAndLineNumber         bool
	isFileAndLineNumberAndMatch bool
	isNormal                    bool
}

type Match struct {
	curOutput    string
	lineNumber   int
	filePath     string
	line         string
	foundMatches []string
}

func normalOutput() func(match *Match) string {
	return func(match *Match) string {
		match.line = colorString(match.line, "red", match.foundMatches)
		return fmt.Sprintf("%s%s:%s\n", match.curOutput, match.filePath, match.line)
	}
}

func fileOutput() func(match *Match) string {
	return func(match *Match) string {
		return fmt.Sprintf("%s%s\n", match.curOutput, match.filePath)
	}
}

func fileAndLineNumberOutput() func(match *Match) string {
	return func(match *Match) string {
		return fmt.Sprintf("%s%s:%d\n", match.curOutput, match.filePath, match.lineNumber)
	}
}

func fileAndLineNumberAndMatchOutput() func(match *Match) string {
	return func(match *Match) string {
		match.line = colorString(match.line, "red", match.foundMatches)
		return fmt.Sprintf("%s%s:%d:%s\n", match.curOutput, match.filePath, match.lineNumber, match.line)
	}
}

func chooseOperation(flags *Flags) func(match *Match) string {

	if flags.isFile {
		return fileOutput()
	} else if flags.isFileAndLineNumber {
		return fileAndLineNumberOutput()
	} else if flags.isFileAndLineNumberAndMatch {
		return fileAndLineNumberAndMatchOutput()
	} else {
		return normalOutput()
	}
}

// find any kind of regex in a list of files and return grep style output
func findRegexInFiles(filesPaths []string, pattern string, flags *Flags) string {

	output := ""

	if flags.isIgnoreCase {
		pattern = fmt.Sprintf("(?i)%s", pattern)
	}

	operation := chooseOperation(flags)

	re := regexp.MustCompile(pattern)

	for _, filePath := range filesPaths {
		file := openFile(filePath)
		defer file.Close()

		lineNumber := 1

		scanner := bufio.NewScanner(file)
		for scanner.Scan() {
			line := scanner.Text()

			regexMatches := re.FindAllString(line, -1)
			numOfMatches := len(regexMatches)
			match := &Match{
				curOutput:    output,
				lineNumber:   lineNumber,
				filePath:     filePath,
				line:         line,
				foundMatches: regexMatches,
			}

			if (numOfMatches > 0 && !flags.isInvert) || (numOfMatches == 0 && flags.isInvert) {
				output = operation(match)
				if flags.isFile {
					break
				}
			}

			if err := scanner.Err(); err != nil {
				fmt.Println(err)
			}

			lineNumber++
		}
	}

	return output
}

// color all the matching substrings that were found on a one line.
// returns that line where its matching substrings have been colored.
func colorString(line, color string, matches []string) string {

	coloredLine := line

	for _, match := range matches {
		coloredMatch := colorstring.Color(fmt.Sprintf("[%s]%s", color, match))
		coloredLine = strings.ReplaceAll(line, match, coloredMatch)
	}
	return coloredLine
}
