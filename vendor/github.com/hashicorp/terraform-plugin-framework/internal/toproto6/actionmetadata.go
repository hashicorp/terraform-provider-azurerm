// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package toproto6

import (
	"context"

	"github.com/hashicorp/terraform-plugin-framework/internal/fwserver"
	"github.com/hashicorp/terraform-plugin-go/tfprotov6"
)

// ActionMetadata returns the tfprotov6.ActionMetadata for a
// fwserver.ActionMetadata.
func ActionMetadata(ctx context.Context, fw fwserver.ActionMetadata) tfprotov6.ActionMetadata {
	return tfprotov6.ActionMetadata{
		TypeName: fw.TypeName,
	}
}
