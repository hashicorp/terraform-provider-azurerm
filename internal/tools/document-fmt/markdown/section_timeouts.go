package markdown

import (
	"strings"
)

type TimeoutsSection struct {
	heading Heading
	content []string
}

var _ SectionWithTemplate = &TimeoutsSection{}

func (s *TimeoutsSection) Match(line string) bool {
	return strings.Contains(strings.ToLower(line), "timeout")
}

func (s *TimeoutsSection) SetHeading(line string) {
	s.heading = NewHeading(line)
}

func (s *TimeoutsSection) GetHeading() Heading {
	return s.heading
}

func (s *TimeoutsSection) SetContent(content []string) {
	s.content = content
}

func (s *TimeoutsSection) GetContent() []string {
	return s.content
}

func (s *TimeoutsSection) Template() string {
	return `## Timeouts

The [bt]timeouts[bt] block allows you to specify [timeouts](https://www.terraform.io/language/resources/syntax#operation-timeouts) for certain actions:

{{ range .Timeouts -}}
{{ .String }}
{{ end -}}
`
}
