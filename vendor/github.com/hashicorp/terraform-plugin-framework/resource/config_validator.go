// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package resource

import "context"

// ConfigValidator describes reusable Resource configuration validation functionality.
type ConfigValidator interface {
	// Description describes the validation in plain text formatting.
	//
	// This information may be automatically added to resource plain text
	// descriptions by external tooling.
	Description(context.Context) string

	// MarkdownDescription describes the validation in Markdown formatting.
	//
	// This information may be automatically added to resource Markdown
	// descriptions by external tooling.
	MarkdownDescription(context.Context) string

	// ValidateResource performs the validation.
	//
	// This method name is separate from the datasource.ConfigValidator
	// interface ValidateDataSource method name and provider.ConfigValidator
	// interface ValidateProvider method name to allow generic validators.
	ValidateResource(context.Context, ValidateConfigRequest, *ValidateConfigResponse)
}
