// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package fwschemadata

import (
	"context"
	"fmt"

	"github.com/hashicorp/terraform-plugin-go/tftypes"
)

// TerraformValueAtTerraformPath returns the tftypes.Value at a given
// tftypes.AttributePath or an error.
func (d Data) TerraformValueAtTerraformPath(_ context.Context, path *tftypes.AttributePath) (tftypes.Value, error) {
	rawValue, remaining, err := tftypes.WalkAttributePath(d.TerraformValue, path)

	if err != nil {
		return tftypes.Value{}, fmt.Errorf("%v still remains in the path: %w", remaining, err)
	}

	attrValue, ok := rawValue.(tftypes.Value)

	if !ok {
		return tftypes.Value{}, fmt.Errorf("got non-tftypes.Value result %v", rawValue)
	}

	return attrValue, err
}
