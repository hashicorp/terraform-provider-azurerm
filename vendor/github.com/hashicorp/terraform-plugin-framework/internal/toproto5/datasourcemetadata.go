// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package toproto5

import (
	"context"

	"github.com/hashicorp/terraform-plugin-framework/internal/fwserver"
	"github.com/hashicorp/terraform-plugin-go/tfprotov5"
)

// DataSourceMetadata returns the tfprotov5.DataSourceMetadata for a
// fwserver.DataSourceMetadata.
func DataSourceMetadata(ctx context.Context, fw fwserver.DataSourceMetadata) tfprotov5.DataSourceMetadata {
	return tfprotov5.DataSourceMetadata{
		TypeName: fw.TypeName,
	}
}
