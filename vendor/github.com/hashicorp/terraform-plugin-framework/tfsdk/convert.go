// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package tfsdk

import (
	"context"
	"fmt"

	"github.com/hashicorp/terraform-plugin-framework/attr"
	"github.com/hashicorp/terraform-plugin-framework/diag"
)

// ConvertValue creates a new attr.Value of the attr.Type `typ`, using the data
// in `val`, which can be of any attr.Type so long as its TerraformType method
// returns a tftypes.Type that `typ`'s ValueFromTerraform method can accept.
func ConvertValue(ctx context.Context, val attr.Value, typ attr.Type) (attr.Value, diag.Diagnostics) {
	newVal, err := val.ToTerraformValue(ctx)
	if err != nil {
		return nil, diag.Diagnostics{diag.NewErrorDiagnostic("Error converting value",
			fmt.Sprintf("An unexpected error was encountered converting a %T to a %s. This is always a problem with the provider. Please tell the provider developers that %T ran into the following error during ToTerraformValue: %s", val, typ, val, err),
		)}
	}
	res, err := typ.ValueFromTerraform(ctx, newVal)
	if err != nil {
		return nil, diag.Diagnostics{diag.NewErrorDiagnostic("Error converting value",
			fmt.Sprintf("An unexpected error was encountered converting a %T to a %s. This is always a problem with the provider. Please tell the provider developers that %s returned the following error when calling ValueFromTerraform: %s", val, typ, typ, err),
		)}
	}
	return res, nil
}
