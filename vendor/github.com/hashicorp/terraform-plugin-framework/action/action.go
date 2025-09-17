// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package action

import "context"

type Action interface {
	// Schema should return the schema for this action.
	Schema(context.Context, SchemaRequest, *SchemaResponse)

	// Metadata should return the full name of the action, such as examplecloud_do_thing.
	Metadata(context.Context, MetadataRequest, *MetadataResponse)

	// Invoke is called to run the logic of the action. Config values should be read from the InvokeRequest
	// and potential diagnostics set in InvokeResponse.
	//
	// The [InvokeResponse.SendProgress] function can be called in the Invoke method to immediately
	// report progress events related to the invocation of the action to Terraform.
	Invoke(context.Context, InvokeRequest, *InvokeResponse)
}

// ActionWithConfigure is an interface type that extends Action to
// include a method which the framework will automatically call so provider
// developers have the opportunity to setup any necessary provider-level data
// or clients in the Action type.
type ActionWithConfigure interface {
	Action

	// Configure enables provider-level data or clients to be set in the
	// provider-defined Action type.
	Configure(context.Context, ConfigureRequest, *ConfigureResponse)
}

// ActionWithModifyPlan represents an action with a ModifyPlan function.
type ActionWithModifyPlan interface {
	Action

	// ModifyPlan is called when the provider has an opportunity to modify
	// the plan for an action: once during the plan phase, and once
	// during the apply phase with any unknown values from configuration
	// filled in with their final values.
	//
	// All action schema types can use the plan as an opportunity to raise early
	// diagnostics to practitioners, such as validation errors.
	ModifyPlan(context.Context, ModifyPlanRequest, *ModifyPlanResponse)
}

// ActionWithConfigValidators is an interface type that extends Action to include declarative validations.
//
// Declaring validation using this methodology simplifies implementation of
// reusable functionality. These also include descriptions, which can be used
// for automating documentation.
//
// Validation will include ConfigValidators and ValidateConfig, if both are
// implemented, in addition to any Attribute or Type validation.
type ActionWithConfigValidators interface {
	Action

	// ConfigValidators returns a list of functions which will all be performed during validation.
	ConfigValidators(context.Context) []ConfigValidator
}

// ActionWithValidateConfig is an interface type that extends Action to include imperative validation.
//
// Declaring validation using this methodology simplifies one-off
// functionality that typically applies to a single action. Any documentation
// of this functionality must be manually added into schema descriptions.
//
// Validation will include ConfigValidators and ValidateConfig, if both are
// implemented, in addition to any Attribute or Type validation.
type ActionWithValidateConfig interface {
	Action

	// ValidateConfig performs the validation.
	ValidateConfig(context.Context, ValidateConfigRequest, *ValidateConfigResponse)
}
