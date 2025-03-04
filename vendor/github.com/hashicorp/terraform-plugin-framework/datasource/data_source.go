// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package datasource

import (
	"context"
)

// DataSource represents an instance of a data source type. This is the core
// interface that all data sources must implement.
//
// Data sources can optionally implement these additional concepts:
//
//   - Configure: Include provider-level data or clients.
//   - Validation: Schema-based or entire configuration
//     via DataSourceWithConfigValidators or DataSourceWithValidateConfig.
type DataSource interface {
	// Metadata should return the full name of the data source, such as
	// examplecloud_thing.
	Metadata(context.Context, MetadataRequest, *MetadataResponse)

	// Schema should return the schema for this data source.
	Schema(context.Context, SchemaRequest, *SchemaResponse)

	// Read is called when the provider must read data source values in
	// order to update state. Config values should be read from the
	// ReadRequest and new state values set on the ReadResponse.
	Read(context.Context, ReadRequest, *ReadResponse)
}

// DataSourceWithConfigure is an interface type that extends DataSource to
// include a method which the framework will automatically call so provider
// developers have the opportunity to setup any necessary provider-level data
// or clients in the DataSource type.
//
// This method is intended to replace the provider.DataSourceType type
// NewDataSource method in a future release.
type DataSourceWithConfigure interface {
	DataSource

	// Configure enables provider-level data or clients to be set in the
	// provider-defined DataSource type. It is separately executed for each
	// ReadDataSource RPC.
	Configure(context.Context, ConfigureRequest, *ConfigureResponse)
}

// DataSourceWithConfigValidators is an interface type that extends DataSource to include declarative validations.
//
// Declaring validation using this methodology simplifies implmentation of
// reusable functionality. These also include descriptions, which can be used
// for automating documentation.
//
// Validation will include ConfigValidators and ValidateConfig, if both are
// implemented, in addition to any Attribute or Type validation.
type DataSourceWithConfigValidators interface {
	DataSource

	// ConfigValidators returns a list of ConfigValidators. Each ConfigValidator's Validate method will be called when validating the data source.
	ConfigValidators(context.Context) []ConfigValidator
}

// DataSourceWithValidateConfig is an interface type that extends DataSource to include imperative validation.
//
// Declaring validation using this methodology simplifies one-off
// functionality that typically applies to a single data source. Any
// documentation of this functionality must be manually added into schema
// descriptions.
//
// Validation will include ConfigValidators and ValidateConfig, if both are
// implemented, in addition to any Attribute or Type validation.
type DataSourceWithValidateConfig interface {
	DataSource

	// ValidateConfig performs the validation.
	ValidateConfig(context.Context, ValidateConfigRequest, *ValidateConfigResponse)
}
