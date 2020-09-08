package topic_parser

import (
	"fmt"
	"regexp"
)

var pathConvert = regexp.MustCompile("{(.*?)}")

func convertPath(path string) string {
	return pathConvert.ReplaceAllStringFunc(path, func(ident string) string {
		return fmt.Sprintf("(?P<%s>.*)", ident[1:len(ident)-1])
	})
}

func topic(origTopic string) string {
	return pathConvert.ReplaceAllString(origTopic, "+")
}
