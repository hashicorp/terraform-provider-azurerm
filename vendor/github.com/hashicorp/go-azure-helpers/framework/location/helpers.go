// Copyright IBM Corp. 2018, 2025
// SPDX-License-Identifier: MPL-2.0

package location

import (
	"strings"

	"github.com/hashicorp/go-azure-helpers/lang/pointer"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-framework/types/basetypes"
)

// Normalize transforms the human readable Azure Region/Location names (e.g. `West US`)
// into the canonical value to allow comparisons between user-code and API Responses
func Normalize(input string) string {
	return strings.ReplaceAll(strings.ToLower(input), " ", "")
}

// NormalizeNilable normalizes the Location field even if it's nil to ensure this field
// can always have a value
func NormalizeNilable(input *string) string {
	return Normalize(pointer.From(input))
}

// NormalizeValue returns a Framework compatible StringValue for the location
func NormalizeValue(input string) types.String {
	return basetypes.NewStringValue(Normalize(input))
}
