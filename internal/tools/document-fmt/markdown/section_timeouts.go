// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package markdown

import (
	"regexp"
	"strings"
)

type TimeoutsSection struct {
	heading Heading
	content []string
}

var _ SectionWithTemplate = &TimeoutsSection{}

func (s *TimeoutsSection) Match(line string) bool {
	return regexp.MustCompile(`#+(\s)*timeout.*`).MatchString(strings.ToLower(line))
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

The [bt]timeouts[bt] block allows you to specify [timeouts](https://developer.hashicorp.com/terraform/language/resources/configure#define-operation-timeouts) for certain actions:

{{ range .Timeouts -}}
{{ .String }}
{{ end -}}
`
}
