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

func (r Resource) Create() sdk.ResourceFunc {
	return sdk.ResourceFunc{
		Timeout: 30 * time.Minute,
		Func: func(ctx context.Context, metadata sdk.ResourceMetaData) error {
			client := metadata.Client.ApiManagement.ProductClient_v2024_05_01

			var model ApiManagementWorkspaceProductModel
			if err := metadata.Decode(&model); err != nil {
				return fmt.Errorf("decoding: %+v", err)
			}

			workspaceId, err := workspace.ParseWorkspaceID(model.ApiManagementWorkspaceId)
			if err != nil {
				return err
			}

			id := product.NewWorkspaceProductID(workspaceId.SubscriptionId, workspaceId.ResourceGroupName, workspaceId.ServiceName, workspaceId.WorkspaceId, model.Name)
			existing, err := client.WorkspaceProductGet(ctx, id)
			if err != nil && !response.WasNotFound(existing.HttpResponse) {
				return fmt.Errorf("checking for presence of existing %s: %+v", id, err)
			}

			if !response.WasNotFound(existing.HttpResponse) {
				return metadata.ResourceRequiresImport(r.ResourceType(), id)
			}

			publishedState := product.ProductStateNotPublished
			if model.PublishedEnabled {
				publishedState = product.ProductStatePublished
			}

			properties := product.ProductContract{
				Properties: &product.ProductContractProperties{
					Description:          pointer.To(model.Description),
					DisplayName:          model.DisplayName,
					State:                pointer.To(publishedState),
					SubscriptionRequired: pointer.To(model.RequireSubscriptionEnabled),
				},
			}

			// Cannot provide values for `require_approval_enabled` and `subscriptions_limit` when `require_subscription_enabled` is set to false in the request payload
			if model.RequireSubscriptionEnabled {
				if model.SubscriptionsLimit > 0 {
					properties.Properties.SubscriptionsLimit = pointer.To(model.SubscriptionsLimit)
				}
				properties.Properties.ApprovalRequired = pointer.To(model.RequireApprovalEnabled)
			}

			if model.Terms != "" {
				properties.Properties.Terms = pointer.To(model.Terms)
			}
			if _, err := client.WorkspaceProductCreateOrUpdate(ctx, id, properties, product.WorkspaceProductCreateOrUpdateOperationOptions{}); err != nil {
				return fmt.Errorf("creating %s: %+v", id, err)
			}

			metadata.SetID(id)
			return nil
		},
	}
}
