package main

import (
	"fmt"
)

func main() {
	regex, flags := parseGrepArgs()
	fmt.Println(findRegexInFiles(regex.files, regex.pattern, flags))
}
