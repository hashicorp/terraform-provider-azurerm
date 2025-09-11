// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package tfprotov6

// ActionSchema is how Terraform defines the shape of action data and
// how the practitioner can interact with the action.
type ActionSchema struct {
	// Schema is the definition for the action data itself, which will be specified in an action block in the user's configuration.
	Schema *Schema

	// Type defines how a practitioner can trigger an action, as well as what effect the action can have on the state
	// of the linked managed resources. There are currently three different types of actions:
	//   - Unlinked actions are actions that cannot cause changes to resource states.
	//   - Lifecycle actions are actions that can cause changes to exactly one resource state.
	//   - Linked actions are actions that can cause changes to one or more resource states.
	Type ActionSchemaType
}

// ActionSchemaType is an intentionally unimplementable interface that
// functions as an enum, allowing us to use different strongly-typed action schema types
// that contain additional, but different data, as a generic "action" type.
//
// An action can only be one type (Unlinked, Lifecycle, or Linked), which are all statically defined in the protocol.
type ActionSchemaType interface {
	isActionSchemaType()
}

var (
	_ ActionSchemaType = UnlinkedActionSchemaType{}
	_ ActionSchemaType = LifecycleActionSchemaType{}
	_ ActionSchemaType = LinkedActionSchemaType{}
)

// UnlinkedActionSchemaType represents an unlinked action, which cannot cause changes to resource states.
type UnlinkedActionSchemaType struct{}

func (a UnlinkedActionSchemaType) isActionSchemaType() {}

// LifecycleActionSchemaType represents a lifecycle action, which can cause changes to exactly one resource state,
// which is the linked resource.
type LifecycleActionSchemaType struct {
	// Executes defines when the lifecycle action must be executed in relation to the linked resource, either before
	// or after the linked resource's plan/apply.
	Executes LifecycleExecutionOrder

	// LinkedResource is the managed resource type that this action can make state changes to.
	// This linked resource is currently restricted to be defined in the same provider as the action is defined.
	LinkedResource *LinkedResourceSchema
}

func (a LifecycleActionSchemaType) isActionSchemaType() {}

const (
	// LifecycleExecutionOrderInvalid is used to indicate an invalid `LifecycleExecutionOrder`.
	// Provider developers should not use it.
	LifecycleExecutionOrderInvalid LifecycleExecutionOrder = 0

	// LifecycleExecutionOrderBefore is used to indicate that the action must be invoked before it's
	// linked resource's plan/apply.
	LifecycleExecutionOrderBefore LifecycleExecutionOrder = 1

	// LifecycleExecutionOrderAfter is used to indicate that the action must be invoked after it's
	// linked resource's plan/apply.
	LifecycleExecutionOrderAfter LifecycleExecutionOrder = 2
)

// LifecycleExecutionOrder is an enum that represents when an action is invoked relative to it's linked resource.
type LifecycleExecutionOrder int32

func (l LifecycleExecutionOrder) String() string {
	switch l {
	case 0:
		return "INVALID"
	case 1:
		return "BEFORE"
	case 2:
		return "AFTER"
	case 3:
	}
	return "UNKNOWN"
}

// LinkedResourceSchema represents information about the schema of a linked resource, which is used by an action schema to describe to Terraform the
// resource types that an action is allowed to change the state of. Linked resources are currently restricted to be defined in the same provider
// as the action is defined.
//
// LinkedResourceSchema does not contain the entire schema definition of the linked resource, which must be obtained by the provider in order to
// decode the linked resource plan/state/identity protocol data during PlanAction and InvokeAction.
type LinkedResourceSchema struct {
	// TypeName is the name of the managed resource which can have it's resource state changed by the action. The name should be prefixed with
	// the provider shortname and an underscore.
	TypeName string

	// Description is a human-readable description of the linked resource.
	Description string
}

// LinkedActionSchemaType represents a linked action, which can cause changes to one or more resource states.
type LinkedActionSchemaType struct {
	// LinkedResources are the managed resource types that this action can make state changes to.
	// These linked resources are currently restricted to be defined in the same provider as the action is defined.
	LinkedResources []*LinkedResourceSchema
}

func (a LinkedActionSchemaType) isActionSchemaType() {}
