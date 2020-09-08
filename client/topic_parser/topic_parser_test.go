package topic_parser

import "testing"

func checkMap(t *testing.T, exp map[string]string, got map[string]string) {
	for k, v := range exp {
		val, ok := got[k]
		if !ok {
			t.Errorf("Expected key %s, not found", k)
			continue
		}

		if val != v {
			t.Errorf("Expected %s: %s, got %s", k, v, val)
		}
	}

	if extras := checkExtraKeys(exp, got); len(extras) > 0 {
		for _, e := range extras {
			t.Errorf("Unexpected key %v found", e)
		}
	}

}

func checkExtraKeys(exp map[string]string, got map[string]string) []string {

	result := make([]string, 0)

	for k, _ := range got {
		if _, ok := exp[k]; !ok {
			result = append(result, k)
		}
	}

	return result
}

func TestTopicParser(t *testing.T) {

	tests := []struct {
		topicFmt string
		topic    string
		exp      map[string]string
	}{
		{
			"devices/{deviceId}",
			"devices/bob",
			map[string]string{
				"deviceId": "bob",
			},
		}, {
			"/{x}/{y}",
			"/hello/world",
			map[string]string{
				"x": "hello",
				"y": "world",
			},
		},
	}

	for _, test := range tests {
		parser := New(test.topicFmt)
		result := parser.ParseTopic(test.topic)

		checkMap(t, test.exp, result)
	}
}
