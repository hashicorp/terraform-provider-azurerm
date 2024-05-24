// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package provider

import (
	"context"

	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/function"
	"github.com/hashicorp/terraform-plugin-framework/resource"
)

// Provider is the core interface that all Terraform providers must implement.
//
// Providers can optionally implement these additional concepts:
//
//   - Validation: Schema-based or entire configuration
//     via ProviderWithConfigValidators or ProviderWithValidateConfig.
//   - Functions: ProviderWithFunctions
//   - Meta Schema: ProviderWithMetaSchema
type Provider interface {
	// Metadata should return the metadata for the provider, such as
	// a type name and version data.
	//
	// Implementing the MetadataResponse.TypeName will populate the
	// datasource.MetadataRequest.ProviderTypeName and
	// resource.MetadataRequest.ProviderTypeName fields automatically.
	Metadata(context.Context, MetadataRequest, *MetadataResponse)

	// Schema should return the schema for this provider.
	Schema(context.Context, SchemaRequest, *SchemaResponse)

	// Configure is called at the beginning of the provider lifecycle, when
	// Terraform sends to the provider the values the user specified in the
	// provider configuration block. These are supplied in the
	// ConfigureProviderRequest argument.
	// Values from provider configuration are often used to initialise an
	// API client, which should be stored on the struct implementing the
	// Provider interface.
	Configure(context.Context, ConfigureRequest, *ConfigureResponse)

	// DataSources returns a slice of functions to instantiate each DataSource
	// implementation.
	//
	// The data source type name is determined by the DataSource implementing
	// the Metadata method. All data sources must have unique names.
	DataSources(context.Context) []func() datasource.DataSource

	// Resources returns a slice of functions to instantiate each Resource
	// implementation.
	//
	// The resource type name is determined by the Resource implementing
	// the Metadata method. All resources must have unique names.
	Resources(context.Context) []func() resource.Resource
}

// ProviderWithConfigValidators is an interface type that extends Provider to include declarative validations.
//
// Declaring validation using this methodology simplifies implementation of
// reusable functionality. These also include descriptions, which can be used
// for automating documentation.
//
// Validation will include ConfigValidators and ValidateConfig, if both are
// implemented, in addition to any Attribute or Type validation.
type ProviderWithConfigValidators interface {
	Provider

	// ConfigValidators returns a list of functions which will all be performed during validation.
	ConfigValidators(context.Context) []ConfigValidator
}

// ProviderWithFunctions is an interface type that extends Provider to
// include provider defined functions for usage in practitioner configurations.
//
// Provider-defined functions are supported in Terraform version 1.8 and later.
type ProviderWithFunctions interface {
	Provider

	// Functions returns a slice of functions to instantiate each Function
	// implementation.
	//
	// The function name is determined by the Function implementing its Metadata
	// method. All functions must have unique names.
	Functions(context.Context) []func() function.Function
}

// ProviderWithMetaSchema is a provider with a provider meta schema, which
// is configured by practitioners via the provider_meta configuration block
// and the configuration data is included with certain data source and resource
// operations. The intended use case is to enable Terraform module authors
// within the same organization of the provider to track module usage in
// requests. Other use cases are explicitly not supported. All provider
// instances (aliases) receive the same data.
//
// This functionality is currently experimental and subject to change or break
// without warning. It is not protected by version compatibility guarantees.
type ProviderWithMetaSchema interface {
	Provider

	// MetaSchema should return the meta schema for this provider.
	//
	// This functionality is currently experimental and subject to change or
	// break without warning. It is not protected by version compatibility
	// guarantees.
	MetaSchema(context.Context, MetaSchemaRequest, *MetaSchemaResponse)
}

// ProviderWithValidateConfig is an interface type that extends Provider to include imperative validation.
//
// Declaring validation using this methodology simplifies one-off
// functionality that typically applies to a single provider. Any documentation
// of this functionality must be manually added into schema descriptions.
//
// Validation will include ConfigValidators and ValidateConfig, if both are
// implemented, in addition to any Attribute or Type validation.
type ProviderWithValidateConfig interface {
	Provider

	// ValidateConfig performs the validation.
	ValidateConfig(context.Context, ValidateConfigRequest, *ValidateConfigResponse)
}
