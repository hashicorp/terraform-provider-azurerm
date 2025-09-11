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
	// calculate a plan for an action. Depending on the type defined in
	// the action schema, Terraform may also pass the plan of linked resources
	// that the action can modify or return unmodified to influence Terraform's plan.
	PlanAction(context.Context, *PlanActionRequest) (*PlanActionResponse, error)

	// InvokeAction is called when Terraform wants to execute the logic of an action.
	// Depending on the type defined in the action schema, Terraform may also pass the
	// state of linked resources. The provider runs the logic of the action, reporting progress
	// events as desired, then sends a final complete event that has the linked resource's resulting
	// state and identity.
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

	// LinkedResources contains the data of the managed resource types that are linked to this action.
	//
	//   - If the action schema type is Unlinked, this field will be empty.
	//   - If the action schema type is Lifecycle, this field will be contain a single linked resource.
	//   - If the action schema type is Linked, this field will be one or more linked resources, which
	//     will be in the same order as the linked resource schemas are defined in the action schema.
	//
	// For Lifecycle actions, the provider may only change computed-only attributes.
	//
	// For Linked actions, the provider may change any attributes as long as it adheres to the resource schema.
	LinkedResources []*ProposedLinkedResource

	// Config is the configuration the user supplied for the action. See
	// the documentation on `DynamicValue` for more information about
	// safely accessing the configuration.
	Config *DynamicValue

	// ClientCapabilities defines optionally supported protocol features for the
	// PlanAction RPC, such as forward-compatible Terraform behavior changes.
	ClientCapabilities *PlanActionClientCapabilities
}

// ProposedLinkedResource represents linked resource data before PlanAction is called.
type ProposedLinkedResource struct {
	// PriorState is the state of the linked resource before the plan is applied,
	// represented as a `DynamicValue`. See the documentation for
	// `DynamicValue` for information about safely accessing the state.
	PriorState *DynamicValue

	// PlannedState is the latest indication of what the state for the
	// linked resource should be after apply, represented as a `DynamicValue`.
	// See the documentation for `DynamicValue` for information about safely
	// accessing the planned state.
	//
	// Since PlannedState is the most recent plan for the linked resource, it could
	// be the result of an RPC call to PlanResourceChange or an RPC call to PlanAction
	// for a predecessor action
	PlannedState *DynamicValue

	// Config is the configuration the user supplied for the linked resource. See
	// the documentation on `DynamicValue` for more information about
	// safely accessing the configuration.
	Config *DynamicValue

	// PriorIdentity is the identity of the resource before the plan is
	// applied, represented as a `ResourceIdentityData`.
	PriorIdentity *ResourceIdentityData
}

// PlanActionResponse is the response from the provider when planning an action. If the action
// has linked resources, it will contain any modifications made to the planned state or identity.
type PlanActionResponse struct {
	// LinkedResources contains the provider modified data of the managed resource types that are linked to this action.
	//
	// For Lifecycle actions, the provider may only change computed-only attributes.
	//
	// For Linked actions, the provider may change any attributes as long as it adheres to the resource schema.
	LinkedResources []*PlannedLinkedResource

	// Diagnostics report errors or warnings related to plannning the action and calculating
	// the planned state of the requested linked resources. Returning an empty slice
	// indicates a successful validation with no warnings or errors generated.
	Diagnostics []*Diagnostic

	// Deferred is used to indicate to Terraform that the PlanAction operation
	// needs to be deferred for a reason.
	Deferred *Deferred
}

// PlannedLinkedResource represents linked resource data that was planned during PlanAction and returned.
type PlannedLinkedResource struct {
	// PlannedState is the provider's indication of what the state for the
	// linked resource should be after apply, represented as a `DynamicValue`. See
	// the documentation for `DynamicValue` for information about safely
	// creating the `DynamicValue`.
	PlannedState *DynamicValue

	// PlannedIdentity is the provider's indication of what the identity for the
	// linked resource should be after apply, represented as a `ResourceIdentityData`
	PlannedIdentity *ResourceIdentityData
}

