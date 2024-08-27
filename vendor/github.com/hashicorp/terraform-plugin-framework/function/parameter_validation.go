// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package function

// ParameterWithBoolValidators is an optional interface on Parameter which
// enables Bool validation support.
type ParameterWithBoolValidators interface {
	Parameter

	// GetValidators should return a list of Bool validators.
	GetValidators() []BoolParameterValidator
}

// ParameterWithInt64Validators is an optional interface on Parameter which
// enables Int64 validation support.
type ParameterWithInt64Validators interface {
	Parameter

	// GetValidators should return a list of Int64 validators.
	GetValidators() []Int64ParameterValidator
}

// ParameterWithFloat64Validators is an optional interface on Parameter which
// enables Float64 validation support.
type ParameterWithFloat64Validators interface {
	Parameter

	// GetValidators should return a list of Float64 validators.
	GetValidators() []Float64ParameterValidator
}

// ParameterWithDynamicValidators is an optional interface on Parameter which
// enables Dynamic validation support.
type ParameterWithDynamicValidators interface {
	Parameter

	// GetValidators should return a list of Dynamic validators.
	GetValidators() []DynamicParameterValidator
}

// ParameterWithListValidators is an optional interface on Parameter which
// enables List validation support.
type ParameterWithListValidators interface {
	Parameter

	// GetValidators should return a list of List validators.
	GetValidators() []ListParameterValidator
}

// ParameterWithMapValidators is an optional interface on Parameter which
// enables Map validation support.
type ParameterWithMapValidators interface {
	Parameter

	// GetValidators should return a list of Map validators.
	GetValidators() []MapParameterValidator
}

// ParameterWithNumberValidators is an optional interface on Parameter which
// enables Number validation support.
type ParameterWithNumberValidators interface {
	Parameter

	// GetValidators should return a list of Map validators.
	GetValidators() []NumberParameterValidator
}

// ParameterWithObjectValidators is an optional interface on Parameter which
// enables Object validation support.
type ParameterWithObjectValidators interface {
	Parameter

	// GetValidators should return a list of Object validators.
	GetValidators() []ObjectParameterValidator
}

// ParameterWithSetValidators is an optional interface on Parameter which
// enables Set validation support.
type ParameterWithSetValidators interface {
	Parameter

	// GetValidators should return a list of Set validators.
	GetValidators() []SetParameterValidator
}

// ParameterWithStringValidators is an optional interface on Parameter which
// enables String validation support.
type ParameterWithStringValidators interface {
	Parameter

	// GetValidators should return a list of String validators.
	GetValidators() []StringParameterValidator
}
