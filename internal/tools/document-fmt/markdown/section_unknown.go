package markdown

type SectionUnknown struct {
	heading Heading
	content []string
}

var _ Section = &SectionUnknown{}

func (s *SectionUnknown) Match(_ string) bool {
	return false // special case, this should never match as this is a catch-all
}

func (s *SectionUnknown) SetHeading(line string) {
	s.heading = NewHeading(line)
}

func (s *SectionUnknown) GetHeading() Heading {
	return s.heading
}

func (s *SectionUnknown) SetContent(content []string) {
	s.content = content
}

func (s *SectionUnknown) GetContent() []string {
	return s.content
}
