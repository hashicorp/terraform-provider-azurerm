// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package toproto5

import (
	"context"

	"github.com/hashicorp/terraform-plugin-framework/internal/fwserver"
	"github.com/hashicorp/terraform-plugin-go/tfprotov5"
)

// ResourceMetadata returns the tfprotov5.ResourceMetadata for a
// fwserver.ResourceMetadata.
func ResourceMetadata(ctx context.Context, fw fwserver.ResourceMetadata) tfprotov5.ResourceMetadata {
	return tfprotov5.ResourceMetadata{
		TypeName: fw.TypeName,
	}
}
