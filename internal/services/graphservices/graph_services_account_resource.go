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
	"github.com/hashicorp/go-azure-helpers/resourcemanager/location"
	"github.com/hashicorp/go-azure-helpers/resourcemanager/tags"
	"github.com/hashicorp/go-azure-sdk/resource-manager/graphservices/2023-04-13/graphservicesprods"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"
	"github.com/hashicorp/terraform-provider-azurerm/internal/sdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
)

var (
	_ sdk.Resource           = AccountResource{}
	_ sdk.ResourceWithUpdate = AccountResource{}
)

type AccountResource struct{}

func (r AccountResource) ModelObject() interface{} {
	return &AccountResourceSchema{}
}

type AccountResourceSchema struct {
	ApplicationId     string                 `tfschema:"application_id"`
	BillingPlanId     string                 `tfschema:"billing_plan_id"`
	Name              string                 `tfschema:"name"`
	ResourceGroupName string                 `tfschema:"resource_group_name"`
	Tags              map[string]interface{} `tfschema:"tags"`
}

func (r AccountResource) IDValidationFunc() pluginsdk.SchemaValidateFunc {
	return graphservicesprods.ValidateAccountID
}

func (r AccountResource) ResourceType() string {
	return "azurerm_graph_services_account"
}

func (r AccountResource) Arguments() map[string]*pluginsdk.Schema {
	return map[string]*pluginsdk.Schema{
		"name": {
			ForceNew:     true,
			Required:     true,
			Type:         pluginsdk.TypeString,
			ValidateFunc: validation.StringIsNotEmpty,
		},
		"resource_group_name": commonschema.ResourceGroupName(),
		"application_id": {
			ForceNew:     true,
			Required:     true,
			Type:         pluginsdk.TypeString,
			ValidateFunc: validation.IsUUID,
		},
		"tags": commonschema.Tags(),
	}
}

func (r AccountResource) Attributes() map[string]*pluginsdk.Schema {
	return map[string]*pluginsdk.Schema{
		"billing_plan_id": {
			Computed: true,
			Type:     pluginsdk.TypeString,
		},
	}
}

func (r AccountResource) Create() sdk.ResourceFunc {
	return sdk.ResourceFunc{
		Timeout: 30 * time.Minute,
		Func: func(ctx context.Context, metadata sdk.ResourceMetaData) error {
			client := metadata.Client.Graph.V20230413.Graphservicesprods

			var config AccountResourceSchema
			if err := metadata.Decode(&config); err != nil {
				return fmt.Errorf("decoding: %+v", err)
			}

			subscriptionId := metadata.Client.Account.SubscriptionId
			id := graphservicesprods.NewAccountID(subscriptionId, config.ResourceGroupName, config.Name)

			existing, err := client.AccountsGet(ctx, id)
			if err != nil {
				if !response.WasNotFound(existing.HttpResponse) {
					return fmt.Errorf("checking for the presence of an existing %s: %+v", id, err)
				}
			}
			if !response.WasNotFound(existing.HttpResponse) {
				return metadata.ResourceRequiresImport("azurerm_graph_services_account", id)
			}

			payload := graphservicesprods.AccountResource{
				Location: pointer.To(location.Normalize("global")),
				Tags:     tags.Expand(config.Tags),
				Properties: graphservicesprods.AccountResourceProperties{
					AppId: config.ApplicationId,
				},
			}

			if err := client.AccountsCreateAndUpdateThenPoll(ctx, id, payload); err != nil {
				return fmt.Errorf("creating %s: %+v", id, err)
			}

			metadata.SetID(id)
			return nil
		},
	}
}

func (r AccountResource) Read() sdk.ResourceFunc {
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

func (r AccountResource) Delete() sdk.ResourceFunc {
	return sdk.ResourceFunc{
		Timeout: 30 * time.Minute,
		Func: func(ctx context.Context, metadata sdk.ResourceMetaData) error {
			client := metadata.Client.Graph.V20230413.Graphservicesprods

			id, err := graphservicesprods.ParseAccountID(metadata.ResourceData.Id())
			if err != nil {
				return err
			}

			if _, err := client.AccountsDelete(ctx, *id); err != nil {
				return fmt.Errorf("deleting %s: %+v", *id, err)
			}

			return nil
		},
	}
}

func (r AccountResource) Update() sdk.ResourceFunc {
	return sdk.ResourceFunc{
		Timeout: 30 * time.Minute,
		Func: func(ctx context.Context, metadata sdk.ResourceMetaData) error {
			client := metadata.Client.Graph.V20230413.Graphservicesprods

			id, err := graphservicesprods.ParseAccountID(metadata.ResourceData.Id())
			if err != nil {
				return err
			}

			var config AccountResourceSchema
			if err := metadata.Decode(&config); err != nil {
				return fmt.Errorf("decoding: %+v", err)
			}

			existing, err := client.AccountsGet(ctx, *id)
			if err != nil {
				return fmt.Errorf("retrieving existing %s: %+v", *id, err)
			}
			if existing.Model == nil {
				return fmt.Errorf("retrieving existing %s: model was nil", *id)
			}

			payload := graphservicesprods.TagUpdate{}
			if metadata.ResourceData.HasChange("tags") {
				payload.Tags = tags.Expand(config.Tags)
			}

			if _, err := client.AccountsUpdate(ctx, *id, payload); err != nil {
				return fmt.Errorf("updating %s: %+v", *id, err)
			}

			return nil
		},
	}
}
