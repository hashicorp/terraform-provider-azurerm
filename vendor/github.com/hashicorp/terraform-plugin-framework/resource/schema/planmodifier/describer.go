// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package planmodifier

import (
	"context"
)

// Describer is the common documentation interface for extensible schema
// plan modifier functionality.
type Describer interface {
	// Description should describe the plan modifier in plain text formatting.
	// This information is used by provider logging and provider tooling such
	// as documentation generation.
	//
	// The description should:
	//  - Begin with a lowercase or other character suitable for the middle of
	//    a sentence.
	//  - End without punctuation.
	Description(context.Context) string

	// MarkdownDescription should describe the plan modifier in Markdown
	// formatting. This information is used by provider logging and provider
	// tooling such as documentation generation.
	//
	// The description should:
	//  - Begin with a lowercase or other character suitable for the middle of
	//    a sentence.
	//  - End without punctuation.
	MarkdownDescription(context.Context) string
}
