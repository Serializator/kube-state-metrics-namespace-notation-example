package parser

import (
	"reflect"
	"testing"
)

type TestCase struct {
	input    string
	expected map[string]map[string]string
}

func TestResourceWithLabelFiltersNotationParser_Parse(t *testing.T) {
	tests := []TestCase{
		{
			input: "*=[environment=production,zone=europe]",
			expected: map[string]map[string]string{
				"*": {
					"environment": "production",
					"zone": "europe",
				},
			},
		},
		{
			input: "my-project-*=[region=eu-west-2]",
			expected: map[string]map[string]string{
				"my-project-*": {
					"region": "eu-west-2",
				},
			},
		},
		{
			// for the sake of backward compatibility of the old "--namespaces" notation
			input: "foobar",
			expected: map[string]map[string]string{
				"foobar": {},
			},
		},
		{
			// for the sake of backward compatibility of the old "--namespaces" notation
			input: "foobar1,foobar2,foobar3",
			expected: map[string]map[string]string{
				"foobar1": {},
				"foobar2": {},
				"foobar3": {},
			},
		},
	}

	for _, test := range tests {
		result, err := NewResourceWithLabelFiltersNotationParser(test.input).Parse()
		if err != nil {
			t.Errorf("unexpected error whilst parsing: %v", err)
		}
		if !reflect.DeepEqual(test.expected, result) {
			t.Errorf("unexpected result from parser, got: %v, expected: %v", result, test.expected)
		}
	}
}