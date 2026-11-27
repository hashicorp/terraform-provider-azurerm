// Copyright IBM Corp. 2014, 2026
// SPDX-License-Identifier: MPL-2.0

package api_management_workspace_policy_fragment

import (
	"context"
	"fmt"
	"time"

	"github.com/hashicorp/go-azure-helpers/lang/pointer"
	"github.com/hashicorp/go-azure-helpers/lang/response"
	"github.com/hashicorp/go-azure-sdk/resource-manager/apimanagement/2024-05-01/policyfragment"
	"github.com/hashicorp/terraform-provider-azurerm/internal/features"
	"github.com/hashicorp/terraform-provider-azurerm/internal/sdk"
	apimanagementClient "github.com/hashicorp/terraform-provider-azurerm/internal/services/apimanagement/client"
)

func (r Resource) Create() sdk.ResourceFunc {
	return sdk.ResourceFunc{
		Timeout: 30 * time.Minute,
		Func: func(ctx context.Context, metadata sdk.ResourceMetaData) error {
			client := apimanagementClient.From(metadata.Client).PolicyFragmentClient_v2024_05_01

			var model ApiManagementWorkspacePolicyFragmentModel
			if err := metadata.Decode(&model); err != nil {
				return fmt.Errorf("decoding: %+v", err)
			}

			workspaceId, err := policyfragment.ParseWorkspaceID(model.ApiManagementWorkspaceId)
			if err != nil {
				return err
			}

			id := policyfragment.NewWorkspacePolicyFragmentID(workspaceId.SubscriptionId, workspaceId.ResourceGroupName, workspaceId.ServiceName, workspaceId.WorkspaceId, model.Name)

			if !features.From(metadata.Client).SkipImportCheckOnCreateAndAllowOverwritingExistingResources {
				existing, err := client.WorkspacePolicyFragmentGet(ctx, id, policyfragment.WorkspacePolicyFragmentGetOperationOptions{
					Format: pointer.ToEnum[policyfragment.PolicyFragmentContentFormat](model.XmlFormat),
				})
				if err != nil && !response.WasNotFound(existing.HttpResponse) {
					return fmt.Errorf("checking for presence of existing %s: %+v", id, err)
				}

				if !response.WasNotFound(existing.HttpResponse) {
					return metadata.ResourceRequiresImport(r.ResourceType(), id)
				}
			}

			parameters := policyfragment.PolicyFragmentContract{
				Properties: &policyfragment.PolicyFragmentContractProperties{
					Format: pointer.ToEnum[policyfragment.PolicyFragmentContentFormat](model.XmlFormat),
					Value:  model.XmlContent,
				},
			}

			if model.Description != "" {
				parameters.Properties.Description = pointer.To(model.Description)
			}

			if err := client.WorkspacePolicyFragmentCreateOrUpdateCallbackThenPoll(ctx, id, parameters, policyfragment.WorkspacePolicyFragmentCreateOrUpdateOperationOptions{}, metadata.SetIDCallback(&id)); err != nil {
				return fmt.Errorf("creating %s: %+v", id, err)
			}

			metadata.SetID(id)
			return nil
		},
	}
}
