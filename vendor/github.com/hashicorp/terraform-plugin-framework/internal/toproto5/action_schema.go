// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package toproto5

import (
	"context"

	"github.com/hashicorp/terraform-plugin-go/tfprotov5"

	actionschema "github.com/hashicorp/terraform-plugin-framework/action/schema"
)

// ActionSchema returns the *tfprotov5.ActionSchema equivalent of a ActionSchema.
func ActionSchema(ctx context.Context, s actionschema.Schema) (*tfprotov5.ActionSchema, error) {
	configSchema, err := Schema(ctx, s)
	if err != nil {
		return nil, err
	}

	result := &tfprotov5.ActionSchema{
		Schema: configSchema,
	}

	return result, nil
}
