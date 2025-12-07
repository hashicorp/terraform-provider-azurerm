// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package markdown

import (
	"regexp"
	"strings"
)

type APISection struct {
	heading Heading
	content []string
}

var _ SectionWithTemplate = &APISection{}

func (s *APISection) Match(line string) bool {
	return regexp.MustCompile(`#+(\s)*api providers.*`).MatchString(strings.ToLower(line))
}

func (s *APISection) SetHeading(line string) {
	s.heading = NewHeading(line)
}

func (s *APISection) GetHeading() Heading {
	return s.heading
}

func (s *APISection) SetContent(content []string) {
	s.content = content
}

func (s *APISection) GetContent() []string {
	return s.content
}

func (s *APISection) Template() string {
	return `## API Providers
<!-- This section is generated, changes will be overwritten -->
This {{ .Type.String | lower }} uses the following Azure API Providers:
{{ range .APIs }}
{{ .String }}
{{ end -}}
`
}
