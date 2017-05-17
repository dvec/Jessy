package tools

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
		if CheckData(test.args, test.expected) != test.expectedResult {
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
		if GetRandomNum(test.first) == GetRandomNum(test.second) {
			t.Fail()
		}
	}
}

func TestContains(t *testing.T) {
	tests := []struct {
		arr []string
		value string
	}{
		{[]string{"1", "2", "3"}, "2"},
		{[]string{"", "", ""}, ""},
		{[]string{"@", "\\", "/"}, "/"},
	}

	for _, test := range tests {
		if !Contains(test.arr, test.value) {
			t.Fail()
		}
	}
}