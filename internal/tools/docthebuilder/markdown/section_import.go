package markdown

import "strings"

type ImportSection struct {
	heading Heading
	content []string
}

var _ SectionWithTemplate = &ImportSection{}

func (s *ImportSection) Match(line string) bool {
	return strings.Contains(strings.ToLower(line), "import")
}

func (s *ImportSection) SetHeading(line string) {
	s.heading = NewHeading(line)
}

func (s *ImportSection) GetHeading() Heading {
	return s.heading
}

func (s *ImportSection) SetContent(content []string) {
	s.content = content
}

func (s *ImportSection) GetContent() []string {
	return s.content
}

func (s *ImportSection) Template() string {
	panic("implement me")
}
