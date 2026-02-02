// Copyright IBM Corp. 2014, 2025
// SPDX-License-Identifier: MPL-2.0

package markdown

func RegisteredSections() []Section {
	return []Section{
		&FrontMatterSection{},
		&TitleSection{},
		&ExampleSection{},
		&ArgumentsSection{},
		&AttributesSection{},
		&TimeoutsSection{},
		&APISection{},
	}
}
