package markdown

import (
	"strings"
)

// TODO: Should this be a section, or a separate "FrontMatter" field in Document struct?
type FrontMatterSection struct {
	heading Heading
	content []string
}

var _ SectionWithTemplate = &FrontMatterSection{}

func (s *FrontMatterSection) Match(line string) bool {
	return strings.HasPrefix(line, "---")
}

func (s *FrontMatterSection) SetContent(content []string) {
	s.content = content
}

func (s *FrontMatterSection) GetContent() []string {
	return s.content
}

func (s *FrontMatterSection) SetHeading(line string) {
	s.heading = Heading{
		Level: 0,
		Text:  line,
	}
}

func (s *FrontMatterSection) GetHeading() Heading {
	return s.heading
}

func (s *FrontMatterSection) Template() string {
	return `
---
subcategory: "TODO"
layout: "{{ .ProviderName }}"
page_title: "Azure Resource Manager: {{ .Name }}"
description: |-
  TODO
---
`
}
