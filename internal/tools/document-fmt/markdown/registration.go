// Copyright (c) HashiCorp, Inc.
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
