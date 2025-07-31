package markdown

import "strings"

type ExampleSection struct {
	heading Heading
	content []string
}

var _ SectionWithTemplate = &ExampleSection{}

func (s *ExampleSection) Match(line string) bool {
	line = strings.ToLower(line)
	// some docs contain "<Name> Usage" rather than "Example Usage"
	return strings.Contains(line, "example") || strings.Contains(line, "usage")
}

func (s *ExampleSection) SetHeading(line string) {
	s.heading = NewHeading(line)
}

func (s *ExampleSection) GetHeading() Heading {
	return s.heading
}

func (s *ExampleSection) SetContent(content []string) {
	s.content = content
}

func (s *ExampleSection) GetContent() []string {
	return s.content
}

func (s *ExampleSection) Template() string {
	// TODO implement me
	panic("implement me")
}
