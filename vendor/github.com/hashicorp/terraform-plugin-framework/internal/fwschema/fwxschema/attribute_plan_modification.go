// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package fwxschema

import (
	"github.com/hashicorp/terraform-plugin-framework/internal/fwschema"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/planmodifier"
)

// AttributeWithBoolPlanModifiers is an optional interface on Attribute which
// enables Bool plan modifier support.
type AttributeWithBoolPlanModifiers interface {
	fwschema.Attribute

	// BoolPlanModifiers should return a list of Bool plan modifiers.
	BoolPlanModifiers() []planmodifier.Bool
}

// AttributeWithFloat64PlanModifiers is an optional interface on Attribute which
// enables Float64 plan modifier support.
type AttributeWithFloat64PlanModifiers interface {
	fwschema.Attribute

	// Float64PlanModifiers should return a list of Float64 plan modifiers.
	Float64PlanModifiers() []planmodifier.Float64
}

// AttributeWithInt64PlanModifiers is an optional interface on Attribute which
// enables Int64 plan modifier support.
type AttributeWithInt64PlanModifiers interface {
	fwschema.Attribute

	// Int64PlanModifiers should return a list of Int64 plan modifiers.
	Int64PlanModifiers() []planmodifier.Int64
}

// AttributeWithListPlanModifiers is an optional interface on Attribute which
// enables List plan modifier support.
type AttributeWithListPlanModifiers interface {
	fwschema.Attribute

	// ListPlanModifiers should return a list of List plan modifiers.
	ListPlanModifiers() []planmodifier.List
}

// AttributeWithMapPlanModifiers is an optional interface on Attribute which
// enables Map plan modifier support.
type AttributeWithMapPlanModifiers interface {
	fwschema.Attribute

	// MapPlanModifiers should return a list of Map plan modifiers.
	MapPlanModifiers() []planmodifier.Map
}

// AttributeWithNumberPlanModifiers is an optional interface on Attribute which
// enables Number plan modifier support.
type AttributeWithNumberPlanModifiers interface {
	fwschema.Attribute

	// NumberPlanModifiers should return a list of Number plan modifiers.
	NumberPlanModifiers() []planmodifier.Number
}

// AttributeWithObjectPlanModifiers is an optional interface on Attribute which
// enables Object plan modifier support.
type AttributeWithObjectPlanModifiers interface {
	fwschema.Attribute

	// ObjectPlanModifiers should return a list of Object plan modifiers.
	ObjectPlanModifiers() []planmodifier.Object
}

// AttributeWithSetPlanModifiers is an optional interface on Attribute which
// enables Set plan modifier support.
type AttributeWithSetPlanModifiers interface {
	fwschema.Attribute

	// SetPlanModifiers should return a list of Set plan modifiers.
	SetPlanModifiers() []planmodifier.Set
}

// AttributeWithStringPlanModifiers is an optional interface on Attribute which
// enables String plan modifier support.
type AttributeWithStringPlanModifiers interface {
	fwschema.Attribute

	// StringPlanModifiers should return a list of String plan modifiers.
	StringPlanModifiers() []planmodifier.String
}

// AttributeWithDynamicPlanModifiers is an optional interface on Attribute which
// enables Dynamic plan modifier support.
type AttributeWithDynamicPlanModifiers interface {
	fwschema.Attribute

	// DynamicPlanModifiers should return a list of Dynamic plan modifiers.
	DynamicPlanModifiers() []planmodifier.Dynamic
}
