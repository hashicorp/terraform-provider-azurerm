// Copyright IBM Corp. 2014, 2025
// SPDX-License-Identifier: MPL-2.0

package api_management_workspace

import (
	"context"
	"fmt"
	"time"

	"github.com/hashicorp/go-azure-helpers/lang/pointer"
	"github.com/hashicorp/go-azure-sdk/resource-manager/apimanagement/2024-05-01/workspace"
	"github.com/hashicorp/terraform-provider-azurerm/internal/sdk"
)

func (r Resource) Update() sdk.ResourceFunc {
	return sdk.ResourceFunc{
		Timeout: 30 * time.Minute,
		Func: func(ctx context.Context, metadata sdk.ResourceMetaData) error {
			client := metadata.Client.ApiManagement.WorkspaceClient

			id, err := workspace.ParseWorkspaceID(metadata.ResourceData.Id())
			if err != nil {
				return err
			}

			var model ApiManagementWorkspaceModel
			if err := metadata.Decode(&model); err != nil {
				return fmt.Errorf("decoding: %+v", err)
			}

			resp, err := client.Get(ctx, *id)
			if err != nil {
				return fmt.Errorf("retrieving %s: %+v", *id, err)
			}

			if resp.Model == nil {
				return fmt.Errorf("retrieving %s: `model` was nil", *id)
			}

			if resp.Model.Properties == nil {
				return fmt.Errorf("retrieving %s: `properties` was nil", *id)
			}

			properties := resp.Model
			if metadata.ResourceData.HasChange("description") {
				if model.Description != "" {
					properties.Properties.Description = pointer.To(model.Description)
				} else {
					properties.Properties.Description = nil
				}
			}

			if metadata.ResourceData.HasChange("display_name") {
				properties.Properties.DisplayName = model.DisplayName
			}

			if _, err := client.CreateOrUpdate(ctx, *id, *properties, workspace.DefaultCreateOrUpdateOperationOptions()); err != nil {
				return fmt.Errorf("updating %s: %+v", *id, err)
			}

			return nil
		},
	}
}
