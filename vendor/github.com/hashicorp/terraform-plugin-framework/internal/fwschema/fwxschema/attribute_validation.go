// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package fwxschema

import (
	"github.com/hashicorp/terraform-plugin-framework/internal/fwschema"
	"github.com/hashicorp/terraform-plugin-framework/schema/validator"
)

// AttributeWithBoolValidators is an optional interface on Attribute which
// enables Bool validation support.
type AttributeWithBoolValidators interface {
	fwschema.Attribute

	// BoolValidators should return a list of Bool validators.
	BoolValidators() []validator.Bool
}

// AttributeWithFloat64Validators is an optional interface on Attribute which
// enables Float64 validation support.
type AttributeWithFloat64Validators interface {
	fwschema.Attribute

	// Float64Validators should return a list of Float64 validators.
	Float64Validators() []validator.Float64
}

// AttributeWithInt64Validators is an optional interface on Attribute which
// enables Int64 validation support.
type AttributeWithInt64Validators interface {
	fwschema.Attribute

	// Int64Validators should return a list of Int64 validators.
	Int64Validators() []validator.Int64
}

// AttributeWithListValidators is an optional interface on Attribute which
// enables List validation support.
type AttributeWithListValidators interface {
	fwschema.Attribute

	// ListValidators should return a list of List validators.
	ListValidators() []validator.List
}

// AttributeWithMapValidators is an optional interface on Attribute which
// enables Map validation support.
type AttributeWithMapValidators interface {
	fwschema.Attribute

	// MapValidators should return a list of Map validators.
	MapValidators() []validator.Map
}

// AttributeWithNumberValidators is an optional interface on Attribute which
// enables Number validation support.
type AttributeWithNumberValidators interface {
	fwschema.Attribute

	// NumberValidators should return a list of Number validators.
	NumberValidators() []validator.Number
}

// AttributeWithObjectValidators is an optional interface on Attribute which
// enables Object validation support.
type AttributeWithObjectValidators interface {
	fwschema.Attribute

	// ObjectValidators should return a list of Object validators.
	ObjectValidators() []validator.Object
}

// AttributeWithSetValidators is an optional interface on Attribute which
// enables Set validation support.
type AttributeWithSetValidators interface {
	fwschema.Attribute

	// SetValidators should return a list of Set validators.
	SetValidators() []validator.Set
}

// AttributeWithStringValidators is an optional interface on Attribute which
// enables String validation support.
type AttributeWithStringValidators interface {
	fwschema.Attribute

	// StringValidators should return a list of String validators.
	StringValidators() []validator.String
}

// AttributeWithDynamicValidators is an optional interface on Attribute which
// enables Dynamic validation support.
type AttributeWithDynamicValidators interface {
	fwschema.Attribute

	// DynamicValidators should return a list of Dynamic validators.
	DynamicValidators() []validator.Dynamic
}
