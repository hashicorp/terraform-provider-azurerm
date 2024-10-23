// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package totftypes

import (
	"context"
	"fmt"

	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-go/tftypes"
)

// AttributePath returns the *tftypes.AttributePath equivalent of a path.Path.
func AttributePath(ctx context.Context, fw path.Path) (*tftypes.AttributePath, diag.Diagnostics) {
	var tfTypeSteps []tftypes.AttributePathStep

	for _, step := range fw.Steps() {
		tfTypeStep, err := AttributePathStep(ctx, step)

		if err != nil {
			return nil, diag.Diagnostics{
				diag.NewErrorDiagnostic(
					"Unable to Convert Attribute Path",
					"An unexpected error occurred while trying to convert an attribute path. "+
						"This is either an error in terraform-plugin-framework or a custom attribute type used by the provider. "+
						"Please report the following to the provider developers.\n\n"+
						// Since this is an error with the attribute path
						// conversion, we cannot return a protocol path-based
						// diagnostic. Returning a framework human-readable
						// representation seems like the next best thing to do.
						fmt.Sprintf("Attribute Path: %s\n", fw.String())+
						fmt.Sprintf("Original Error: %s", err),
				),
			}
		}

		tfTypeSteps = append(tfTypeSteps, tfTypeStep)
	}

	return tftypes.NewAttributePathWithSteps(tfTypeSteps), nil
}
