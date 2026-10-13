// Copyright IBM Corp. 2014, 2025
// SPDX-License-Identifier: MPL-2.0

package api_management_workspace_api_version_set

import (
	"context"
	"fmt"
	"time"

	"github.com/hashicorp/go-azure-helpers/lang/pointer"
	"github.com/hashicorp/go-azure-helpers/lang/response"
	"github.com/hashicorp/go-azure-sdk/resource-manager/apimanagement/2024-05-01/apiversionset"
	"github.com/hashicorp/go-azure-sdk/resource-manager/apimanagement/2024-05-01/workspace"
	"github.com/hashicorp/terraform-provider-azurerm/internal/sdk"
)

func (r Resource) Create() sdk.ResourceFunc {
	return sdk.ResourceFunc{
		Timeout: 30 * time.Minute,
		Func: func(ctx context.Context, metadata sdk.ResourceMetaData) error {
			client := metadata.Client.ApiManagement.ApiVersionSetClient_v2024_05_01

			var model ApiManagementWorkspaceApiVersionSetModel
			if err := metadata.Decode(&model); err != nil {
				return fmt.Errorf("decoding: %+v", err)
			}

			workspaceId, err := workspace.ParseWorkspaceID(model.ApiManagementWorkspaceId)
			if err != nil {
				return err
			}

			id := apiversionset.NewWorkspaceApiVersionSetID(workspaceId.SubscriptionId, workspaceId.ResourceGroupName, workspaceId.ServiceName, workspaceId.WorkspaceId, model.Name)

			if !metadata.Client.Features.SkipImportCheckOnCreateAndAllowOverwritingExistingResources {
				existing, err := client.WorkspaceApiVersionSetGet(ctx, id)
				if err != nil && !response.WasNotFound(existing.HttpResponse) {
					return fmt.Errorf("checking for presence of existing %s: %+v", id, err)
				}

				if !response.WasNotFound(existing.HttpResponse) {
					return metadata.ResourceRequiresImport(r.ResourceType(), id)
				}
			}

			parameters := apiversionset.ApiVersionSetContract{
				Properties: &apiversionset.ApiVersionSetContractProperties{
					DisplayName:      model.DisplayName,
					VersioningScheme: apiversionset.VersioningScheme(model.VersioningScheme),
				},
			}

			if model.Description != "" {
				parameters.Properties.Description = pointer.To(model.Description)
			}

			if model.VersionHeaderName != "" {
				parameters.Properties.VersionHeaderName = pointer.To(model.VersionHeaderName)
			}

			if model.VersionQueryName != "" {
				parameters.Properties.VersionQueryName = pointer.To(model.VersionQueryName)
			}

			if _, err = client.WorkspaceApiVersionSetCreateOrUpdate(ctx, id, parameters, apiversionset.DefaultWorkspaceApiVersionSetCreateOrUpdateOperationOptions()); err != nil {
				return fmt.Errorf("creating %s: %+v", id, err)
			}

			metadata.SetID(id)

			return nil
		},
	}
}
