// Copyright IBM Corp. 2014, 2025
// SPDX-License-Identifier: MPL-2.0

package cosmos

import (
	"context"
	"errors"
	"fmt"
	"time"

	"github.com/hashicorp/go-azure-helpers/lang/pointer"
	"github.com/hashicorp/go-azure-helpers/lang/response"
	"github.com/hashicorp/go-azure-helpers/resourcemanager/commonschema"
	"github.com/hashicorp/go-azure-helpers/resourcemanager/location"
	"github.com/hashicorp/go-azure-helpers/resourcemanager/resourceids"
	"github.com/hashicorp/go-azure-sdk/resource-manager/cosmosdb/2025-10-15/fleets"
	"github.com/hashicorp/terraform-provider-azurerm/internal/sdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/services/cosmos/validate"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/validation"
)

//go:generate go run ../../tools/generator-tests resourceidentity -resource-name cosmosdb_fleetspace -service-package-name cosmos -properties "name,resource_group_name,fleet_name" -known-values "subscription_id:data.Subscriptions.Primary"

type CosmosDbFleetspaceResource struct{}

var (
	_ sdk.ResourceWithCustomizeDiff = CosmosDbFleetspaceResource{}
	_ sdk.ResourceWithIdentity      = CosmosDbFleetspaceResource{}
	_ sdk.ResourceWithUpdate        = CosmosDbFleetspaceResource{}
)

type CosmosDbFleetspaceModel struct {
	Name              string   `tfschema:"name"`
	ResourceGroupName string   `tfschema:"resource_group_name"`
	DataRegions       []string `tfschema:"data_regions"`
	FleetName         string   `tfschema:"fleet_name"`
	ServiceTier       string   `tfschema:"service_tier"`
	MaximumThroughput int64    `tfschema:"maximum_throughput"`
	MinimumThroughput int64    `tfschema:"minimum_throughput"`
}

func (CosmosDbFleetspaceResource) Arguments() map[string]*pluginsdk.Schema {
	return map[string]*pluginsdk.Schema{
		"name": {
			Type:         pluginsdk.TypeString,
			Required:     true,
			ForceNew:     true,
			ValidateFunc: validate.FleetName,
		},

		"resource_group_name": commonschema.ResourceGroupName(),

		"data_regions": {
			Type:     pluginsdk.TypeList,
			Required: true,
			// `ForceNew` according to https://learn.microsoft.com/en-us/azure/cosmos-db/how-to-create-fleet?pivots=azure-portal
			ForceNew: true,
			MinItems: 1,
			Elem: &pluginsdk.Schema{
				Type:             pluginsdk.TypeString,
				ValidateFunc:     location.EnhancedValidate,
				StateFunc:        location.StateFunc,
				DiffSuppressFunc: location.DiffSuppressFunc,
			},
		},

		"fleet_name": {
			Type:         pluginsdk.TypeString,
			Required:     true,
			ForceNew:     true,
			ValidateFunc: validate.FleetName,
		},

		"service_tier": {
			Type:     pluginsdk.TypeString,
			Required: true,
			// `ForceNew` according to https://learn.microsoft.com/en-us/azure/cosmos-db/how-to-create-fleet?pivots=azure-portal
			ForceNew:     true,
			ValidateFunc: validation.StringInSlice(fleets.PossibleValuesForServiceTier(), true),
		},

		"maximum_throughput": {
			Type:         pluginsdk.TypeInt,
			Optional:     true,
			ValidateFunc: validation.IntDivisibleBy(1000),
			RequiredWith: []string{
				"minimum_throughput",
			},
		},

		"minimum_throughput": {
			Type:     pluginsdk.TypeInt,
			Optional: true,
			ValidateFunc: validation.All(
				validation.IntBetween(100000, 10000000),
				validation.IntDivisibleBy(1000),
			),
			RequiredWith: []string{
				"maximum_throughput",
			},
		},
	}
}

