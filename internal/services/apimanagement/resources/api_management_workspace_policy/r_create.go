// Copyright IBM Corp. 2014, 2026
// SPDX-License-Identifier: MPL-2.0

package api_management_workspace_policy

import (
	"context"
	"fmt"
	"time"

	"github.com/hashicorp/go-azure-helpers/lang/pointer"
	"github.com/hashicorp/go-azure-helpers/lang/response"
	"github.com/hashicorp/go-azure-sdk/resource-manager/apimanagement/2024-05-01/workspacepolicy"
	"github.com/hashicorp/terraform-provider-azurerm/internal/features"
	"github.com/hashicorp/terraform-provider-azurerm/internal/sdk"
	apimanagementClient "github.com/hashicorp/terraform-provider-azurerm/internal/services/apimanagement/client"
)

func (r Resource) Create() sdk.ResourceFunc {
	return sdk.ResourceFunc{
		Timeout: 30 * time.Minute,
		Func: func(ctx context.Context, metadata sdk.ResourceMetaData) error {
			client := apimanagementClient.From(metadata.Client).WorkspacePolicyClient

			var model ApiManagementWorkspacePolicyModel
			if err := metadata.Decode(&model); err != nil {
				return fmt.Errorf("decoding: %+v", err)
			}

			workspaceId, err := workspacepolicy.ParseWorkspaceID(model.ApiManagementWorkspaceId)
			if err != nil {
				return err
			}

			if !features.From(metadata.Client).SkipImportCheckOnCreateAndAllowOverwritingExistingResources {
				existing, err := client.Get(ctx, *workspaceId, workspacepolicy.GetOperationOptions{})
				if err != nil && !response.WasNotFound(existing.HttpResponse) {
					return fmt.Errorf("checking for presence of existing %s: %+v", *workspaceId, err)
				}

				if !response.WasNotFound(existing.HttpResponse) {
					return metadata.ResourceRequiresImport(r.ResourceType(), *workspaceId)
				}
			}

			parameters := workspacepolicy.PolicyContract{}

			if model.XmlLink != "" {
				parameters.Properties = &workspacepolicy.PolicyContractProperties{
					Format: pointer.To(workspacepolicy.PolicyContentFormatRawxmlNegativelink),
					Value:  model.XmlLink,
				}
			}

			if model.XmlContent != "" {
				parameters.Properties = &workspacepolicy.PolicyContractProperties{
					Format: pointer.To(workspacepolicy.PolicyContentFormatRawxml),
					Value:  model.XmlContent,
				}
			}

			if _, err := client.CreateOrUpdate(ctx, *workspaceId, parameters, workspacepolicy.DefaultCreateOrUpdateOperationOptions()); err != nil {
				return fmt.Errorf("creating %s: %+v", *workspaceId, err)
			}

			metadata.SetID(*workspaceId)
			return nil
		},
	}
}
