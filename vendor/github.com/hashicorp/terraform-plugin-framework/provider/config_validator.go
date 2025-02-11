// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package provider

import "context"

// ConfigValidator describes reusable Provider configuration validation functionality.
type ConfigValidator interface {
	// Description describes the validation in plain text formatting.
	//
	// This information may be automatically added to provider plain text
	// descriptions by external tooling.
	Description(context.Context) string

	// MarkdownDescription describes the validation in Markdown formatting.
	//
	// This information may be automatically added to provider Markdown
	// descriptions by external tooling.
	MarkdownDescription(context.Context) string

	// ValidateProvider performs the validation.
	//
	// This method name is separate from the ConfigValidator
	// interface ValidateDataSource method name and ResourceConfigValidator
	// interface ValidateResource method name to allow generic validators.
	ValidateProvider(context.Context, ValidateConfigRequest, *ValidateConfigResponse)
}
