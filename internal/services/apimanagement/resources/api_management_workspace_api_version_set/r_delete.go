// Copyright IBM Corp. 2014, 2025
// SPDX-License-Identifier: MPL-2.0

package api_management_workspace_api_version_set

import (
	"context"
	"fmt"
	"time"

	"github.com/hashicorp/go-azure-sdk/resource-manager/apimanagement/2024-05-01/apiversionsets"
	"github.com/hashicorp/terraform-provider-azurerm/internal/sdk"
)

func (r Resource) Delete() sdk.ResourceFunc {
	return sdk.ResourceFunc{
		Timeout: 30 * time.Minute,
		Func: func(ctx context.Context, metadata sdk.ResourceMetaData) error {
			client := metadata.Client.ApiManagement.ApiVersionSetsClient_v2024_05_01

			id, err := apiversionsets.ParseWorkspaceApiVersionSetID(metadata.ResourceData.Id())
			if err != nil {
				return err
			}

			if _, err = client.WorkspaceApiVersionSetDelete(ctx, *id, apiversionsets.DefaultWorkspaceApiVersionSetDeleteOperationOptions()); err != nil {
				return fmt.Errorf("deleting %s: %+v", *id, err)
			}

			return nil
		},
	}
}
