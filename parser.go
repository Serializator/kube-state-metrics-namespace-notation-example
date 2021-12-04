package parser

import (
	"errors"
	"strings"
)

var syntaxError = errors.New("invalid format, resource=[label=value,label=value,...],resource=[]")

type ResourceWithLabelFiltersNotationParser struct {
	source string
	result map[string]map[string]string
}

func (parser ResourceWithLabelFiltersNotationParser) Parse() (map[string]map[string]string, error) {
	// taken from the text/scanner EOF constant
	const EOF = -1

	var (
		identOrValueStartPos int
		nextRune             rune
		enclosedByIdentifier bool
		enclosingIdentifier string
		labelIdentifier     string
	)

	for i, currentRune := range parser.source {
		if i+1 < len(parser.source) {
			nextRune = []rune(parser.source)[i+1]
		} else {
			nextRune = EOF
		}

		switch currentRune {
		case '=':
			if enclosedByIdentifier {
				if nextRune == ']' {
					// indicates a "my-project-*=[label=]" structure
					return nil, syntaxError
				}
				// in this case we found a "=" after a label ("my-project-*=[label=]")
				identifier := strings.TrimSpace(string([]rune(parser.source)[identOrValueStartPos:i]))
				parser.result[enclosingIdentifier][identifier] = ""
				labelIdentifier = identifier
				identOrValueStartPos = i+1
			} else {
				// my-project-*=[ <- validate that "=" is followed by a "["
				if nextRune != '[' {
					return nil, syntaxError
				}

				identifier := strings.TrimSpace(string([]rune(parser.source)[identOrValueStartPos:i]))
				parser.result[identifier] = map[string]string{}
				enclosingIdentifier = identifier
				identOrValueStartPos = i+1
			}
		case '[':
			if enclosedByIdentifier {
				return nil, syntaxError
			}
			enclosedByIdentifier = true
			identOrValueStartPos = i+1
		case ',':
			if enclosedByIdentifier {
				value := strings.TrimSpace(string([]rune(parser.source)[identOrValueStartPos:i]))
				parser.result[enclosingIdentifier][labelIdentifier] = value

				labelIdentifier = ""
				identOrValueStartPos = i+1
			} else {
				identifier := strings.TrimSpace(string([]rune(parser.source)[identOrValueStartPos:i]))
				parser.result[identifier] = map[string]string{}
				identOrValueStartPos = i+1
			}
		case ']':
			if !enclosedByIdentifier {
				return nil, syntaxError
			}
			value := strings.TrimSpace(string([]rune(parser.source)[identOrValueStartPos:i]))
			parser.result[enclosingIdentifier][labelIdentifier] = value

			enclosedByIdentifier = false
			enclosingIdentifier = ""
			labelIdentifier = ""
		default:
			if nextRune == EOF && !enclosedByIdentifier {
				identifier := strings.TrimSpace(string([]rune(parser.source)[identOrValueStartPos:i+1]))
				parser.result[identifier] = map[string]string{}
				// no need to update the "identOrValueStartPos" as we are at the end
			}
		}
	}

	return parser.result, nil
}

func NewResourceWithLabelFiltersNotationParser(value string) ResourceWithLabelFiltersNotationParser {
	return ResourceWithLabelFiltersNotationParser{
		source: value,
		result: map[string]map[string]string{},
	}
}