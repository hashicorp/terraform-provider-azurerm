// Copyright IBM Corp. 2014, 2026
// SPDX-License-Identifier: MPL-2.0

package api_management_workspace_named_value

import (
	"context"
	"fmt"
	"time"

	"github.com/hashicorp/go-azure-helpers/lang/pointer"
	"github.com/hashicorp/go-azure-helpers/lang/response"
	"github.com/hashicorp/go-azure-sdk/resource-manager/apimanagement/2024-05-01/namedvalue"
	"github.com/hashicorp/go-azure-sdk/resource-manager/apimanagement/2024-05-01/workspace"
	"github.com/hashicorp/terraform-provider-azurerm/internal/sdk"
)

func (r Resource) Create() sdk.ResourceFunc {
	return sdk.ResourceFunc{
		Timeout: 30 * time.Minute,
		Func: func(ctx context.Context, metadata sdk.ResourceMetaData) error {
			client := metadata.Client.ApiManagement.NamedValueClient_v2024_05_01

			var model ApiManagementWorkspaceNamedValueModel
			if err := metadata.Decode(&model); err != nil {
				return fmt.Errorf("decoding: %+v", err)
			}

			workspaceId, err := workspace.ParseWorkspaceID(model.ApiManagementWorkspaceId)
			if err != nil {
				return err
			}

			id := namedvalue.NewWorkspaceNamedValueID(workspaceId.SubscriptionId, workspaceId.ResourceGroupName, workspaceId.ServiceName, workspaceId.WorkspaceId, model.Name)

			if !metadata.Client.Features.SkipImportCheckOnCreateAndAllowOverwritingExistingResources {
				existing, err := client.WorkspaceNamedValueGet(ctx, id)
				if err != nil && !response.WasNotFound(existing.HttpResponse) {
					return fmt.Errorf("checking for presence of existing %s: %+v", id, err)
				}

				if !response.WasNotFound(existing.HttpResponse) {
					return metadata.ResourceRequiresImport(r.ResourceType(), id)
				}
			}

			parameters := namedvalue.NamedValueCreateContract{
				Properties: &namedvalue.NamedValueCreateContractProperties{
					DisplayName: model.DisplayName,
					KeyVault:    expandApiManagementWorkspaceNamedValueKeyVault(model.ValueFromKeyVault),
					Secret:      pointer.To(model.Secret),
					Tags:        pointer.To(model.Tags),
				},
			}

			if model.Value != "" {
				parameters.Properties.Value = pointer.To(model.Value)
			}

			if err := client.WorkspaceNamedValueCreateOrUpdateCallbackThenPoll(ctx, id, parameters, namedvalue.DefaultWorkspaceNamedValueCreateOrUpdateOperationOptions(), metadata.SetIDCallback(&id)); err != nil {
				return fmt.Errorf("creating %s: %+v", id, err)
			}

			metadata.SetID(id)

			return nil
		},
	}
}

func expandApiManagementWorkspaceNamedValueKeyVault(inputs []KeyVault) *namedvalue.KeyVaultContractCreateProperties {
	if len(inputs) == 0 {
		return nil
	}
	input := inputs[0]

	result := namedvalue.KeyVaultContractCreateProperties{
		SecretIdentifier: pointer.To(input.SecretId),
	}

	if input.IdentityClientId != "" {
		result.IdentityClientId = pointer.To(input.IdentityClientId)
	}

	return &result
}
