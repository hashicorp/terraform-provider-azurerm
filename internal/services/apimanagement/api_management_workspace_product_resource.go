// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package apimanagement

import (
	"context"
	"fmt"
	"regexp"
	"time"

	"github.com/hashicorp/go-azure-helpers/lang/pointer"
	"github.com/hashicorp/go-azure-helpers/lang/response"
	"github.com/hashicorp/go-azure-helpers/resourcemanager/commonschema"
	"github.com/hashicorp/go-azure-sdk/resource-manager/apimanagement/2024-05-01/product"
	"github.com/hashicorp/go-azure-sdk/resource-manager/apimanagement/2024-05-01/workspace"
	"github.com/hashicorp/terraform-provider-azurerm/internal/sdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/validation"
)

type ApiManagementWorkspaceProductModel struct {
	Name                       string `tfschema:"name"`
	ApiManagementWorkspaceId   string `tfschema:"api_management_workspace_id"`
	RequireApprovalEnabled     bool   `tfschema:"require_approval_enabled"`
	Description                string `tfschema:"description"`
	DisplayName                string `tfschema:"display_name"`
	PublishedEnabled           bool   `tfschema:"published_enabled"`
	RequireSubscriptionEnabled bool   `tfschema:"require_subscription_enabled"`
	SubscriptionsLimit         int64  `tfschema:"subscriptions_limit"`
	Terms                      string `tfschema:"terms"`
}

type ApiManagementWorkspaceProductResource struct{}

var _ sdk.ResourceWithUpdate = ApiManagementWorkspaceProductResource{}

var _ sdk.ResourceWithCustomizeDiff = ApiManagementWorkspaceProductResource{}

func (r ApiManagementWorkspaceProductResource) ResourceType() string {
	return "azurerm_api_management_workspace_product"
}

func (r ApiManagementWorkspaceProductResource) ModelObject() interface{} {
	return &ApiManagementWorkspaceProductModel{}
}

func (r ApiManagementWorkspaceProductResource) IDValidationFunc() pluginsdk.SchemaValidateFunc {
	return product.ValidateWorkspaceProductID
}

func (r ApiManagementWorkspaceProductResource) Arguments() map[string]*pluginsdk.Schema {
	return map[string]*pluginsdk.Schema{
		"name": {
			Type:     pluginsdk.TypeString,
			Required: true,
			ForceNew: true,
			ValidateFunc: validation.StringMatch(
				regexp.MustCompile(`^[a-zA-Z0-9](?:[a-zA-Z0-9-]{0,78}[a-zA-Z0-9])?$`),
				"`name` must be 1–80 characters, using only letters, numbers, or hyphens, and not starting or ending with a hyphen."),
		},

		"api_management_workspace_id": commonschema.ResourceIDReferenceRequiredForceNew(&workspace.WorkspaceId{}),

		"display_name": {
			Type:         pluginsdk.TypeString,
			Required:     true,
			ValidateFunc: validation.StringIsNotEmpty,
		},

		"description": {
			Type:         pluginsdk.TypeString,
			Optional:     true,
			ValidateFunc: validation.StringIsNotEmpty,
		},

		"published_enabled": {
			Type:     pluginsdk.TypeBool,
			Optional: true,
			Default:  false,
		},

		"require_approval_enabled": {
			Type:     pluginsdk.TypeBool,
			Optional: true,
			Default:  false,
		},

		"require_subscription_enabled": {
			Type:     pluginsdk.TypeBool,
			Optional: true,
			Default:  true,
		},

		"subscriptions_limit": {
			Type:         pluginsdk.TypeInt,
			Optional:     true,
			ValidateFunc: validation.IntAtLeast(1),
		},

		"terms": {
			Type:         pluginsdk.TypeString,
			Optional:     true,
			ValidateFunc: validation.StringIsNotEmpty,
		},
	}
}

func (r ApiManagementWorkspaceProductResource) Attributes() map[string]*pluginsdk.Schema {
	return map[string]*pluginsdk.Schema{}
}

func (r ApiManagementWorkspaceProductResource) Create() sdk.ResourceFunc {
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

func (r ApiManagementWorkspaceProductResource) Read() sdk.ResourceFunc {
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

func (r ApiManagementWorkspaceProductResource) Update() sdk.ResourceFunc {
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

			if metadata.ResourceData.HasChange("require_approval_enabled") || metadata.ResourceData.HasChange("subscriptions_limit") {
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

func (r ApiManagementWorkspaceProductResource) Delete() sdk.ResourceFunc {
	return sdk.ResourceFunc{
		Timeout: 30 * time.Minute,
		Func: func(ctx context.Context, metadata sdk.ResourceMetaData) error {
			client := metadata.Client.ApiManagement.ProductClient_v2024_05_01

			id, err := product.ParseWorkspaceProductID(metadata.ResourceData.Id())
			if err != nil {
				return err
			}

			if resp, err := client.WorkspaceProductDelete(ctx, *id, product.WorkspaceProductDeleteOperationOptions{DeleteSubscriptions: pointer.To(true)}); err != nil {
				if !response.WasNotFound(resp.HttpResponse) {
					return fmt.Errorf("deleting %s: %+v", *id, err)
				}
			}

			return nil
		},
	}
}

func (r ApiManagementWorkspaceProductResource) CustomizeDiff() sdk.ResourceFunc {
	return sdk.ResourceFunc{
		Func: func(ctx context.Context, metadata sdk.ResourceMetaData) error {
			if !metadata.ResourceDiff.Get("require_subscription_enabled").(bool) {
				if !metadata.ResourceDiff.GetRawConfig().AsValueMap()["subscriptions_limit"].IsNull() {
					return fmt.Errorf("`require_subscription_enabled` must be set to `true` when `subscriptions_limit` is specified")
				}

				if !metadata.ResourceDiff.GetRawConfig().AsValueMap()["require_approval_enabled"].IsNull() {
					return fmt.Errorf("`require_subscription_enabled` must be set to `true` when `require_approval_enabled` is specified")
				}
			}

			return nil
		},
		Timeout: 30 * time.Minute,
	}
}