func (CosmosDbFleetspaceResource) Attributes() map[string]*pluginsdk.Schema {
	return map[string]*pluginsdk.Schema{}
}

func (CosmosDbFleetspaceResource) ModelObject() interface{} {
	return &CosmosDbFleetspaceModel{}
}

func (CosmosDbFleetspaceResource) ResourceType() string {
	return "azurerm_cosmosdb_fleetspace"
}

func (r CosmosDbFleetspaceResource) Create() sdk.ResourceFunc {
	return sdk.ResourceFunc{
		Timeout: 30 * time.Minute,
		Func: func(ctx context.Context, metadata sdk.ResourceMetaData) error {
			client := metadata.Client.Cosmos.FleetsClient
			subscriptionId := metadata.Client.Account.SubscriptionId

			var config CosmosDbFleetspaceModel
			if err := metadata.Decode(&config); err != nil {
				return fmt.Errorf("decoding: %+v", err)
			}
			id := fleets.NewFleetspaceID(subscriptionId, config.ResourceGroupName, config.FleetName, config.Name)

			existing, err := client.FleetspaceGet(ctx, id)
			if err != nil && !response.WasNotFound(existing.HttpResponse) {
				return fmt.Errorf("checking for presence of existing %s: %+v", id, err)
			}
			if !response.WasNotFound(existing.HttpResponse) {
				return metadata.ResourceRequiresImport(r.ResourceType(), id)
			}

			param := fleets.FleetspaceResource{
				Properties: &fleets.FleetspaceProperties{
					DataRegions:                 pointer.To(config.DataRegions),
					ServiceTier:                 pointer.ToEnum[fleets.ServiceTier](config.ServiceTier),
					ThroughputPoolConfiguration: r.expandFleetspaceThroughputPoolConfiguration(config),
				},
			}

			if err := client.FleetspaceCreateThenPoll(ctx, id, param); err != nil {
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

func (r CosmosDbFleetspaceResource) Update() sdk.ResourceFunc {
	return sdk.ResourceFunc{
		Timeout: 30 * time.Minute,
		Func: func(ctx context.Context, metadata sdk.ResourceMetaData) error {
			client := metadata.Client.Cosmos.FleetsClient

			id, err := fleets.ParseFleetspaceID(metadata.ResourceData.Id())
			if err != nil {
				return err
			}

			var config CosmosDbFleetspaceModel
			if err := metadata.Decode(&config); err != nil {
				return fmt.Errorf("decoding: %+v", err)
			}

			param := fleets.FleetspaceUpdate{
				Properties: &fleets.FleetspaceProperties{},
			}

			if metadata.ResourceData.HasChange("minimum_throughput") || metadata.ResourceData.HasChange("maximum_throughput") {
				param.Properties.ThroughputPoolConfiguration = r.expandFleetspaceThroughputPoolConfiguration(config)
			}

			if err := client.FleetspaceUpdateThenPoll(ctx, *id, param); err != nil {
				return fmt.Errorf("updating %s: %+v", *id, err)
			}

			return nil
		},
	}
}

func (r CosmosDbFleetspaceResource) Read() sdk.ResourceFunc {
	return sdk.ResourceFunc{
		Timeout: 5 * time.Minute,
		Func: func(ctx context.Context, metadata sdk.ResourceMetaData) error {
			client := metadata.Client.Cosmos.FleetsClient

			id, err := fleets.ParseFleetspaceID(metadata.ResourceData.Id())
			if err != nil {
				return err
			}

			resp, err := client.FleetspaceGet(ctx, *id)
			if err != nil {
				if response.WasNotFound(resp.HttpResponse) {
					return metadata.MarkAsGone(id)
				}
				return fmt.Errorf("retrieving %s: %+v", id, err)
			}

			return r.flatten(metadata, id, resp.Model)
		},
	}
}

func (CosmosDbFleetspaceResource) Delete() sdk.ResourceFunc {
	return sdk.ResourceFunc{
		Timeout: 30 * time.Minute,
		Func: func(ctx context.Context, metadata sdk.ResourceMetaData) error {
			client := metadata.Client.Cosmos.FleetsClient

			id, err := fleets.ParseFleetspaceID(metadata.ResourceData.Id())
			if err != nil {
				return err
			}

			if err := client.FleetspaceDeleteThenPoll(ctx, *id); err != nil {
				return fmt.Errorf("deleting %s: %+v", *id, err)
			}
			return nil
		},
	}
}

func (CosmosDbFleetspaceResource) CustomizeDiff() sdk.ResourceFunc {
	return sdk.ResourceFunc{
		Timeout: 5 * time.Minute,
		Func: func(ctx context.Context, metadata sdk.ResourceMetaData) error {
			var config CosmosDbFleetspaceModel
			if err := metadata.DecodeDiff(&config); err != nil {
				return fmt.Errorf("DecodeDiff: %+v", err)
			}

			if config.MaximumThroughput > config.MinimumThroughput*10 {
				return fmt.Errorf("`maximum_throughput` should be less than or equal to 10 times of `minimum_throughput` (%d). Refer to https://learn.microsoft.com/en-us/azure/cosmos-db/fleet-pools#how-pooled-throughput-works", config.MinimumThroughput*10)
			}

			if config.MinimumThroughput > config.MaximumThroughput {
				return errors.New("`minimum_throughput` should be less than or equal to `maximum_throughput`")
			}

			return nil
		},
	}
}

func (CosmosDbFleetspaceResource) IDValidationFunc() pluginsdk.SchemaValidateFunc {
	return fleets.ValidateFleetspaceID
}

func (r CosmosDbFleetspaceResource) Identity() resourceids.ResourceId {
	return &fleets.FleetspaceId{}
}

func (CosmosDbFleetspaceResource) flatten(metadata sdk.ResourceMetaData, id *fleets.FleetspaceId, model *fleets.FleetspaceResource) error {
	state := CosmosDbFleetspaceModel{
		Name:              id.FleetspaceName,
		ResourceGroupName: id.ResourceGroupName,
		FleetName:         id.FleetName,
	}

	if model != nil {
		if props := model.Properties; props != nil {
			state.DataRegions = pointer.From(props.DataRegions)
			state.ServiceTier = string(pointer.From(props.ServiceTier))
			state.MinimumThroughput, state.MaximumThroughput = flattenFleetspaceThroughputPoolConfiguration(props.ThroughputPoolConfiguration)
		}
	}

	if err := pluginsdk.SetResourceIdentityData(metadata.ResourceData, id); err != nil {
		return err
	}

	return metadata.Encode(&state)
}

func (r CosmosDbFleetspaceResource) expandFleetspaceThroughputPoolConfiguration(config CosmosDbFleetspaceModel) *fleets.FleetspacePropertiesThroughputPoolConfiguration {
	if config.MaximumThroughput == 0 && config.MinimumThroughput == 0 {
		return nil
	}

	throughputPoolConfiguration := &fleets.FleetspacePropertiesThroughputPoolConfiguration{}
	if config.MaximumThroughput != 0 {
		throughputPoolConfiguration.MaxThroughput = pointer.To(config.MaximumThroughput)
	}
	if config.MinimumThroughput != 0 {
		throughputPoolConfiguration.MinThroughput = pointer.To(config.MinimumThroughput)
	}

	return throughputPoolConfiguration
}

func (CosmosDbFleetspaceResource) flattenFleetspaceThroughputPoolConfiguration(input *fleets.FleetspacePropertiesThroughputPoolConfiguration) (int64, int64) {
	if input == nil {
		return 0, 0
	}

	return pointer.From(input.MinThroughput), pointer.From(input.MaxThroughput)
}
