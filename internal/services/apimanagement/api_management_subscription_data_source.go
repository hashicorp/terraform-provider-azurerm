package apimanagement

import (
	"context"
	"fmt"
	"strings"
	"time"

	"github.com/hashicorp/go-azure-helpers/lang/pointer"
	"github.com/hashicorp/go-azure-helpers/lang/response"
	"github.com/hashicorp/go-azure-helpers/resourcemanager/commonschema"
	"github.com/hashicorp/go-azure-sdk/resource-manager/apimanagement/2022-08-01/api"
	"github.com/hashicorp/go-azure-sdk/resource-manager/apimanagement/2022-08-01/product"
	"github.com/hashicorp/go-azure-sdk/resource-manager/apimanagement/2022-08-01/subscription"
	"github.com/hashicorp/terraform-provider-azurerm/internal/sdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/services/apimanagement/schemaz"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
)

var _ sdk.DataSource = ApiManagementSubscriptionDataSource{}

type ApiManagementSubscriptionDataSource struct{}

type ApiManagementSubscriptionDataSourceModel struct {
	ApiManagementName string `tfschema:"api_management_name"`
	ResourceGroupName string `tfschema:"resource_group_name"`
	SubscriptionId    string `tfschema:"subscription_id"`
	AllowTracing      bool   `tfschema:"allow_tracing"`
	ApiId             string `tfschema:"api_id"`
	DisplayName       string `tfschema:"display_name"`
	PrimaryKey        string `tfschema:"primary_key"`
	ProductId         string `tfschema:"product_id"`
	SecondaryKey      string `tfschema:"secondary_key"`
	State             string `tfschema:"state"`
	UserId            string `tfschema:"user_id"`
}

func (ApiManagementSubscriptionDataSource) Arguments() map[string]*pluginsdk.Schema {
	return map[string]*pluginsdk.Schema{
		"api_management_name": schemaz.SchemaApiManagementDataSourceName(),

		"resource_group_name": commonschema.ResourceGroupName(),

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
			subscriptionId := metadata.Client.Account.SubscriptionId

			var state ApiManagementSubscriptionDataSourceModel

			if err := metadata.Decode(&state); err != nil {
				return fmt.Errorf("decoding: %+v", err)
			}

			id := subscription.NewSubscriptions2ID(subscriptionId, state.ResourceGroupName, state.ApiManagementName, state.SubscriptionId)

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
					productId := ""
					apiId := ""
					// check if the subscription is for all apis or a specific product/ api
					if props.Scope != "" && !strings.HasSuffix(props.Scope, "/apis") {
						// the scope is either a product or api id
						parseId, err := product.ParseProductIDInsensitively(props.Scope)
						if err == nil {
							productId = parseId.ID()
						} else {
							parsedApiId, err := api.ParseApiIDInsensitively(props.Scope)
							if err != nil {
								return fmt.Errorf("parsing scope into product/ api id %q: %+v", props.Scope, err)
							}
							apiId = parsedApiId.ID()
						}
					}
					state.AllowTracing = pointer.From(props.AllowTracing)
					state.ApiId = apiId
					state.DisplayName = pointer.From(props.DisplayName)
					state.ProductId = productId
					state.State = string(props.State)
					state.UserId = pointer.From(props.OwnerId)
				}
			}

			// Primary and secondary keys must be got from this additional api
			keyResp, err := client.ListSecrets(ctx, id)
			if err != nil {
				return fmt.Errorf("listing Subscription %q Primary and Secondary Keys (API Management Service %q / Resource Group %q): %+v", id.SubscriptionId, id.ServiceName, id.ResourceGroupName, err)
			}
			if model := keyResp.Model; model != nil {
				state.SecondaryKey = pointer.From(model.SecondaryKey)
				state.PrimaryKey = pointer.From(model.PrimaryKey)
			}

			return metadata.Encode(&state)
		},
	}
}
