package markdown

type TitleSection struct {
	heading Heading
	content []string
}

var _ SectionWithTemplate = &TitleSection{}

func (s *TitleSection) Match(line string) bool {
	panic("todo")
}

func (s *TitleSection) SetHeading(line string) {
	s.heading = NewHeading(line)
}

func (s *TitleSection) GetHeading() Heading {
	return s.heading
}

func (s *TitleSection) SetContent(content []string) {
	s.content = content
}

func (s *TitleSection) GetContent() []string {
	return s.content
}

func (s *TitleSection) Template() string {
	// TODO implement me
	panic("implement me")
}
