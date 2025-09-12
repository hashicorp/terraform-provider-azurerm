// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package schema

import (
	"context"

	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/internal/fwschema"
)

// TODO:Actions: Implement lifecycle and linked schemas
//
// SchemaType is the interface that an action schema type must implement. Action
// schema types are statically definined in the protocol, so all implementations
// are defined in this package.
//
// SchemaType implementations define how a practitioner can trigger an action, as well
// as what effect the action can have on the state. There are currently three different
// types of actions:
//   - [UnlinkedSchema] actions are actions that cannot cause changes to resource states.
//   - [LifecycleSchema] actions are actions that can cause changes to exactly one resource state.
//   - [LinkedSchema] actions are actions that can cause changes to one or more resource states.
type SchemaType interface {
	fwschema.Schema

	// MAINTAINER NOTE: Action schemas are unique to other schema types in framework in that the
	// exported methods all return a schema interface ([SchemaType]) rather than a schema struct,
	// due to the multiple different types of action schema implementations.
	//
	// As a result, there are certain methods that all schema structs implement that aren't defined in
	// the [fwschema.Schema] interface, such as the ValidateImplementation method. So we are adding that
	// here to the action schema interface to avoid additional internal interfaces and unnecessary
	// type assertions.
	ValidateImplementation(context.Context) diag.Diagnostics

	// Action schema types are statically defined in the protocol, so this
	// interface is not meant to be implemented outside of this package
	isActionSchemaType()
}
