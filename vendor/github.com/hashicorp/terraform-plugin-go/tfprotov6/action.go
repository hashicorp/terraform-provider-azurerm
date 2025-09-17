// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package tfprotov6

import (
	"context"
	"iter"
)

// ActionMetadata describes metadata for an action in the GetMetadata RPC.
type ActionMetadata struct {
	// TypeName is the name of the action.
	TypeName string
}

// ActionServer is an interface containing the methods an action implementation needs to fill.
type ActionServer interface {
	// ValidateActionConfig is called when Terraform is checking that an
	// action configuration is valid. It is guaranteed to have types
	// conforming to your schema, but it is not guaranteed that all values
	// will be known. This is your opportunity to do custom or advanced
	// validation prior to an action being planned/invoked.
	ValidateActionConfig(context.Context, *ValidateActionConfigRequest) (*ValidateActionConfigResponse, error)

	// PlanAction is called when Terraform is attempting to
	// calculate a plan for an action.
	PlanAction(context.Context, *PlanActionRequest) (*PlanActionResponse, error)

	// InvokeAction is called when Terraform wants to execute the logic of an action.
	// The provider runs the logic of the action, reporting progress
	// events as desired, then sends a final complete event.
	//
	// If an error occurs, the provider sends a complete event with the relevant diagnostics.
	InvokeAction(context.Context, *InvokeActionRequest) (*InvokeActionServerStream, error)
}

// ValidateActionConfigRequest is the request Terraform sends when it
// wants to validate an action's configuration.
type ValidateActionConfigRequest struct {
	// ActionType is the type of action Terraform is validating.
	ActionType string

	// Config is the configuration the user supplied for that action. See
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

// ValidateActionConfigResponse is the response from the provider about
// the validity of an action's configuration.
type ValidateActionConfigResponse struct {
	// Diagnostics report errors or warnings related to the given
	// configuration. Returning an empty slice indicates a successful
	// validation with no warnings or errors generated.
	Diagnostics []*Diagnostic
}

// PlanActionRequest is the request Terraform sends when it is attempting to
// calculate a plan for an action.
type PlanActionRequest struct {
	// ActionType is the name of the action being called.
	ActionType string

	// Config is the configuration the user supplied for the action. See
	// the documentation on `DynamicValue` for more information about
	// safely accessing the configuration.
	Config *DynamicValue

	// ClientCapabilities defines optionally supported protocol features for the
	// PlanAction RPC, such as forward-compatible Terraform behavior changes.
	ClientCapabilities *PlanActionClientCapabilities
}

// PlanActionResponse is the response from the provider when planning an action.
type PlanActionResponse struct {
	// Diagnostics report errors or warnings related to planning the action. Returning an empty slice
	// indicates a successful validation with no warnings or errors generated.
	Diagnostics []*Diagnostic

	// Deferred is used to indicate to Terraform that the PlanAction operation
	// needs to be deferred for a reason.
	Deferred *Deferred
}

// InvokeActionRequest is the request Terraform sends when it wants to execute
// the logic of an action.
type InvokeActionRequest struct {
	// ActionType is the name of the action being called.
	ActionType string

	// Config is the configuration the user supplied for the action. See
	// the documentation on `DynamicValue` for more information about
	// safely accessing the configuration.
	Config *DynamicValue

	// ClientCapabilities defines optionally supported protocol features for the
	// InvokeAction RPC, such as forward-compatible Terraform behavior changes.
	ClientCapabilities *InvokeActionClientCapabilities
}

// InvokeActionServerStream represents a streaming response to an
// InvokeActionRequest.  An instance of this struct is supplied as an argument
// to the provider's InvokeAction implementation. The provider should set an
// Events iterator function that pushes zero or more events of type InvokeActionEvent.
type InvokeActionServerStream struct {
	// Events is the iterator that the provider can stream progress messages back to Terraform
	// as the action is executing. Once the provider has completed the action invocation, the provider must
	// respond with a completed event. If the action failed, the completed event must contain
	// diagnostics explaining why the action failed.
	Events iter.Seq[InvokeActionEvent]
}

// InvokeActionEvent is an event sent back to Terraform during the InvokeAction RPC.
type InvokeActionEvent struct {
	// Type is the type of event that is being sent back during InvokeAction, either a Progress event
	// or a Completed event.
	Type InvokeActionEventType
}

// InvokeActionEventType is an intentionally unimplementable interface that
// functions as an enum, allowing us to use different strongly-typed event types
// that contain additional, but different data, as a generic "event" type.
type InvokeActionEventType interface {
	isInvokeActionEventType() // this interface is only implementable in this package
}

var (
	_ InvokeActionEventType = ProgressInvokeActionEventType{}
	_ InvokeActionEventType = CompletedInvokeActionEventType{}
)

// ProgressInvokeActionEventType represents a progress update that should be displayed in the Terraform
// CLI or external system running Terraform.
type ProgressInvokeActionEventType struct {
	// Message is the human-readable message to display about the progress of the action invocation.
	Message string
}

func (a ProgressInvokeActionEventType) isInvokeActionEventType() {}

// CompletedInvokeActionEventType represents the final completed event, along with
// potential diagnostics about an action invocation failure.
type CompletedInvokeActionEventType struct {
	// Diagnostics report errors or warnings related to invoking an action.
	// Returning an empty slice indicates a successful invocation with no warnings
	// or errors generated.
	Diagnostics []*Diagnostic
}

func (a CompletedInvokeActionEventType) isInvokeActionEventType() {}
