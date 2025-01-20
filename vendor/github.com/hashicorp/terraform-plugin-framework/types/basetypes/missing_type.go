// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package basetypes

import (
	"context"
	"fmt"

	"github.com/hashicorp/terraform-plugin-framework/attr"
	"github.com/hashicorp/terraform-plugin-go/tftypes"
)

var _ attr.Type = missingType{}

// missingType is a placeholder attr.Type implementation for when type
// information is missing. This type is never valid for real usage, which is why
// it is unexported, but it is primarily used by other base types when an
// expected attr.Type field is nil for panic prevention and troubleshooting.
// Ideally those other base type implementations would make it impossible to
// create a situation which needs this, but those exported APIs are protected by
// compatibility promises until a major version.
type missingType struct{}

// ApplyTerraform5AttributePathStep always returns an error.
func (t missingType) ApplyTerraform5AttributePathStep(step tftypes.AttributePathStep) (interface{}, error) {
	return nil, fmt.Errorf("cannot apply AttributePathStep %T to %s", step, t.String())
}

// Equal returns true if the given type is equivalent.
func (t missingType) Equal(o attr.Type) bool {
	_, ok := o.(missingType)

	return ok
}

// String returns a human readable string of the type.
func (t missingType) String() string {
	return "!!! MISSING TYPE !!!"
}

// TerraformType returns DynamicPseudoType.
func (t missingType) TerraformType(_ context.Context) tftypes.Type {
	// Ideally, upstream would implement an "invalid" primitive type for this
	// situation, but DynamicPseudoType is an alternative unexpected type in
	// the framework until it potentially implements its own dynamic type
	// handling.
	return tftypes.DynamicPseudoType
}

// ValueFromTerraform always returns an error.
func (t missingType) ValueFromTerraform(_ context.Context, _ tftypes.Value) (attr.Value, error) {
	return missingValue{}, fmt.Errorf("missing type information; cannot create value")
}

// ValueType returns the missingValue type.
func (t missingType) ValueType(_ context.Context) attr.Value {
	return missingValue{}
}
