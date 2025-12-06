// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package markdown

import (
	"bufio"
	"fmt"
	"strings"

	"github.com/spf13/afero"
)

type Document struct {
	Path     string
	Sections []Section

	Exists    bool
	HasChange bool
}

func NewDocument(path string) *Document {
	return &Document{
		Sections: []Section{},
		Path:     path,
	}
}

func (d *Document) Write(fs afero.Fs) error {
	docContent := make([]string, 0)
	for _, section := range d.Sections {
		docContent = append(docContent, section.GetContent()...)
	}

	if err := afero.WriteFile(fs, d.Path, []byte(strings.Join(docContent, "\n")), 0o644); err != nil {
		return fmt.Errorf("writing to file `%s`: %+v", d.Path, err)
	}

	return nil
}

func (d *Document) Parse(fs afero.Fs, shouldNormalizeMd bool) error {
	var current Section
	var content []string

	file, err := fs.Open(d.Path)
	if err != nil {
		return fmt.Errorf("opening file `%s`: %+v", d.Path, err)
	}

	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		t := scanner.Text()

		// if length of `d.Content` is 0 and line has `---` prefix, we're parsing FrontMatter
		if strings.HasPrefix(t, "#") || (len(content) == 0 && strings.HasPrefix(t, "---")) {
			if current != nil {
				current.SetContent(content)
				d.Sections = append(d.Sections, current)

				current = nil
				content = nil
			}

			for _, s := range RegisteredSections() {
				if s.Match(t) {
					current = s
				}
			}

			// if we didn't match any, default to unknown section
			if current == nil {
				current = &SectionUnknown{}
			}

			current.SetHeading(t)
		}

		content = append(content, t)
	}

	// handle final section
	if current != nil {
		content = append(content, "")
		current.SetContent(content)
		d.Sections = append(d.Sections, current)
	}

	// Apply normalization to section if requested
	if shouldNormalizeMd {
		for _, section := range d.Sections {
			if normalizable, ok := section.(SectionWithNormalize); ok {
				normalizedContent, hasChange := normalizable.Normalize()
				if hasChange {
					section.SetContent(normalizedContent)
					d.HasChange = true
				}
			}
		}
	}

	return nil
}

// GetContent returns the full document content as a string
func (d *Document) GetContent() string {
	docContent := make([]string, 0)
	for _, section := range d.Sections {
		docContent = append(docContent, section.GetContent()...)
	}
	return strings.Join(docContent, "\n")
}

func (d *Document) GetArgumentsSection() *ArgumentsSection {
	for _, section := range d.Sections {
		if args, ok := section.(*ArgumentsSection); ok {
			return args
		}
	}
	return nil
}
