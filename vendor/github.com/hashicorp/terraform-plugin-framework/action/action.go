// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package action

import "context"

type Action interface {
	// Schema should return the schema for this action.
	Schema(context.Context, SchemaRequest, *SchemaResponse)

	// Metadata should return the full name of the action, such as examplecloud_do_thing.
	Metadata(context.Context, MetadataRequest, *MetadataResponse)

	// Invoke is called to run the logic of the action and update linked resources if applicable.
	// Config, linked resource planned state, and linked resource prior state values should
	// be read from the InvokeRequest and new linked resource state values set on the InvokeResponse.
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
	// Actions do not have computed attributes that can be modified during the plan,
	// but linked and lifecycle actions can modify the plan of linked resources.
	//
	// All action schema types can use the plan as an opportunity to raise early
	// diagnostics to practitioners, such as validation errors.
	ModifyPlan(context.Context, ModifyPlanRequest, *ModifyPlanResponse)
}
