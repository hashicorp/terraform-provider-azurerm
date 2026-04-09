// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package toproto6

import (
	"context"

	"github.com/hashicorp/terraform-plugin-framework/internal/fwserver"
	"github.com/hashicorp/terraform-plugin-go/tfprotov6"
)

// ResourceMetadata returns the tfprotov6.ResourceMetadata for a
// fwserver.ResourceMetadata.
func ResourceMetadata(ctx context.Context, fw fwserver.ResourceMetadata) tfprotov6.ResourceMetadata {
	return tfprotov6.ResourceMetadata{
		TypeName: fw.TypeName,
	}
}
