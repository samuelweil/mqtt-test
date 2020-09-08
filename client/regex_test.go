package main

import (
	"testing"
)

func checkString(t *testing.T, exp string, got string) {
	if exp != got {
		t.Errorf("Expected \"%s\", got \"%s\"", exp, got)
	}
}

func TestPathRegex(t *testing.T) {
	result := pathConvert.FindAllStringSubmatch("{first}/{second}/{third}", -1)

	if len(result) != 3 {
		t.Fatalf("Regex failed to return 4 values, got %d", len(result))
	}

	checkString(t, "first", result[0][1])
	checkString(t, "second", result[1][1])
	checkString(t, "third", result[2][1])
}
func TestConvertPath(t *testing.T) {

	tests := []struct {
		inp, exp string
	}{
		{
			"/{first}/{second}/{third}",
			"/(?P<first>.*)/(?P<second>.*)/(?P<third>.*)",
		},
	}

	for _, test := range tests {
		checkString(t, test.exp, convertPath(test.inp))
	}

}

func TestTopicChange(t *testing.T) {

	tests := []struct {
		inp, exp string
	}{
		{
			"{first}/{second}/messages/",
			"+/+/messages/",
		},
	}

	for _, test := range tests {
		checkString(t, test.exp, topic(test.inp))
	}

}
