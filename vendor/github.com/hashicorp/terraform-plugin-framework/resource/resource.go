// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package resource

import (
	"context"
)

// Resource represents an instance of a managed resource type. This is the core
// interface that all resources must implement.
//
// Resources can optionally implement these additional concepts:
//
//   - Configure: Include provider-level data or clients.
//   - Import: ResourceWithImportState
//   - Validation: Schema-based or entire configuration
//     via ResourceWithConfigValidators or ResourceWithValidateConfig.
//   - Plan Modification: Schema-based or entire plan
//     via ResourceWithModifyPlan.
//   - State Upgrades: ResourceWithUpgradeState
//
// Although not required, it is conventional for resources to implement the
// ResourceWithImportState interface.
type Resource interface {
	// Metadata should return the full name of the resource, such as
	// examplecloud_thing.
	Metadata(context.Context, MetadataRequest, *MetadataResponse)

	// Schema should return the schema for this resource.
	Schema(context.Context, SchemaRequest, *SchemaResponse)

	// Create is called when the provider must create a new resource. Config
	// and planned state values should be read from the
	// CreateRequest and new state values set on the CreateResponse.
	Create(context.Context, CreateRequest, *CreateResponse)

	// Read is called when the provider must read resource values in order
	// to update state. Planned state values should be read from the
	// ReadRequest and new state values set on the ReadResponse.
	Read(context.Context, ReadRequest, *ReadResponse)

	// Update is called to update the state of the resource. Config, planned
	// state, and prior state values should be read from the
	// UpdateRequest and new state values set on the UpdateResponse.
	Update(context.Context, UpdateRequest, *UpdateResponse)

	// Delete is called when the provider must delete the resource. Config
	// values may be read from the DeleteRequest.
	//
	// If execution completes without error, the framework will automatically
	// call DeleteResponse.State.RemoveResource(), so it can be omitted
	// from provider logic.
	Delete(context.Context, DeleteRequest, *DeleteResponse)
}

// ResourceWithConfigure is an interface type that extends Resource to
// include a method which the framework will automatically call so provider
// developers have the opportunity to setup any necessary provider-level data
// or clients in the Resource type.
//
// This method is intended to replace the provider.ResourceType type
// NewResource method in a future release.
type ResourceWithConfigure interface {
	Resource

	// Configure enables provider-level data or clients to be set in the
	// provider-defined Resource type. It is separately executed for each
	// ReadResource RPC.
	Configure(context.Context, ConfigureRequest, *ConfigureResponse)
}

// ResourceWithConfigValidators is an interface type that extends Resource to include declarative validations.
//
// Declaring validation using this methodology simplifies implmentation of
// reusable functionality. These also include descriptions, which can be used
// for automating documentation.
//
// Validation will include ConfigValidators and ValidateConfig, if both are
// implemented, in addition to any Attribute or Type validation.
type ResourceWithConfigValidators interface {
	Resource

	// ConfigValidators returns a list of functions which will all be performed during validation.
	ConfigValidators(context.Context) []ConfigValidator
}

// Optional interface on top of Resource that enables provider control over
// the ImportResourceState RPC. This RPC is called by Terraform when the
// `terraform import` command is executed. Afterwards, the ReadResource RPC
// is executed to allow providers to fully populate the resource state.
type ResourceWithImportState interface {
	Resource

	// ImportState is called when the provider must import the state of a
	// resource instance. This method must return enough state so the Read
	// method can properly refresh the full resource.
	//
	// If setting an attribute with the import identifier, it is recommended
	// to use the ImportStatePassthroughID() call in this method.
	ImportState(context.Context, ImportStateRequest, *ImportStateResponse)
}

// ResourceWithModifyPlan represents a resource instance with a ModifyPlan
// function.
type ResourceWithModifyPlan interface {
	Resource

	// ModifyPlan is called when the provider has an opportunity to modify
	// the plan: once during the plan phase when Terraform is determining
	// the diff that should be shown to the user for approval, and once
	// during the apply phase with any unknown values from configuration
	// filled in with their final values.
	//
	// The planned new state is represented by
	// ModifyPlanResponse.Plan. It must meet the following
	// constraints:
	// 1. Any non-Computed attribute set in config must preserve the exact
	// config value or return the corresponding attribute value from the
	// prior state (ModifyPlanRequest.State).
	// 2. Any attribute with a known value must not have its value changed
	// in subsequent calls to ModifyPlan or Create/Read/Update.
	// 3. Any attribute with an unknown value may either remain unknown
	// or take on any value of the expected type.
	//
	// Any errors will prevent further resource-level plan modifications.
	ModifyPlan(context.Context, ModifyPlanRequest, *ModifyPlanResponse)
}

// Optional interface on top of [Resource] that enables provider control over
// the MoveResourceState RPC. This RPC is called by Terraform when there is a
// `moved` configuration block that changes the resource type and where this
// [Resource] is the target resource type. Since state data operations can cause
// data loss for practitioners, this support is explicitly opt-in to ensure that
// all data transformation logic is explicitly defined by the provider.
//
// If the [Resource] does not implement this interface and Terraform sends a
// MoveResourceState request, the framework will automatically return an error
// diagnostic notifying the practitioner that this resource does not support the
// requested operation.
//
// This functionality is only supported in Terraform 1.8 and later.
type ResourceWithMoveState interface {
	Resource

	// An ordered list of source resource to current schema version state move
	// implementations. Only the first [StateMover] implementation that returns
	// state data or error diagnostics will be used, otherwise the framework
	// considers the [StateMover] as skipped and will try the next [StateMover].
	// If all implementations return without state and error diagnostics, the
	// framework will return an implementation not found error.
	//
	// It is strongly recommended that implementations be overly cautious and
	// return no state data if the source provider address, resource type,
	// or schema version is not fully implemented.
	MoveState(context.Context) []StateMover
}

// Optional interface on top of Resource that enables provider control over
// the UpgradeResourceState RPC. This RPC is automatically called by Terraform
// when the current Schema type Version field is greater than the stored state.
// Terraform does not store previous Schema information, so any breaking
// changes to state data types must be handled by providers.
//
// Terraform CLI can execute the UpgradeResourceState RPC even when the prior
// state version matches the current schema version. The framework will
// automatically intercept this request and attempt to respond with the
// existing state. In this situation the framework will not execute any
// provider defined logic, so declaring it for this version is extraneous.
type ResourceWithUpgradeState interface {
	Resource

	// A mapping of prior state version to current schema version state upgrade
	// implementations. Only the specified state upgrader for the prior state
	// version is called, rather than each version in between, so it must
	// encapsulate all logic to convert the prior state to the current schema
	// version.
	//
	// Version keys begin at 0, which is the default schema version when
	// undefined. The framework will return an error diagnostic should the
	// requested state version not be implemented.
	UpgradeState(context.Context) map[int64]StateUpgrader
}

// ResourceWithValidateConfig is an interface type that extends Resource to include imperative validation.
//
// Declaring validation using this methodology simplifies one-off
// functionality that typically applies to a single resource. Any documentation
// of this functionality must be manually added into schema descriptions.
//
// Validation will include ConfigValidators and ValidateConfig, if both are
// implemented, in addition to any Attribute or Type validation.
type ResourceWithValidateConfig interface {
	Resource

	// ValidateConfig performs the validation.
	ValidateConfig(context.Context, ValidateConfigRequest, *ValidateConfigResponse)
}
