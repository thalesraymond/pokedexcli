package main

import (
	"testing"
)

func TestCleanInput(t *testing.T) {
	cases := []struct {
		input    string
		expected []string
	}{
		{
			input:    "  hello  world  ",
			expected: []string{"hello", "world"},
		},
		{
			input:    "  foo   bar   baz  ",
			expected: []string{"foo", "bar", "baz"},
		},
		{
			input:    "   ",
			expected: []string{},
		},
		{
			input:    "   singleword   ",
			expected: []string{"singleword"},
		},
		{
			input: "UPPERCASE  lowercase   MixedCase  ",
			expected: []string{"uppercase", "lowercase", "mixedcase"},
		},
	}

	for _, c := range cases {
		t.Run(c.input, func(t *testing.T) {
			result := cleanInput(c.input)
			if len(result) != len(c.expected) {
				t.Errorf("Expected %d elements, got %d", len(c.expected), len(result))
			}
			for i, v := range result {
				if v != c.expected[i] {
					t.Errorf("Expected %s at index %d, got %s", c.expected[i], i, v)
				}
			}
		})
	}	
}
