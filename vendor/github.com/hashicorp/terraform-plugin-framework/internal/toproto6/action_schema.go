// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package toproto6

import (
	"context"

	"github.com/hashicorp/terraform-plugin-go/tfprotov6"

	actionschema "github.com/hashicorp/terraform-plugin-framework/action/schema"
)

// ActionSchema returns the *tfprotov6.ActionSchema equivalent of a ActionSchema.
func ActionSchema(ctx context.Context, s actionschema.Schema) (*tfprotov6.ActionSchema, error) {
	configSchema, err := Schema(ctx, s)
	if err != nil {
		return nil, err
	}

	result := &tfprotov6.ActionSchema{
		Schema: configSchema,
	}

	return result, nil
}
