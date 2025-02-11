// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package basetypes

import (
	"context"

	"github.com/hashicorp/terraform-plugin-framework/attr"
	"github.com/hashicorp/terraform-plugin-go/tftypes"
)

var _ attr.Value = missingValue{}

// missingValue is a placeholder attr.Value implementation for when type
// information is missing. This type is never valid for real usage, which is why
// it is unexported, but it is primarily used by other base types when an
// expected attr.Value field is nil for panic prevention and troubleshooting.
// Ideally those other base type implementations would make it impossible to
// create a situation which needs this, but those exported APIs are protected by
// compatibility promises until a major version.
type missingValue struct{}

// Equal returns true if the given value is a missingValue.
func (v missingValue) Equal(o attr.Value) bool {
	_, ok := o.(missingValue)

	return ok
}

// IsNull returns false.
func (v missingValue) IsNull() bool {
	// Short of causing a panic, this method must choose a return value and
	// false was chosen so it is always "known".
	return false
}

// IsUnknown returns false.
func (v missingValue) IsUnknown() bool {
	// Short of causing a panic, this method must choose a return value and
	// false was chosen so it is always "known".
	return false
}

// String returns a human-readable representation of the value.
//
// The string returned here is not protected by any compatibility guarantees,
// and is intended for logging and error reporting.
func (v missingValue) String() string {
	return "!!! MISSING VALUE !!!"
}

// ToTerraformValue always returns an error.
func (v missingValue) ToTerraformValue(_ context.Context) (tftypes.Value, error) {
	return tftypes.Value{}, nil
}

// Type returns missingType.
func (v missingValue) Type(_ context.Context) attr.Type {
	return missingType{}
}
