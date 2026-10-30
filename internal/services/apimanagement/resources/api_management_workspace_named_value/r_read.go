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

func (r Resource) Read() sdk.ResourceFunc {
	return sdk.ResourceFunc{
		Timeout: 5 * time.Minute,
		Func: func(ctx context.Context, metadata sdk.ResourceMetaData) error {
			client := metadata.Client.ApiManagement.NamedValueClient_v2024_05_01

			id, err := namedvalue.ParseWorkspaceNamedValueID(metadata.ResourceData.Id())
			if err != nil {
				return err
			}

			var config ApiManagementWorkspaceNamedValueModel
			if err := metadata.Decode(&config); err != nil {
				return fmt.Errorf("decoding: %+v", err)
			}

			resp, err := client.WorkspaceNamedValueGet(ctx, *id)
			if err != nil {
				if response.WasNotFound(resp.HttpResponse) {
					return metadata.MarkAsGone(*id)
				}
				return fmt.Errorf("retrieving %s: %+v", *id, err)
			}

			model := ApiManagementWorkspaceNamedValueModel{
				Name:                     id.NamedValueId,
				ApiManagementWorkspaceId: workspace.NewWorkspaceID(id.SubscriptionId, id.ResourceGroupName, id.ServiceName, id.WorkspaceId).ID(),
			}

			if respModel := resp.Model; respModel != nil {
				if props := respModel.Properties; props != nil {
					model.DisplayName = props.DisplayName
					model.Secret = pointer.From(props.Secret)
					model.Tags = pointer.From(props.Tags)
					if model.Secret {
						// `value` is not retrievable when `secret` is `true`, so we use the config value
						model.Value = config.Value
					} else {
						model.Value = pointer.From(props.Value)
					}
					model.ValueFromKeyVault = flattenApiManagementWorkspaceNamedValueKeyVault(props.KeyVault)
				}
			}

			return metadata.Encode(&model)
		},
	}
}

func flattenApiManagementWorkspaceNamedValueKeyVault(input *namedvalue.KeyVaultContractProperties) []KeyVault {
	if input == nil {
		return []KeyVault{}
	}

	return []KeyVault{
		{
			SecretId:         pointer.From(input.SecretIdentifier),
			IdentityClientId: pointer.From(input.IdentityClientId),
		},
	}
}

type KeyVault struct {
	IdentityClientId string `tfschema:"identity_client_id"`
	SecretId         string `tfschema:"secret_id"`
}
