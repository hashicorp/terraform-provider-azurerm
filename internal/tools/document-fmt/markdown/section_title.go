// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package markdown

import "regexp"

type TitleSection struct {
	heading Heading
	content []string
}

var _ SectionWithTemplate = &TitleSection{}

func (s *TitleSection) Match(line string) bool {
	match, _ := regexp.MatchString(`#+[\s\t]*([\w\s\t]*:\s*)?\w*_+[\w_]*`, line)
	return match
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
	return `# {{ if eq .Type.String "Data Source" }}Data Source: {{ end }}{{ .Name }}
{{ if eq .Type.String "Data Source" }}
Use this data source to access information about an existing <brandName>.
{{- else }}
Manages a <brandName>.
{{- end }}
`
}
