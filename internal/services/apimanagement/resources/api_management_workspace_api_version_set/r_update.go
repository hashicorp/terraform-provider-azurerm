// Copyright IBM Corp. 2014, 2025
// SPDX-License-Identifier: MPL-2.0

package api_management_workspace_api_version_set

import (
	"context"
	"fmt"
	"time"

	"github.com/hashicorp/go-azure-helpers/lang/pointer"
	"github.com/hashicorp/go-azure-sdk/resource-manager/apimanagement/2024-05-01/apiversionset"
	"github.com/hashicorp/terraform-provider-azurerm/internal/sdk"
)

func (r Resource) Update() sdk.ResourceFunc {
	return sdk.ResourceFunc{
		Timeout: 30 * time.Minute,
		Func: func(ctx context.Context, metadata sdk.ResourceMetaData) error {
			client := metadata.Client.ApiManagement.ApiVersionSetClient_v2024_05_01

			var model ApiManagementWorkspaceApiVersionSetModel
			if err := metadata.Decode(&model); err != nil {
				return fmt.Errorf("decoding: %+v", err)
			}

			id, err := apiversionset.ParseWorkspaceApiVersionSetID(metadata.ResourceData.Id())
			if err != nil {
				return err
			}

			resp, err := client.WorkspaceApiVersionSetGet(ctx, *id)
			if err != nil {
				return fmt.Errorf("retrieving %s: %+v", id, err)
			}

			if resp.Model == nil {
				return fmt.Errorf("retrieving %s: `model` was nil", id)
			}

			if resp.Model.Properties == nil {
				return fmt.Errorf("retrieving %s: `properties` was nil", id)
			}

			properties := resp.Model.Properties
			if metadata.ResourceData.HasChange("display_name") {
				properties.DisplayName = model.DisplayName
			}

			if metadata.ResourceData.HasChange("versioning_scheme") {
				properties.VersioningScheme = apiversionset.VersioningScheme(model.VersioningScheme)
			}

			if metadata.ResourceData.HasChange("description") {
				properties.Description = pointer.To(model.Description)
			}

			if metadata.ResourceData.HasChange("version_header_name") {
				properties.VersionHeaderName = nil
				if model.VersionHeaderName != "" {
					properties.VersionHeaderName = pointer.To(model.VersionHeaderName)
				}
			}

			if metadata.ResourceData.HasChange("version_query_name") {
				properties.VersionQueryName = nil
				if model.VersionQueryName != "" {
					properties.VersionQueryName = pointer.To(model.VersionQueryName)
				}
			}

			resp.Model.Properties = properties
			if _, err = client.WorkspaceApiVersionSetCreateOrUpdate(ctx, *id, *resp.Model, apiversionset.DefaultWorkspaceApiVersionSetCreateOrUpdateOperationOptions()); err != nil {
				return fmt.Errorf("updating %s: %+v", id, err)
			}

			return nil
		},
	}
}
