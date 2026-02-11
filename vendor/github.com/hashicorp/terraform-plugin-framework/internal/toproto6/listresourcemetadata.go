// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package toproto6

import (
	"context"

	"github.com/hashicorp/terraform-plugin-framework/internal/fwserver"
	"github.com/hashicorp/terraform-plugin-go/tfprotov6"
)

// ListResourceMetadata returns the tfprotov6.ListResourceMetadata for a
// fwserver.ListResourceMetadata.
func ListResourceMetadata(ctx context.Context, fw fwserver.ListResourceMetadata) tfprotov6.ListResourceMetadata {
	return tfprotov6.ListResourceMetadata{
		TypeName: fw.TypeName,
	}
}
