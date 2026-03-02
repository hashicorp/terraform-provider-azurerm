// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package toproto5

import (
	"context"
	"sort"

	"github.com/hashicorp/terraform-plugin-framework/internal/fwschema"
	"github.com/hashicorp/terraform-plugin-go/tfprotov5"
	"github.com/hashicorp/terraform-plugin-go/tftypes"
)

// IdentitySchema returns the *tfprotov5.ResourceIdentitySchema equivalent of a Schema.
func IdentitySchema(ctx context.Context, s fwschema.Schema) (*tfprotov5.ResourceIdentitySchema, error) {
	if s == nil {
		return nil, nil
	}

	result := &tfprotov5.ResourceIdentitySchema{
		Version: s.GetVersion(),
	}

	attrs := make([]*tfprotov5.ResourceIdentitySchemaAttribute, 0)

	for name, attr := range s.GetAttributes() {
		a, err := IdentitySchemaAttribute(ctx, name, tftypes.NewAttributePath().WithAttributeName(name), attr)

		if err != nil {
			return nil, err
		}

		attrs = append(attrs, a)
	}

	sort.Slice(attrs, func(i, j int) bool {
		if attrs[i] == nil {
			return true
		}

		if attrs[j] == nil {
			return false
		}

		return attrs[i].Name < attrs[j].Name
	})

	result.IdentityAttributes = attrs

	return result, nil
}
