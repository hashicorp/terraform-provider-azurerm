// Copyright IBM Corp. 2014, 2026
// SPDX-License-Identifier: MPL-2.0

package api_management_workspace_product

import (
	"context"
	"fmt"
	"time"

	"github.com/hashicorp/go-azure-helpers/lang/pointer"
	"github.com/hashicorp/go-azure-helpers/lang/response"
	"github.com/hashicorp/go-azure-sdk/resource-manager/apimanagement/2024-05-01/product"
	"github.com/hashicorp/go-azure-sdk/resource-manager/apimanagement/2024-05-01/workspace"
	"github.com/hashicorp/terraform-provider-azurerm/internal/sdk"
)

func (r Resource) Read() sdk.ResourceFunc {
	return sdk.ResourceFunc{
		Timeout: 5 * time.Minute,
		Func: func(ctx context.Context, metadata sdk.ResourceMetaData) error {
			client := metadata.Client.ApiManagement.ProductClient_v2024_05_01

			id, err := product.ParseWorkspaceProductID(metadata.ResourceData.Id())
			if err != nil {
				return err
			}

			resp, err := client.WorkspaceProductGet(ctx, *id)
			if err != nil {
				if response.WasNotFound(resp.HttpResponse) {
					return metadata.MarkAsGone(*id)
				}
				return fmt.Errorf("retrieving %s: %+v", *id, err)
			}

			state := ApiManagementWorkspaceProductModel{
				Name:                     id.ProductId,
				ApiManagementWorkspaceId: workspace.NewWorkspaceID(id.SubscriptionId, id.ResourceGroupName, id.ServiceName, id.WorkspaceId).ID(),
			}

			if model := resp.Model; model != nil {
				if props := model.Properties; props != nil {
					state.RequireApprovalEnabled = pointer.From(props.ApprovalRequired)
					state.Description = pointer.From(props.Description)
					state.DisplayName = props.DisplayName
					state.PublishedEnabled = pointer.From(props.State) == product.ProductStatePublished
					state.RequireSubscriptionEnabled = pointer.From(props.SubscriptionRequired)
					state.SubscriptionsLimit = pointer.From(props.SubscriptionsLimit)
					state.Terms = pointer.From(props.Terms)
				}
			}

			return metadata.Encode(&state)
		},
	}
}
