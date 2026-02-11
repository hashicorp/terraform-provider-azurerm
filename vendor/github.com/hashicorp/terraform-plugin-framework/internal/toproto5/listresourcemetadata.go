// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package toproto5

import (
	"context"

	"github.com/hashicorp/terraform-plugin-framework/internal/fwserver"
	"github.com/hashicorp/terraform-plugin-go/tfprotov5"
)

// ListResourceMetadata returns the tfprotov5.ListResourceMetadata for a
// fwserver.ListResourceMetadata.
func ListResourceMetadata(ctx context.Context, fw fwserver.ListResourceMetadata) tfprotov5.ListResourceMetadata {
	return tfprotov5.ListResourceMetadata{
		TypeName: fw.TypeName,
	}
}
