// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package toproto5

import (
	"context"

	"github.com/hashicorp/terraform-plugin-framework/internal/fwserver"
	"github.com/hashicorp/terraform-plugin-go/tfprotov5"
)

// ActionMetadata returns the tfprotov5.ActionMetadata for a
// fwserver.ActionMetadata.
func ActionMetadata(ctx context.Context, fw fwserver.ActionMetadata) tfprotov5.ActionMetadata {
	return tfprotov5.ActionMetadata{
		TypeName: fw.TypeName,
	}
}
