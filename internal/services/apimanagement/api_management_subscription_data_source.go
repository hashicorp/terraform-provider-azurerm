// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package apimanagement

import (
	"context"
	"fmt"
	"strings"
	"time"

	"github.com/hashicorp/go-azure-helpers/lang/pointer"
	"github.com/hashicorp/go-azure-helpers/lang/response"
	"github.com/hashicorp/go-azure-sdk/resource-manager/apimanagement/2022-08-01/api"
	"github.com/hashicorp/go-azure-sdk/resource-manager/apimanagement/2022-08-01/product"
	"github.com/hashicorp/go-azure-sdk/resource-manager/apimanagement/2022-08-01/subscription"
	"github.com/hashicorp/go-azure-sdk/resource-manager/apimanagement/2024-05-01/apimanagementservice"
	"github.com/hashicorp/terraform-provider-azurerm/internal/sdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/services/apimanagement/schemaz"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
)

var _ sdk.DataSource = ApiManagementSubscriptionDataSource{}

type ApiManagementSubscriptionDataSource struct{}

type ApiManagementSubscriptionDataSourceModel struct {
	ApiManagementId string `tfschema:"api_management_id"`
	SubscriptionId  string `tfschema:"subscription_id"`
	AllowTracing    bool   `tfschema:"allow_tracing"`
	ApiId           string `tfschema:"api_id"`
	DisplayName     string `tfschema:"display_name"`
	PrimaryKey      string `tfschema:"primary_key"`
	ProductId       string `tfschema:"product_id"`
	SecondaryKey    string `tfschema:"secondary_key"`
	State           string `tfschema:"state"`
	UserId          string `tfschema:"user_id"`
}

func (ApiManagementSubscriptionDataSource) Arguments() map[string]*pluginsdk.Schema {
	return map[string]*pluginsdk.Schema{
		"api_management_id": {
			Type:         pluginsdk.TypeString,
			Required:     true,
			ForceNew:     true,
			ValidateFunc: apimanagementservice.ValidateServiceID,
		},

		"subscription_id": schemaz.SchemaApiManagementChildDataSourceName(),
	}
}

func (ApiManagementSubscriptionDataSource) Attributes() map[string]*pluginsdk.Schema {
	return map[string]*pluginsdk.Schema{
		"allow_tracing": {
			Type:     pluginsdk.TypeBool,
			Computed: true,
		},

		"api_id": {
			Type:     pluginsdk.TypeString,
			Computed: true,
		},

		"display_name": {
			Type:     pluginsdk.TypeString,
			Computed: true,
		},

		"primary_key": {
			Type:      pluginsdk.TypeString,
			Computed:  true,
			Sensitive: true,
		},

		"product_id": {
			Type:     pluginsdk.TypeString,
			Computed: true,
		},

		"secondary_key": {
			Type:      pluginsdk.TypeString,
			Computed:  true,
			Sensitive: true,
		},

		"state": {
			Type:     pluginsdk.TypeString,
			Computed: true,
		},
		"user_id": {
			Type:     pluginsdk.TypeString,
			Computed: true,
		},
	}
}

func (ApiManagementSubscriptionDataSource) ModelObject() interface{} {
	return &ApiManagementSubscriptionDataSourceModel{}
}

func (ApiManagementSubscriptionDataSource) ResourceType() string {
	return "azurerm_api_management_subscription"
}

func (ApiManagementSubscriptionDataSource) Read() sdk.ResourceFunc {
	return sdk.ResourceFunc{
		Timeout: 5 * time.Minute,
		Func: func(ctx context.Context, metadata sdk.ResourceMetaData) error {
			client := metadata.Client.ApiManagement.SubscriptionsClient

			var state ApiManagementSubscriptionDataSourceModel

			if err := metadata.Decode(&state); err != nil {
				return fmt.Errorf("decoding: %+v", err)
			}

			api_management_id, err := apimanagementservice.ParseServiceID(state.ApiManagementId)
			if err != nil {
				return fmt.Errorf("parsing: %+v", err)
			}

			id := subscription.NewSubscriptions2ID(api_management_id.SubscriptionId, api_management_id.ResourceGroupName, api_management_id.ServiceName, state.SubscriptionId)

			resp, err := client.Get(ctx, id)
			if err != nil {
				if response.WasNotFound(resp.HttpResponse) {
					return fmt.Errorf("%s was not found", id)
				}

				return fmt.Errorf("retrieving %s: %+v", id, err)
			}

			metadata.SetID(id)

			if model := resp.Model; model != nil {
				if props := model.Properties; props != nil {
					// Check if the subscription is for all apis or a specific product/ api, excluding the built-in master subscription.
					// The scope of the built-in subscription is the API Management service itself (service ID).
					if props.Scope != "" && !strings.HasSuffix(props.Scope, "/apis") {
						// the scope is either a product, api id or service id
						if productId, err := product.ParseProductIDInsensitively(props.Scope); err == nil {
							state.ProductId = productId.ID()
						}
						if apiId, err := api.ParseApiIDInsensitively(props.Scope); err == nil {
							state.ApiId = apiId.ID()
						}
					}
					state.AllowTracing = pointer.From(props.AllowTracing)
					state.DisplayName = pointer.From(props.DisplayName)
					state.State = string(props.State)
					state.UserId = pointer.From(props.OwnerId)
				}
			}

			// Primary and secondary keys must be got from this additional api
			keyResp, err := client.ListSecrets(ctx, id)
			if err != nil {
				return fmt.Errorf("listing Primary and Secondary Keys for %s: %+v", id, err)
			}
			if model := keyResp.Model; model != nil {
				state.SecondaryKey = pointer.From(model.SecondaryKey)
				state.PrimaryKey = pointer.From(model.PrimaryKey)
			}

			return metadata.Encode(&state)
		},
	}
}

// Temporary comment to ensure PR has changes

// Temporary comment to ensure PR has changes
