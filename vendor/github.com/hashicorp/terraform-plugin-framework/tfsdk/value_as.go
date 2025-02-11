// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package tfsdk

import (
	"context"
	"fmt"

	"github.com/hashicorp/terraform-plugin-framework/attr"
	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/internal/reflect"
	"github.com/hashicorp/terraform-plugin-framework/path"
)

// ValueAs takes the attr.Value `val` and populates the Go value `target` with its content.
//
// This is achieved using reflection rules provided by the internal/reflect package.
func ValueAs(ctx context.Context, val attr.Value, target interface{}) diag.Diagnostics {
	if reflect.IsGenericAttrValue(ctx, target) {
		//nolint:forcetypeassert // Type assertion is guaranteed by the above `reflect.IsGenericAttrValue` function
		*(target.(*attr.Value)) = val
		return nil
	}
	raw, err := val.ToTerraformValue(ctx)
	if err != nil {
		return diag.Diagnostics{diag.NewErrorDiagnostic("Error converting value",
			fmt.Sprintf("An unexpected error was encountered converting a %T to its equivalent Terraform representation. This is always a bug in the provider.\n\nError: %s", val, err))}
	}
	return reflect.Into(ctx, val.Type(ctx), raw, target, reflect.Options{}, path.Empty())
}
