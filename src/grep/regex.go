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
	isInvert              bool
	isIgnoreCase          bool
	isFile                bool
	isFileAndLine         bool
	isFileAndLineAndMatch bool
	isNormal              bool
}

// find any kind of regex in a list of files and return grep style output
func findRegexInFiles(filesPaths []string, pattern string, flags *Flags) string {

	output := ""

	if flags.isIgnoreCase {
		pattern = fmt.Sprintf("(?i)%s", pattern)
	}

	re := regexp.MustCompile(pattern)

	for _, filePath := range filesPaths {
		file := openFile(filePath)
		defer file.Close()

		lineNumber := 1

		scanner := bufio.NewScanner(file)
		for scanner.Scan() {
			line := scanner.Text()

			matches := re.FindAllString(line, -1)
			numOfMatches := len(matches)

			if numOfMatches > 0 {
				line = colorString(line, "red", matches)
			}

			if flags.isNormal {
				if (numOfMatches > 0 && !flags.isInvert) || (numOfMatches == 0 && flags.isInvert) {
					output = fmt.Sprintf("%s%s:%s\n", output, filePath, line)
				}
			}

			if flags.isFile {
				if (numOfMatches > 0 && !flags.isInvert) || (numOfMatches == 0 && flags.isInvert) {
					output = fmt.Sprintf("%s%s\n", output, filePath)
					break
				}
			}

			if flags.isFileAndLine {
				if (numOfMatches > 0 && !flags.isInvert) || (numOfMatches == 0 && flags.isInvert) {
					output = fmt.Sprintf("%s%s:%d\n", output, filePath, lineNumber)
				}
			}

			if flags.isFileAndLineAndMatch {
				if (numOfMatches > 0 && !flags.isInvert) || (numOfMatches == 0 && flags.isInvert) {
					output = fmt.Sprintf("%s%s:%d:%s\n", output, filePath, lineNumber, line)
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
