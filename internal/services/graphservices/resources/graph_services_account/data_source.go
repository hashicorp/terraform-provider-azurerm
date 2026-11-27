// Copyright IBM Corp. 2014, 2026
// SPDX-License-Identifier: MPL-2.0

package graph_services_account

import (
	"context"
	"fmt"
	"time"

	"github.com/hashicorp/go-azure-helpers/lang/pointer"
	"github.com/hashicorp/go-azure-helpers/lang/response"
	"github.com/hashicorp/go-azure-helpers/resourcemanager/commonschema"
	"github.com/hashicorp/go-azure-helpers/resourcemanager/tags"
	"github.com/hashicorp/go-azure-sdk/resource-manager/graphservices/2023-04-13/graphservicesprods"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"
	"github.com/hashicorp/terraform-provider-azurerm/internal/clients/account"
	"github.com/hashicorp/terraform-provider-azurerm/internal/sdk"
	graphservicesClient "github.com/hashicorp/terraform-provider-azurerm/internal/services/graphservices/client"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
)

var _ sdk.DataSource = DataSource{}

type DataSource struct{}

type AccountDataSourceModel struct {
	ApplicationId     string         `tfschema:"application_id"`
	BillingPlanId     string         `tfschema:"billing_plan_id"`
	Name              string         `tfschema:"name"`
	ResourceGroupName string         `tfschema:"resource_group_name"`
	Tags              map[string]any `tfschema:"tags"`
}

func (r DataSource) ModelObject() any {
	return &AccountDataSourceModel{}
}

func (r DataSource) ResourceType() string {
	return "azurerm_graph_services_account"
}

func (r DataSource) Description() string {
	return "Use this data source to access information about an existing Graph Services Account."
}

func (r DataSource) Arguments() map[string]*pluginsdk.Schema {
	return map[string]*pluginsdk.Schema{
		"name": {
			Description:  "The name of the Graph Services Account.",
			Type:         pluginsdk.TypeString,
			Required:     true,
			ValidateFunc: validation.StringIsNotEmpty,
		},

		"resource_group_name": commonschema.ResourceGroupNameForDataSource(),
	}
}

func (r DataSource) Attributes() map[string]*pluginsdk.Schema {
	return map[string]*pluginsdk.Schema{
		"application_id": {
			Description: "The application ID associated with the Graph Services Account.",
			Type:        pluginsdk.TypeString,
			Computed:    true,
		},

		"billing_plan_id": {
			Description: "The billing plan ID for the Graph Services Account.",
			Type:        pluginsdk.TypeString,
			Computed:    true,
		},

		"tags": commonschema.TagsDataSource(),
	}
}

func (r DataSource) Read() sdk.ResourceFunc {
	return sdk.ResourceFunc{
		Timeout: 5 * time.Minute,
		Func: func(ctx context.Context, metadata sdk.ResourceMetaData) error {
			client := graphservicesClient.From(metadata.Client).V20230413.Graphservicesprods
			subscriptionId := account.From(metadata.Client).SubscriptionId

			var state AccountDataSourceModel
			if err := metadata.Decode(&state); err != nil {
				return fmt.Errorf("decoding: %+v", err)
			}

			id := graphservicesprods.NewAccountID(subscriptionId, state.ResourceGroupName, state.Name)

			resp, err := client.AccountsGet(ctx, id)
			if err != nil {
				if response.WasNotFound(resp.HttpResponse) {
					return fmt.Errorf("%s was not found", id)
				}
				return fmt.Errorf("retrieving %s: %+v", id, err)
			}

			metadata.SetID(id)

			if model := resp.Model; model != nil {
				state.ApplicationId = model.Properties.AppId
				state.Name = id.AccountName
				state.ResourceGroupName = id.ResourceGroupName
				state.Tags = tags.Flatten(model.Tags)

				if model.Properties.BillingPlanId != nil {
					state.BillingPlanId = pointer.From(model.Properties.BillingPlanId)
				}
			}

			return metadata.Encode(&state)
		},
	}
}
