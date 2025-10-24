// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package graphservices

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
	"github.com/hashicorp/terraform-provider-azurerm/internal/sdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
)

var _ sdk.DataSource = AccountDataSource{}

type AccountDataSource struct{}

type AccountDataSourceModel struct {
	ApplicationId     string                 `tfschema:"application_id"`
	BillingPlanId     string                 `tfschema:"billing_plan_id"`
	Name              string                 `tfschema:"name"`
	ResourceGroupName string                 `tfschema:"resource_group_name"`
	Tags              map[string]interface{} `tfschema:"tags"`
}

func (r AccountDataSource) ModelObject() interface{} {
	return &AccountDataSourceModel{}
}

func (r AccountDataSource) ResourceType() string {
	return "azurerm_graph_services_account"
}

func (r AccountDataSource) Arguments() map[string]*pluginsdk.Schema {
	return map[string]*pluginsdk.Schema{
		"name": {
			Type:         pluginsdk.TypeString,
			Required:     true,
			ValidateFunc: validation.StringIsNotEmpty,
		},

		"resource_group_name": commonschema.ResourceGroupNameForDataSource(),
	}
}

func (r AccountDataSource) Attributes() map[string]*pluginsdk.Schema {
	return map[string]*pluginsdk.Schema{
		"application_id": {
			Type:     pluginsdk.TypeString,
			Computed: true,
		},

		"billing_plan_id": {
			Type:     pluginsdk.TypeString,
			Computed: true,
		},

		"tags": commonschema.TagsDataSource(),
	}
}

func (r AccountDataSource) Read() sdk.ResourceFunc {
	return sdk.ResourceFunc{
		Timeout: 5 * time.Minute,
		Func: func(ctx context.Context, metadata sdk.ResourceMetaData) error {
			client := metadata.Client.Graph.V20230413.Graphservicesprods
			subscriptionId := metadata.Client.Account.SubscriptionId

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
