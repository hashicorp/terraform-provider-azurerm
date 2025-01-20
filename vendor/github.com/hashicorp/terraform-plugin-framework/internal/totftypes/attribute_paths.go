// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package totftypes

import (
	"context"

	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-go/tftypes"
)

// AttributePaths returns the []*tftypes.AttributePath equivalent of a path.Paths.
func AttributePaths(ctx context.Context, fw path.Paths) ([]*tftypes.AttributePath, diag.Diagnostics) {
	if fw == nil {
		return nil, nil
	}

	result := make([]*tftypes.AttributePath, 0, len(fw))

	for _, path := range fw {
		tfType, diags := AttributePath(ctx, path)

		if diags.HasError() {
			return result, diags
		}

		result = append(result, tfType)
	}

	return result, nil
}
