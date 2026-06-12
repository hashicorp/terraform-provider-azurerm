// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package list

import "context"

// ConfigValidator describes reusable ListResource configuration validation
// functionality.
type ConfigValidator interface {
	// Description describes the validation in plain text formatting.
	//
	// This information may be automatically added to list resource plain text
	// descriptions by external tooling.
	Description(context.Context) string

	// MarkdownDescription describes the validation in Markdown formatting.
	//
	// This information may be automatically added to list resource Markdown
	// descriptions by external tooling.
	MarkdownDescription(context.Context) string

	// ValidateResource performs the validation.
	//
	// This method name is separate from ConfigValidators in resource and other packages in
	// order to allow generic validators.
	ValidateListResourceConfig(context.Context, ValidateConfigRequest, *ValidateConfigResponse)
}
