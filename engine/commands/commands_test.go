package commands

import "testing"

func TestParseData(t *testing.T) {
	tests := []struct {
		args []string
		expected []string
		expectedResult bool
	}{{[]string{"just", "a", "test"}, []string{"s", "s", "s"}, true},
	{[]string{"just", "a", "test"}, []string{"s", "s", "s", "s"}, false},
	{[]string{"just", "a", "test", "9"}, []string{"s", "s", "s", "i"}, true},
	{[]string{"000", "000", "000", "000"}, []string{"i", "i", "i", "i"}, true},}
	for _, test := range tests {
		if parseData(test.args, test.expected) != test.expectedResult {
			t.Fail()
		}
	}
}

func TestToSimpleText(t *testing.T) {
	tests := []struct {
		text string
		filter []string
		expectedResult string
	}{
		{"First test!", []string{"!"}, "First test"},
		{"Second test?!", []string{"!", "?"}, "Second test"},
		{"Third test? Oh, no!", []string{"!"}, "Third test? Oh, no"},
		{"...", []string{"."}, ""},
	}
	for _, test := range tests {
		if toSimpleText(test.text, test.filter) != test.expectedResult {
			t.Fail()
		}
	}
}

func TestGetRandomNum(t *testing.T) {
	tests := []struct {
		first string
		second string
	}{
		{"test1", "test2"},
		{" ", ""},
		{"123", "6"},
	}
	for _, test := range tests {
		if getRandomNum(test.first) == getRandomNum(test.second) {
			t.Fail()
		}
	}
}