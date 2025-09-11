// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package toproto

import (
	"fmt"

	"github.com/hashicorp/terraform-plugin-go/tfprotov5"
	"github.com/hashicorp/terraform-plugin-go/tfprotov5/internal/tfplugin5"
)

func ActionSchema(in *tfprotov5.ActionSchema) *tfplugin5.ActionSchema {
	if in == nil {
		return nil
	}

	resp := &tfplugin5.ActionSchema{
		Schema: Schema(in.Schema),
	}

	switch actionSchemaType := in.Type.(type) {
	case tfprotov5.UnlinkedActionSchemaType:
		resp.Type = &tfplugin5.ActionSchema_Unlinked_{
			Unlinked: &tfplugin5.ActionSchema_Unlinked{},
		}
	case tfprotov5.LifecycleActionSchemaType:
		resp.Type = &tfplugin5.ActionSchema_Lifecycle_{
			Lifecycle: &tfplugin5.ActionSchema_Lifecycle{
				Executes:       tfplugin5.ActionSchema_Lifecycle_ExecutionOrder(actionSchemaType.Executes),
				LinkedResource: LinkedResourceSchema(actionSchemaType.LinkedResource),
			},
		}
	case tfprotov5.LinkedActionSchemaType:
		resp.Type = &tfplugin5.ActionSchema_Linked_{
			Linked: &tfplugin5.ActionSchema_Linked{
				LinkedResources: LinkedResourceSchemas(actionSchemaType.LinkedResources),
			},
		}
	default:
		// It is not currently possible to create tfprotov5.ActionSchemaType
		// implementations outside the tfprotov5 package. If this panic was reached,
		// it implies that a new event type was introduced and needs to be implemented
		// as a new case above.
		panic(fmt.Sprintf("unimplemented tfprotov5.ActionSchemaType type: %T", in.Type))
	}

	return resp
}
func LinkedResourceSchemas(in []*tfprotov5.LinkedResourceSchema) []*tfplugin5.ActionSchema_LinkedResource {
	resp := make([]*tfplugin5.ActionSchema_LinkedResource, 0, len(in))

	for _, schema := range in {
		resp = append(resp, LinkedResourceSchema(schema))
	}

	return resp
}

func LinkedResourceSchema(in *tfprotov5.LinkedResourceSchema) *tfplugin5.ActionSchema_LinkedResource {
	if in == nil {
		return nil
	}

	return &tfplugin5.ActionSchema_LinkedResource{
		TypeName:    in.TypeName,
		Description: in.Description,
	}
}
