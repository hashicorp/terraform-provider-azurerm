package markdown

import "strings"

type ArgumentsSection struct {
	heading Heading
	content []string
}

var _ SectionWithTemplate = &ArgumentsSection{}

func (s *ArgumentsSection) Match(line string) bool {
	return strings.Contains(strings.ToLower(line), "arguments")
}

func (s *ArgumentsSection) SetHeading(line string) {
	s.heading = NewHeading(line)
}

func (s *ArgumentsSection) GetHeading() Heading {
	return s.heading
}

func (s *ArgumentsSection) SetContent(content []string) {
	s.content = content
}

func (s *ArgumentsSection) GetContent() []string {
	return s.content
}

func (s *ArgumentsSection) Template() string {
	// TODO implement me
	panic("implement me")
}
