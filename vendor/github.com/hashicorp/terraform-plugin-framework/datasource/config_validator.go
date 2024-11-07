// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package datasource

import "context"

// ConfigValidator describes reusable data source configuration validation functionality.
type ConfigValidator interface {
	// Description describes the validation in plain text formatting.
	//
	// This information may be automatically added to data source plain text
	// descriptions by external tooling.
	Description(context.Context) string

	// MarkdownDescription describes the validation in Markdown formatting.
	//
	// This information may be automatically added to data source Markdown
	// descriptions by external tooling.
	MarkdownDescription(context.Context) string

	// ValidateDataSource performs the validation.
	//
	// This method name is separate from the provider.ConfigValidator
	// interface ValidateProvider method name and resource.ConfigValidator
	// interface ValidateResource method name to allow generic validators.
	ValidateDataSource(context.Context, ValidateConfigRequest, *ValidateConfigResponse)
}
