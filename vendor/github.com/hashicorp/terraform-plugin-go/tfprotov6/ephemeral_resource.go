// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package tfprotov6

import (
	"context"
	"time"
)

// EphemeralResourceMetadata describes metadata for an ephemeral resource in the GetMetadata
// RPC.
type EphemeralResourceMetadata struct {
	// TypeName is the name of the ephemeral resource.
	TypeName string
}

// EphemeralResourceServer is an interface containing the methods an ephemeral resource
// implementation needs to fill.
type EphemeralResourceServer interface {
	// ValidateEphemeralResourceConfig is called when Terraform is checking that an
	// ephemeral resource configuration is valid. It is guaranteed to have types
	// conforming to your schema, but it is not guaranteed that all values
	// will be known. This is your opportunity to do custom or advanced
	// validation prior to an ephemeral resource being opened.
	ValidateEphemeralResourceConfig(context.Context, *ValidateEphemeralResourceConfigRequest) (*ValidateEphemeralResourceConfigResponse, error)

	// OpenEphemeralResource is called when Terraform wants to open the ephemeral resource,
	// usually during planning. If the config for the ephemeral resource contains unknown
	// values, Terraform will defer the OpenEphemeralResource call until apply.
	OpenEphemeralResource(context.Context, *OpenEphemeralResourceRequest) (*OpenEphemeralResourceResponse, error)

	// RenewEphemeralResource is called when Terraform detects that the previously specified
	// RenewAt timestamp has passed. The RenewAt timestamp is supplied either from the
	// OpenEphemeralResource call or a previous RenewEphemeralResource call.
	RenewEphemeralResource(context.Context, *RenewEphemeralResourceRequest) (*RenewEphemeralResourceResponse, error)

	// CloseEphemeralResource is called when Terraform is closing the ephemeral resource.
	CloseEphemeralResource(context.Context, *CloseEphemeralResourceRequest) (*CloseEphemeralResourceResponse, error)
}

// ValidateEphemeralResourceConfigRequest is the request Terraform sends when it
// wants to validate an ephemeral resource's configuration.
type ValidateEphemeralResourceConfigRequest struct {
	// TypeName is the type of resource Terraform is validating.
	TypeName string

	// Config is the configuration the user supplied for that ephemeral resource. See
	// the documentation on `DynamicValue` for more information about
	// safely accessing the configuration.
	//
	// The configuration is represented as a tftypes.Object, with each
	// attribute and nested block getting its own key and value.
	//
	// This configuration may contain unknown values if a user uses
	// interpolation or other functionality that would prevent Terraform
	// from knowing the value at request time. Any attributes not directly
	// set in the configuration will be null.
	Config *DynamicValue
}

// ValidateEphemeralResourceConfigResponse is the response from the provider about
// the validity of an ephemeral resource's configuration.
type ValidateEphemeralResourceConfigResponse struct {
	// Diagnostics report errors or warnings related to the given
	// configuration. Returning an empty slice indicates a successful
	// validation with no warnings or errors generated.
	Diagnostics []*Diagnostic
}

// OpenEphemeralResourceRequest is the request Terraform sends when it
// wants to open an ephemeral resource.
type OpenEphemeralResourceRequest struct {
	// TypeName is the type of resource Terraform is opening.
	TypeName string

	// Config is the configuration the user supplied for that ephemeral resource. See
	// the documentation on `DynamicValue` for more information about
	// safely accessing the configuration.
	//
	// The configuration is represented as a tftypes.Object, with each
	// attribute and nested block getting its own key and value.
	//
	// This configuration will always be fully known. If Config contains unknown values,
	// Terraform will defer the OpenEphemeralResource RPC until apply.
	Config *DynamicValue

	// ClientCapabilities defines optionally supported protocol features for the
	// OpenEphemeralResource RPC, such as forward-compatible Terraform behavior changes.
	ClientCapabilities *OpenEphemeralResourceClientCapabilities
}

// OpenEphemeralResourceResponse is the response from the provider about the current
// state of the opened ephemeral resource.
type OpenEphemeralResourceResponse struct {
	// Result is the provider's understanding of what the ephemeral resource's
	// data is after it has been opened, represented as a `DynamicValue`.
	// See the documentation for `DynamicValue` for information about
	// safely creating the `DynamicValue`.
	//
	// Any attribute, whether computed or not, that has a known value in
	// the Config in the OpenEphemeralResourceRequest must be preserved
	// exactly as it was in Result.
	//
	// Any attribute in the Config in the OpenEphemeralResourceRequest
	// that is unknown must take on a known value at this time. No unknown
	// values are allowed in the Result.
	//
	// The result should be represented as a tftypes.Object, with each
	// attribute and nested block getting its own key and value.
	Result *DynamicValue

	// Diagnostics report errors or warnings related to opening the
	// requested ephemeral resource. Returning an empty slice
	// indicates a successful creation with no warnings or errors
	// generated.
	Diagnostics []*Diagnostic

	// Private should be set to any private data that the provider would like to be
	// sent to the next Renew or Close call.
	Private []byte

	// RenewAt indicates to Terraform that the ephemeral resource
	// needs to be renewed at the specified time. Terraform will
	// call the RenewEphemeralResource RPC when the specified time has passed.
	RenewAt time.Time

	// Deferred is used to indicate to Terraform that the OpenEphemeralResource operation
	// needs to be deferred for a reason.
	Deferred *Deferred
}

// RenewEphemeralResourceRequest is the request Terraform sends when it
// wants to renew an ephemeral resource.
type RenewEphemeralResourceRequest struct {
	// TypeName is the type of resource Terraform is renewing.
	TypeName string

	// Private is any provider-defined private data stored with the
	// ephemeral resource from the most recent Open or Renew call.
	//
	// To ensure private data is preserved, copy any necessary data to
	// the RenewEphemeralResourceResponse type Private field.
	Private []byte
}

// RenewEphemeralResourceResponse is the response from the provider after an ephemeral resource
// has been renewed.
type RenewEphemeralResourceResponse struct {
	// Diagnostics report errors or warnings related to renewing the
	// requested ephemeral resource. Returning an empty slice
	// indicates a successful creation with no warnings or errors
	// generated.
	Diagnostics []*Diagnostic

	// Private should be set to any private data that the provider would like to be
	// sent to the next Renew or Close call.
	Private []byte

	// RenewAt indicates to Terraform that the ephemeral resource
	// needs to be renewed at the specified time. Terraform will
	// call the RenewEphemeralResource RPC when the specified time has passed.
	RenewAt time.Time
}

// CloseEphemeralResourceRequest is the request Terraform sends when it
// wants to close an ephemeral resource.
type CloseEphemeralResourceRequest struct {
	// TypeName is the type of resource Terraform is closing.
	TypeName string

	// Private is any provider-defined private data stored with the
	// ephemeral resource from the most recent Open or Renew call.
	Private []byte
}

// CloseEphemeralResourceResponse is the response from the provider about
// the closed ephemeral resource.
type CloseEphemeralResourceResponse struct {
	// Diagnostics report errors or warnings related to closing the
	// requested ephemeral resource. Returning an empty slice
	// indicates a successful creation with no warnings or errors
	// generated.
	Diagnostics []*Diagnostic
}
