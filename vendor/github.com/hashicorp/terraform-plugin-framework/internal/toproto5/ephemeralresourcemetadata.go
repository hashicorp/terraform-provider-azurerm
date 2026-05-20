// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package toproto5

import (
	"context"

	"github.com/hashicorp/terraform-plugin-framework/internal/fwserver"
	"github.com/hashicorp/terraform-plugin-go/tfprotov5"
)

// EphemeralResourceMetadata returns the tfprotov5.EphemeralResourceMetadata for a
// fwserver.EphemeralResourceMetadata.
func EphemeralResourceMetadata(ctx context.Context, fw fwserver.EphemeralResourceMetadata) tfprotov5.EphemeralResourceMetadata {
	return tfprotov5.EphemeralResourceMetadata{
		TypeName: fw.TypeName,
	}
}
