// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package validator

import (
	"context"
)

// Describer is the common documentation interface for extensible schema
// validation functionality.
type Describer interface {
	// Description should describe the validation in plain text formatting.
	// This information is used by provider logging and provider tooling such
	// as documentation generation.
	//
	// The description should:
	//  - Begin with a lowercase or other character suitable for the middle of
	//    a sentence.
	//  - End without punctuation.
	//  - Use actionable language, such as "must" or "cannot".
	//  - Avoid newlines. Prefer separate validators instead.
	//
	// For example, "size must be less than 50 elements".
	Description(context.Context) string

	// MarkdownDescription should describe the validation in Markdown
	// formatting. This information is used by provider logging and provider
	// tooling such as documentation generation.
	//
	// The description should:
	//  - Begin with a lowercase or other character suitable for the middle of
	//    a sentence.
	//  - End without punctuation.
	//  - Use actionable language, such as "must" or "cannot".
	//  - Avoid newlines. Prefer separate validators instead.
	//
	// For example, "value must be `one` or `two`".
	MarkdownDescription(context.Context) string
}
