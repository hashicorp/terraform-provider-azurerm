// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package totftypes

import (
	"context"
	"fmt"

	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-go/tftypes"
)

// AttributePathStep returns the tftypes.AttributePathStep equivalent of an
// path.PathStep. An error is returned instead of diag.Diagnostics so callers
// can include appropriate logical context about when the error occurred.
func AttributePathStep(ctx context.Context, fw path.PathStep) (tftypes.AttributePathStep, error) {
	switch fw := fw.(type) {
	case path.PathStepAttributeName:
		return tftypes.AttributeName(string(fw)), nil
	case path.PathStepElementKeyInt:
		return tftypes.ElementKeyInt(int64(fw)), nil
	case path.PathStepElementKeyString:
		return tftypes.ElementKeyString(string(fw)), nil
	case path.PathStepElementKeyValue:
		tfTypesValue, err := fw.Value.ToTerraformValue(ctx)

		if err != nil {
			return nil, fmt.Errorf("unable to convert attr.Value (%s) to tftypes.Value: %w", fw.Value.String(), err)
		}

		return tftypes.ElementKeyValue(tfTypesValue), nil
	default:
		return nil, fmt.Errorf("unknown path.PathStep: %#v", fw)
	}
}
