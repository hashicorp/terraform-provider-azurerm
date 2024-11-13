// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package defaults

import "context"

// Describer is the common documentation interface for extensible schema
// default value functionality.
type Describer interface {
	// Description should describe the default in plain text formatting.
	// This information is used by provider logging and provider tooling such
	// as documentation generation.
	//
	// The description should:
	//  - Begin with a lowercase or other character suitable for the middle of
	//    a sentence.
	//  - End without punctuation.
	Description(ctx context.Context) string

	// MarkdownDescription should describe the default in Markdown
	// formatting. This information is used by provider logging and provider
	// tooling such as documentation generation.
	//
	// The description should:
	//  - Begin with a lowercase or other character suitable for the middle of
	//    a sentence.
	//  - End without punctuation.
	MarkdownDescription(ctx context.Context) string
}
