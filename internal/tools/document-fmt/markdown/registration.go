// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package markdown

func RegisteredSections() []Section {
	return []Section{
		&FrontMatterSection{},
		&ExampleSection{},
		&ArgumentsSection{},
		&AttributesSection{},
		&TimeoutsSection{},
		&APISection{},
	}
}
