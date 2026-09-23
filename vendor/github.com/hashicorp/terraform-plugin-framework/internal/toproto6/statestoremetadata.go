// Copyright IBM Corp. 2021, 2026
// SPDX-License-Identifier: MPL-2.0

package toproto6

import (
	"context"

	"github.com/hashicorp/terraform-plugin-framework/internal/fwserver"
	"github.com/hashicorp/terraform-plugin-go/tfprotov6"
)

// StateStoreMetadata returns the tfprotov6.StateStoreMetadata for a
// fwserver.StateStoreMetadata.
func StateStoreMetadata(ctx context.Context, fw fwserver.StateStoreMetadata) tfprotov6.StateStoreMetadata {
	return tfprotov6.StateStoreMetadata{
		TypeName: fw.TypeName,
	}
}
