package topic_parser

import (
	"regexp"
)

type TopicParser struct {
	regex *regexp.Regexp
}

func (tp *TopicParser) ParseTopic(topic string) map[string]string {
	result := make(map[string]string)

	matches := tp.regex.FindStringSubmatch(topic)
	names := tp.regex.SubexpNames()

	for i := 1; i < len(matches); i++ {
		result[names[i]] = matches[i]
	}

	return result
}

func New(topicFmt string) TopicParser {
	re, err := regexp.Compile(convertPath(topicFmt))
	if err != nil {
		panic(err)
	}

	return TopicParser{re}
}
