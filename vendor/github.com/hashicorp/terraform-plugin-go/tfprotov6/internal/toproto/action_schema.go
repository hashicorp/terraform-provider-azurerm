// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package toproto

import (
	"fmt"

	"github.com/hashicorp/terraform-plugin-go/tfprotov6"
	"github.com/hashicorp/terraform-plugin-go/tfprotov6/internal/tfplugin6"
)

func ActionSchema(in *tfprotov6.ActionSchema) *tfplugin6.ActionSchema {
	if in == nil {
		return nil
	}

	resp := &tfplugin6.ActionSchema{
		Schema: Schema(in.Schema),
	}

	switch actionSchemaType := in.Type.(type) {
	case tfprotov6.UnlinkedActionSchemaType:
		resp.Type = &tfplugin6.ActionSchema_Unlinked_{
			Unlinked: &tfplugin6.ActionSchema_Unlinked{},
		}
	case tfprotov6.LifecycleActionSchemaType:
		resp.Type = &tfplugin6.ActionSchema_Lifecycle_{
			Lifecycle: &tfplugin6.ActionSchema_Lifecycle{
				Executes:       tfplugin6.ActionSchema_Lifecycle_ExecutionOrder(actionSchemaType.Executes),
				LinkedResource: LinkedResourceSchema(actionSchemaType.LinkedResource),
			},
		}
	case tfprotov6.LinkedActionSchemaType:
		resp.Type = &tfplugin6.ActionSchema_Linked_{
			Linked: &tfplugin6.ActionSchema_Linked{
				LinkedResources: LinkedResourceSchemas(actionSchemaType.LinkedResources),
			},
		}
	default:
		// It is not currently possible to create tfprotov6.ActionSchemaType
		// implementations outside the tfprotov6 package. If this panic was reached,
		// it implies that a new event type was introduced and needs to be implemented
		// as a new case above.
		panic(fmt.Sprintf("unimplemented tfprotov6.ActionSchemaType type: %T", in.Type))
	}

	return resp
}
func LinkedResourceSchemas(in []*tfprotov6.LinkedResourceSchema) []*tfplugin6.ActionSchema_LinkedResource {
	resp := make([]*tfplugin6.ActionSchema_LinkedResource, 0, len(in))

	for _, schema := range in {
		resp = append(resp, LinkedResourceSchema(schema))
	}

	return resp
}

func LinkedResourceSchema(in *tfprotov6.LinkedResourceSchema) *tfplugin6.ActionSchema_LinkedResource {
	if in == nil {
		return nil
	}

	return &tfplugin6.ActionSchema_LinkedResource{
		TypeName:    in.TypeName,
		Description: in.Description,
	}
}
