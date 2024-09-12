// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package fwschemadata

import (
	"context"

	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/internal/reflect"
	"github.com/hashicorp/terraform-plugin-framework/path"
)

// Get populates the struct passed as `target` with the entire state.
func (d Data) Get(ctx context.Context, target any) diag.Diagnostics {
	return reflect.Into(ctx, d.Schema.Type(), d.TerraformValue, target, reflect.Options{}, path.Empty())
}