// InvokeActionRequest is the request Terraform sends when it wants to execute
// the logic of an action.
type InvokeActionRequest struct {
	// ActionType is the name of the action being called.
	ActionType string

	// LinkedResources contains the data of the managed resource types that are linked to this action.
	//
	//   - If the action schema type is Unlinked, this field will be empty.
	//   - If the action schema type is Lifecycle, this field will be contain a single linked resource.
	//   - If the action schema type is Linked, this field will be one or more linked resources, which
	//     will be in the same order as the linked resource schemas are defined in the action schema.
	//
	// For Lifecycle actions, the provider may only change computed-only attributes.
	//
	// For Linked actions, the provider may change any attributes as long as it adheres to the resource schema.
	LinkedResources []*InvokeLinkedResource

	// Config is the configuration the user supplied for the action. See
	// the documentation on `DynamicValue` for more information about
	// safely accessing the configuration.
	Config *DynamicValue
}

// InvokeLinkedResource represents linked resource data before InvokeAction is called.
type InvokeLinkedResource struct {
	// PriorState is the state of the linked resource before changes are applied,
	// represented as a `DynamicValue`. See the documentation for
	// `DynamicValue` for information about safely accessing the state.
	PriorState *DynamicValue

	// PlannedState is the latest indication of what the state for the
	// linked resource should look like after changes are applied, represented
	// as a `DynamicValue`. See the documentation for `DynamicValue` for
	// information about safely accessing the planned state.
	//
	// Since PlannedState is the most recent state for the linked resource, it could
	// be the result of an RPC call to ApplyResourceChange or an RPC call to InvokeAction
	// for a predecessor action.
	PlannedState *DynamicValue

	// Config is the configuration the user supplied for the linked resource. See
	// the documentation on `DynamicValue` for more information about
	// safely accessing the configuration.
	Config *DynamicValue

	// PlannedIdentity is Terraform's plan for what the linked resource identity should
	// look like after the changes are applied, represented as a `ResourceIdentityData`.
	PlannedIdentity *ResourceIdentityData
}

// InvokeActionServerStream represents a streaming response to an
// InvokeActionRequest.  An instance of this struct is supplied as an argument
// to the provider's InvokeAction implementation. The provider should set an
// Events iterator function that pushes zero or more events of type InvokeActionEvent.
type InvokeActionServerStream struct {
	// Events is the iterator that the provider can stream progress messages back to Terraform
	// as the action is executing. Once the provider has completed the action invocation, the provider must
	// respond with a completed event with the new linked resource state or diagnostics explaining why
	// the action failed.
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

// CompletedInvokeActionEventType represents the final completed event, along with all of the linked resource
// data modified by the provider or diagnostics about an action invocation failure.
type CompletedInvokeActionEventType struct {
	// LinkedResources contains the provider modified data of the managed resource types that are linked to this action.
	//
	// For Lifecycle actions, the provider may only change computed-only attributes.
	//
	// For Linked actions, the provider may change any attributes as long as it adheres to the resource schema.
	LinkedResources []*NewLinkedResource

	// Diagnostics report errors or warnings related to invoking an action.
	// Returning an empty slice indicates a successful invocation with no warnings
	// or errors generated.
	Diagnostics []*Diagnostic
}

func (a CompletedInvokeActionEventType) isInvokeActionEventType() {}

// NewLinkedResource represents linked resource data that was changed during InvokeAction and returned.
//
// Depending on how the action was invoked, the modified state data will either be immediately recorded in
// state or reconcicled in a future terraform apply operation.
type NewLinkedResource struct {
	// NewState is the provider's understanding of what the linked resource's
	// state is after changes are applied, represented as a `DynamicValue`.
	// See the documentation for `DynamicValue` for information about
	// safely creating the `DynamicValue`.
	//
	// Any attribute, whether computed or not, that has a known value in
	// the PlannedState in the InvokeActionRequest must be preserved
	// exactly as it was in NewState.
	NewState *DynamicValue

	// NewIdentity is the provider's understanding of what the linked resource's
	// identity is after changes are applied, represented as a `ResourceIdentityData`.
	NewIdentity *ResourceIdentityData

	// RequiresReplace can only be set if diagnostics are returned for the action and indicate
	// the linked resource must be replaced as a result of the action invocation error.
	RequiresReplace bool
}
