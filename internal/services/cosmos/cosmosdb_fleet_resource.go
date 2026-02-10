// Copyright IBM Corp. 2014, 2025
// SPDX-License-Identifier: MPL-2.0

package cosmos

import (
	"context"
	"fmt"
	"regexp"
	"time"

	"github.com/hashicorp/go-azure-helpers/lang/pointer"
	"github.com/hashicorp/go-azure-helpers/lang/response"
	"github.com/hashicorp/go-azure-helpers/resourcemanager/commonschema"
	"github.com/hashicorp/go-azure-helpers/resourcemanager/location"
	"github.com/hashicorp/go-azure-helpers/resourcemanager/resourceids"
	"github.com/hashicorp/go-azure-helpers/resourcemanager/tags"
	"github.com/hashicorp/go-azure-sdk/resource-manager/cosmosdb/2025-10-15/fleets"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-provider-azurerm/internal/sdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/validation"
)

//go:generate go run ../../tools/generator-tests resourceidentity -resource-name cosmosdb_fleet -service-package-name cosmos -properties "name,resource_group_name" -known-values "subscription_id:data.Subscriptions.Primary"

type CosmosDbFleetResource struct{}

var _ sdk.ResourceWithIdentity = CosmosDbFleetResource{}

type CosmosDbFleetModel struct {
	Name              string            `tfschema:"name"`
	ResourceGroupName string            `tfschema:"resource_group_name"`
	Location          string            `tfschema:"location"`
	Tags              map[string]string `tfschema:"tags"`
}

func (CosmosDbFleetResource) Arguments() map[string]*pluginsdk.Schema {
	return map[string]*pluginsdk.Schema{
		"name": {
			Type:     pluginsdk.TypeString,
			Required: true,
			ForceNew: true,
			ValidateFunc: validation.StringMatch(
				regexp.MustCompile("^(?=.{3,50}$)[a-z0-9]+(?:-[a-z0-9]+)*$"),
				"CosmosDB Fleet name must be 3 - 50 characters long, contain only lowercase letters, numbers and hyphens.",
			),
		},

		"resource_group_name": commonschema.ResourceGroupName(),

		"location": commonschema.Location(),

		"tags": {
			Type:     schema.TypeMap,
			Optional: true,
			// `ForceNew` behavior is added as `tags` property is absent in update model
			ForceNew:     true,
			ValidateFunc: tags.Validate,
			Elem: &schema.Schema{
				Type: schema.TypeString,
			},
		},
	}
}

func (CosmosDbFleetResource) Attributes() map[string]*pluginsdk.Schema {
	return map[string]*pluginsdk.Schema{}
}

func (CosmosDbFleetResource) ModelObject() interface{} {
	return &CosmosDbFleetResource{}
}

func (CosmosDbFleetResource) ResourceType() string {
	return "azurerm_cosmosdb_fleet"
}

func (r CosmosDbFleetResource) Create() sdk.ResourceFunc {
	return sdk.ResourceFunc{
		Timeout: 30 * time.Minute,
		Func: func(ctx context.Context, metadata sdk.ResourceMetaData) error {
			client := metadata.Client.Cosmos.FleetsClient
			subscriptionId := metadata.Client.Account.SubscriptionId

			var config CosmosDbFleetModel
			if err := metadata.Decode(&config); err != nil {
				return fmt.Errorf("decoding: %+v", err)
			}
			id := fleets.NewFleetID(subscriptionId, config.ResourceGroupName, config.Name)

			existing, err := client.FleetGet(ctx, id)
			if err != nil && !response.WasNotFound(existing.HttpResponse) {
				return fmt.Errorf("checking for presence of existing %s: %+v", id, err)
			}
			if !response.WasNotFound(existing.HttpResponse) {
				return metadata.ResourceRequiresImport(r.ResourceType(), id)
			}

			param := fleets.FleetResource{
				Location: location.Normalize(config.Location),
				Tags:     pointer.To(config.Tags),
			}
			if _, err := client.FleetCreate(ctx, id, param); err != nil {
				return fmt.Errorf("creating %s: %+v", id, err)
			}

			metadata.SetID(id)
			if err := pluginsdk.SetResourceIdentityData(metadata.ResourceData, &id); err != nil {
				return err
			}

			return nil
		},
	}
}

func (r CosmosDbFleetResource) Read() sdk.ResourceFunc {
	return sdk.ResourceFunc{
		Timeout: 5 * time.Minute,
		Func: func(ctx context.Context, metadata sdk.ResourceMetaData) error {
			client := metadata.Client.Cosmos.FleetsClient
			id, err := fleets.ParseFleetID(metadata.ResourceData.Id())
			if err != nil {
				return err
			}

			resp, err := client.FleetGet(ctx, *id)
			if err != nil {
				if response.WasNotFound(resp.HttpResponse) {
					return metadata.MarkAsGone(id)
				}

				return fmt.Errorf("retrieving %s: %+v", id, err)
			}

			state := CosmosDbFleetModel{
				Name:              id.FleetName,
				ResourceGroupName: id.ResourceGroupName,
			}

			if model := resp.Model; model != nil {
				state.Location = location.NormalizeNilable(&model.Location)
				state.Tags = pointer.From(model.Tags)
			}

			if err := pluginsdk.SetResourceIdentityData(metadata.ResourceData, id); err != nil {
				return err
			}

			return metadata.Encode(&state)
		},
	}
}

func (CosmosDbFleetResource) Delete() sdk.ResourceFunc {
	return sdk.ResourceFunc{
		Timeout: 30 * time.Minute,
		Func: func(ctx context.Context, metadata sdk.ResourceMetaData) error {
			client := metadata.Client.Cosmos.FleetsClient

			id, err := fleets.ParseFleetID(metadata.ResourceData.Id())
			if err != nil {
				return err
			}

			if err := client.FleetDeleteThenPoll(ctx, *id); err != nil {
				return fmt.Errorf("deleting %s: %+v", *id, err)
			}
			return nil
		},
	}
}

func (CosmosDbFleetResource) IDValidationFunc() pluginsdk.SchemaValidateFunc {
	return fleets.ValidateFleetID
}

func (CosmosDbFleetResource) Identity() resourceids.ResourceId {
	return &fleets.FleetId{}
}
