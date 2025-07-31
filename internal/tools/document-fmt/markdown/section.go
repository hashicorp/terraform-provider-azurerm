package markdown

import (
	"fmt"
	"reflect"
	"slices"
	"strings"
)

type Section interface {
	// Match returns whether the provided heading text matches the expected section heading
	Match(string) bool

	SetHeading(string)
	GetHeading() Heading

	SetContent([]string)
	GetContent() []string
}

type SectionWithTemplate interface {
	Section
	// Template should return a string able to be parsed as a `text/template` template.
	// - when rendered, the entire `data.ResourceData` struct is available to the template, so placeholders such as `{{ .Name }}` can be used.
	// - backticks can be included using `[bt]`, these will all be replaced with a literal backtick (`).
	Template() string
}

// InsertAfterSection inserts a provided section after a specified section and returns an updated Section slice.
// if a matching section wasn't found, it returns an error
func InsertAfterSection(newSection Section, sections []Section, after Section) ([]Section, error) {
	if len(sections) == 0 {
		return nil, fmt.Errorf("received an empty sections slice")
	}

	for idx, s := range sections {
		if reflect.TypeOf(s) == reflect.TypeOf(after) {
			sections = slices.Insert(sections, idx+1, newSection)
			return sections, nil
		}
	}

	return nil, fmt.Errorf("did not find a section of type `%T`", after)
}

// InsertBeforeSection inserts a provided section before a specified section and returns an updated Section slice.
// if a matching section wasn't found, it returns an error
func InsertBeforeSection(newSection Section, sections []Section, before Section) ([]Section, error) {
	if len(sections) == 0 {
		return nil, fmt.Errorf("received an empty sections slice")
	}

	for idx, s := range sections {
		if reflect.TypeOf(s) == reflect.TypeOf(before) {
			sections = slices.Insert(sections, idx, newSection)
			return sections, nil
		}
	}

	return nil, fmt.Errorf("did not find a section of type `%T`", before)
}

func FindSectionByHeading(sections []Section, search string) Section {
	for _, s := range sections {
		h := s.GetHeading().Text
		if strings.Contains(strings.ToLower(h), strings.ToLower(search)) { // case-insensitive for flexibility
			return s
		}
	}

	return nil
}
