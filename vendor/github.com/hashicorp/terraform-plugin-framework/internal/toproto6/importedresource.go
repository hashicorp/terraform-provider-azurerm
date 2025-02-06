// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package toproto6

import (
	"context"

	"github.com/hashicorp/terraform-plugin-go/tfprotov6"

	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/internal/fwserver"
)

// ImportedResource returns the *tfprotov6.ImportedResource equivalent of a
// *fwserver.ImportedResource.
func ImportedResource(ctx context.Context, fw *fwserver.ImportedResource) (*tfprotov6.ImportedResource, diag.Diagnostics) {
	if fw == nil {
		return nil, nil
	}

	proto6 := &tfprotov6.ImportedResource{
		TypeName: fw.TypeName,
	}

	state, diags := State(ctx, &fw.State)

	proto6.State = state

	newPrivate, privateDiags := fw.Private.Bytes(ctx)

	diags = append(diags, privateDiags...)
	proto6.Private = newPrivate

	return proto6, diags
}
