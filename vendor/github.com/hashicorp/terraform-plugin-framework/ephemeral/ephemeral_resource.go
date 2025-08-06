// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package ephemeral

import (
	"context"
)

// EphemeralResource represents an instance of an ephemeral resource type. This is the core
// interface that all ephemeral resources must implement.
//
// Ephemeral resources can optionally implement these additional concepts:
//
//   - Configure: Include provider-level data or clients via EphemeralResourceWithConfigure
//
//   - Validation: Schema-based or entire configuration via EphemeralResourceWithConfigValidators
//     or EphemeralResourceWithValidateConfig.
//
//   - Renew: Handle renewal of an expired remote object via EphemeralResourceWithRenew.
//     Ephemeral resources can indicate to Terraform when a renewal must occur via the RenewAt
//     response field of the Open/Renew methods. Renew cannot return new result data for the
//     ephemeral resource instance, so this logic is only appropriate for remote objects like
//     HashiCorp Vault leases, which can be renewed without changing their data.
//
//   - Close: Allows providers to clean up the ephemeral resource via EphemeralResourceWithClose.
type EphemeralResource interface {
	// Metadata should return the full name of the ephemeral resource, such as
	// examplecloud_thing.
	Metadata(context.Context, MetadataRequest, *MetadataResponse)

	// Schema should return the schema for this ephemeral resource.
	Schema(context.Context, SchemaRequest, *SchemaResponse)

	// Open is called when the provider must generate a new ephemeral resource. Config values
	// should be read from the OpenRequest and new response values set on the OpenResponse.
	Open(context.Context, OpenRequest, *OpenResponse)
}

// EphemeralResourceWithRenew is an interface type that extends EphemeralResource to
// include a method which the framework will call when Terraform detects that the
// provider-defined returned RenewAt time for an ephemeral resource has passed. This RenewAt
// response field can be set in the OpenResponse and RenewResponse.
type EphemeralResourceWithRenew interface {
	EphemeralResource

	// Renew is called when the provider must renew the ephemeral resource based on
	// the provided RenewAt time. This RenewAt response field can be set in the OpenResponse and RenewResponse.
	//
	// Renew cannot return new result data for the ephemeral resource instance, so this logic is only appropriate
	// for remote objects like HashiCorp Vault leases, which can be renewed without changing their data.
	Renew(context.Context, RenewRequest, *RenewResponse)
}

// EphemeralResourceWithClose is an interface type that extends
// EphemeralResource to include a method which the framework will call when
// Terraform determines that the ephemeral resource can be safely cleaned up.
type EphemeralResourceWithClose interface {
	EphemeralResource

	// Close is called when the provider can clean up the ephemeral resource.
	// Config values may be read from the CloseRequest.
	Close(context.Context, CloseRequest, *CloseResponse)
}

// EphemeralResourceWithConfigure is an interface type that extends EphemeralResource to
// include a method which the framework will automatically call so provider
// developers have the opportunity to setup any necessary provider-level data
// or clients in the EphemeralResource type.
type EphemeralResourceWithConfigure interface {
	EphemeralResource

	// Configure enables provider-level data or clients to be set in the
	// provider-defined EphemeralResource type.
	Configure(context.Context, ConfigureRequest, *ConfigureResponse)
}

// EphemeralResourceWithConfigValidators is an interface type that extends EphemeralResource to include declarative validations.
//
// Declaring validation using this methodology simplifies implementation of
// reusable functionality. These also include descriptions, which can be used
// for automating documentation.
//
// Validation will include ConfigValidators and ValidateConfig, if both are
// implemented, in addition to any Attribute or Type validation.
type EphemeralResourceWithConfigValidators interface {
	EphemeralResource

	// ConfigValidators returns a list of functions which will all be performed during validation.
	ConfigValidators(context.Context) []ConfigValidator
}

// EphemeralResourceWithValidateConfig is an interface type that extends EphemeralResource to include imperative validation.
//
// Declaring validation using this methodology simplifies one-off
// functionality that typically applies to a single ephemeral resource. Any documentation
// of this functionality must be manually added into schema descriptions.
//
// Validation will include ConfigValidators and ValidateConfig, if both are
// implemented, in addition to any Attribute or Type validation.
type EphemeralResourceWithValidateConfig interface {
	EphemeralResource

	// ValidateConfig performs the validation.
	ValidateConfig(context.Context, ValidateConfigRequest, *ValidateConfigResponse)
}
