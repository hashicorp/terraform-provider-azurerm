// Copyright IBM Corp. 2014, 2026
// SPDX-License-Identifier: MPL-2.0

package graph_services_account

import (
	"context"
	"fmt"
	"time"

	"github.com/hashicorp/go-azure-helpers/lang/pointer"
	"github.com/hashicorp/go-azure-helpers/lang/response"
	"github.com/hashicorp/go-azure-helpers/resourcemanager/tags"
	"github.com/hashicorp/go-azure-sdk/resource-manager/graphservices/2023-04-13/graphservicesprods"
	"github.com/hashicorp/terraform-provider-azurerm/internal/sdk"
)

func (r Resource) Read() sdk.ResourceFunc {
	return sdk.ResourceFunc{
		Timeout: 5 * time.Minute,
		Func: func(ctx context.Context, metadata sdk.ResourceMetaData) error {
			client := metadata.Client.Graph.V20230413.Graphservicesprods

			id, err := graphservicesprods.ParseAccountID(metadata.ResourceData.Id())
			if err != nil {
				return err
			}

			resp, err := client.AccountsGet(ctx, *id)
			if err != nil {
				if response.WasNotFound(resp.HttpResponse) {
					return metadata.MarkAsGone(*id)
				}
				return fmt.Errorf("retrieving %s: %+v", *id, err)
			}

			model := resp.Model
			if model == nil {
				return fmt.Errorf("retrieving %s: model was nil", *id)
			}

			schema := AccountResourceSchema{
				ApplicationId:     model.Properties.AppId,
				Name:              id.AccountName,
				ResourceGroupName: id.ResourceGroupName,
				Tags:              tags.Flatten(model.Tags),
			}
			if model.Properties.BillingPlanId != nil {
				schema.BillingPlanId = pointer.From(model.Properties.BillingPlanId)
			}

			return metadata.Encode(&schema)
		},
	}
}

type AccountResourceSchema struct {
	ApplicationId     string         `tfschema:"application_id"`
	BillingPlanId     string         `tfschema:"billing_plan_id"`
	Name              string         `tfschema:"name"`
	ResourceGroupName string         `tfschema:"resource_group_name"`
	Tags              map[string]any `tfschema:"tags"`
}
