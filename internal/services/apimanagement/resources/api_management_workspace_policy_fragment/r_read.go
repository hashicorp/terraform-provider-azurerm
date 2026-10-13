// Copyright IBM Corp. 2014, 2025
// SPDX-License-Identifier: MPL-2.0

package api_management_workspace_policy_fragment

import (
	"context"
	"fmt"
	"time"

	"github.com/hashicorp/go-azure-helpers/lang/pointer"
	"github.com/hashicorp/go-azure-helpers/lang/response"
	"github.com/hashicorp/go-azure-sdk/resource-manager/apimanagement/2024-05-01/policyfragment"
	"github.com/hashicorp/terraform-provider-azurerm/internal/sdk"
)

func (r Resource) Read() sdk.ResourceFunc {
	return sdk.ResourceFunc{
		Timeout: 5 * time.Minute,
		Func: func(ctx context.Context, metadata sdk.ResourceMetaData) error {
			client := metadata.Client.ApiManagement.PolicyFragmentClient_v2024_05_01

			id, err := policyfragment.ParseWorkspacePolicyFragmentID(metadata.ResourceData.Id())
			if err != nil {
				return err
			}

			format := metadata.ResourceData.Get("xml_format").(string)
			if format == "" {
				format = string(policyfragment.PolicyFragmentContentFormatXml)
			}
			resp, err := client.WorkspacePolicyFragmentGet(ctx, *id, policyfragment.WorkspacePolicyFragmentGetOperationOptions{
				Format: pointer.ToEnum[policyfragment.PolicyFragmentContentFormat](format),
			})
			if err != nil {
				if response.WasNotFound(resp.HttpResponse) {
					return metadata.MarkAsGone(*id)
				}

				return fmt.Errorf("retrieving %s: %+v", *id, err)
			}

			state := ApiManagementWorkspacePolicyFragmentModel{
				Name:                     id.PolicyFragmentName,
				ApiManagementWorkspaceId: policyfragment.NewWorkspaceID(id.SubscriptionId, id.ResourceGroupName, id.ServiceName, id.WorkspaceId).ID(),
			}

			if model := resp.Model; model != nil {
				if props := model.Properties; props != nil {
					state.Description = pointer.From(props.Description)
					state.XmlContent = props.Value

					// The API only returns `xml_format` when set to "rawxml"; the default "xml" is intentionally never returned.
					format := policyfragment.PolicyFragmentContentFormatXml
					if props.Format != nil {
						format = pointer.From(props.Format)
					}
					state.XmlFormat = string(format)
				}
			}

			return metadata.Encode(&state)
		},
	}
}
