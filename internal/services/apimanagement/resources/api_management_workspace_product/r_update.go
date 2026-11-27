// Copyright IBM Corp. 2014, 2026
// SPDX-License-Identifier: MPL-2.0

package api_management_workspace_product

import (
	"context"
	"fmt"
	"time"

	"github.com/hashicorp/go-azure-helpers/lang/pointer"
	"github.com/hashicorp/go-azure-sdk/resource-manager/apimanagement/2024-05-01/product"
	"github.com/hashicorp/terraform-provider-azurerm/internal/sdk"
)

func (r Resource) Update() sdk.ResourceFunc {
	return sdk.ResourceFunc{
		Timeout: 30 * time.Minute,
		Func: func(ctx context.Context, metadata sdk.ResourceMetaData) error {
			client := metadata.Client.ApiManagement.ProductClient_v2024_05_01

			id, err := product.ParseWorkspaceProductID(metadata.ResourceData.Id())
			if err != nil {
				return err
			}

			var model ApiManagementWorkspaceProductModel
			if err := metadata.Decode(&model); err != nil {
				return fmt.Errorf("decoding: %+v", err)
			}

			existing, err := client.WorkspaceProductGet(ctx, *id)
			if err != nil {
				return fmt.Errorf("retrieving %s: %+v", *id, err)
			}

			if existing.Model == nil {
				return fmt.Errorf("retrieving %s: `model` was nil", *id)
			}

			if existing.Model.Properties == nil {
				return fmt.Errorf("retrieving %s: `properties` was nil", *id)
			}

			properties := existing.Model
			if metadata.ResourceData.HasChange("display_name") {
				properties.Properties.DisplayName = model.DisplayName
			}

			if metadata.ResourceData.HasChange("description") {
				properties.Properties.Description = pointer.To(model.Description)
			}

			if metadata.ResourceData.HasChange("published_enabled") {
				publishedState := product.ProductStateNotPublished
				if model.PublishedEnabled {
					publishedState = product.ProductStatePublished
				}
				properties.Properties.State = pointer.To(publishedState)
			}

			if metadata.ResourceData.HasChange("require_subscription_enabled") {
				properties.Properties.SubscriptionRequired = pointer.To(model.RequireSubscriptionEnabled)
			}

			if metadata.ResourceData.HasChange("terms") {
				properties.Properties.Terms = pointer.To(model.Terms)
			}

			if metadata.ResourceData.HasChanges("require_approval_enabled", "subscriptions_limit") {
				if model.RequireSubscriptionEnabled {
					if model.SubscriptionsLimit > 0 {
						properties.Properties.SubscriptionsLimit = pointer.To(model.SubscriptionsLimit)
					} else {
						properties.Properties.SubscriptionsLimit = nil
					}
					properties.Properties.ApprovalRequired = pointer.To(model.RequireApprovalEnabled)
				} else {
					// Cannot provide values for `require_approval_enabled` and `subscriptions_limit` when `require_subscription_enabled` is set to false in the request payload
					properties.Properties.ApprovalRequired = nil
					properties.Properties.SubscriptionsLimit = nil
				}
			}

			if _, err := client.WorkspaceProductCreateOrUpdate(ctx, *id, *properties, product.WorkspaceProductCreateOrUpdateOperationOptions{}); err != nil {
				return fmt.Errorf("updating %s: %+v", *id, err)
			}

			return nil
		},
	}
}
