// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package toproto6

import (
	"context"
	"fmt"

	actionschema "github.com/hashicorp/terraform-plugin-framework/action/schema"
	"github.com/hashicorp/terraform-plugin-go/tfprotov6"
)

// ActionSchema returns the *tfprotov6.ActionSchema equivalent of a ActionSchema.
func ActionSchema(ctx context.Context, s actionschema.SchemaType) (*tfprotov6.ActionSchema, error) {
	if s == nil {
		return nil, nil
	}

	configSchema, err := Schema(ctx, s)
	if err != nil {
		return nil, err
	}

	result := &tfprotov6.ActionSchema{
		Schema: configSchema,
	}

	// TODO:Actions: Implement linked and lifecycle action schema types
	switch s.(type) {
	case actionschema.UnlinkedSchema:
		result.Type = tfprotov6.UnlinkedActionSchemaType{}
	default:
		// It is not currently possible to create [actionschema.SchemaType]
		// implementations outside the "action/schema" package. If this error was reached,
		// it implies that a new event type was introduced and needs to be implemented
		// as a new case above.
		return nil, fmt.Errorf("unimplemented schema.SchemaType type: %T", s)
	}

	return result, nil
}
