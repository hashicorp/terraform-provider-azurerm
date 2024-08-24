// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package toproto6

import (
	"context"

	"github.com/hashicorp/terraform-plugin-framework/internal/fwserver"
	"github.com/hashicorp/terraform-plugin-go/tfprotov6"
)

// DataSourceMetadata returns the tfprotov6.DataSourceMetadata for a
// fwserver.DataSourceMetadata.
func DataSourceMetadata(ctx context.Context, fw fwserver.DataSourceMetadata) tfprotov6.DataSourceMetadata {
	return tfprotov6.DataSourceMetadata{
		TypeName: fw.TypeName,
	}
}
