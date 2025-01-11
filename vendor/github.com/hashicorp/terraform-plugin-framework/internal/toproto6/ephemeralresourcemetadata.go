// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package toproto6

import (
	"context"

	"github.com/hashicorp/terraform-plugin-framework/internal/fwserver"
	"github.com/hashicorp/terraform-plugin-go/tfprotov6"
)

// EphemeralResourceMetadata returns the tfprotov6.EphemeralResourceMetadata for a
// fwserver.EphemeralResourceMetadata.
func EphemeralResourceMetadata(ctx context.Context, fw fwserver.EphemeralResourceMetadata) tfprotov6.EphemeralResourceMetadata {
	return tfprotov6.EphemeralResourceMetadata{
		TypeName: fw.TypeName,
	}
}
