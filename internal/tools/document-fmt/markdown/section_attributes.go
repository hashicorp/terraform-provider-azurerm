package markdown

import "strings"

type AttributesSection struct {
	heading Heading
	content []string
}

var _ SectionWithTemplate = &AttributesSection{}

func (s *AttributesSection) Match(line string) bool {
	return strings.Contains(strings.ToLower(line), "attributes")
}

func (s *AttributesSection) SetHeading(line string) {
	s.heading = NewHeading(line)
}

func (s *AttributesSection) GetHeading() Heading {
	return s.heading
}

func (s *AttributesSection) SetContent(content []string) {
	s.content = content
}

func (s *AttributesSection) GetContent() []string {
	return s.content
}

func (s *AttributesSection) Template() string {
	// TODO implement me
	panic("implement me")
}
