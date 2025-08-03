// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package ephemeral

import "context"

// ConfigValidator describes reusable EphemeralResource configuration validation functionality.
type ConfigValidator interface {
	// Description describes the validation in plain text formatting.
	//
	// This information may be automatically added to ephemeral resource plain text
	// descriptions by external tooling.
	Description(context.Context) string

	// MarkdownDescription describes the validation in Markdown formatting.
	//
	// This information may be automatically added to ephemeral resource Markdown
	// descriptions by external tooling.
	MarkdownDescription(context.Context) string

	// ValidateEphemeralResource performs the validation.
	//
	// This method name is separate from the datasource.ConfigValidator
	// interface ValidateDataSource method name, provider.ConfigValidator
	// interface ValidateProvider method name, and resource.ConfigValidator
	// interface ValidateResource method name to allow generic validators.
	ValidateEphemeralResource(context.Context, ValidateConfigRequest, *ValidateConfigResponse)
}
