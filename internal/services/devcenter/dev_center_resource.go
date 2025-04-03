// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package devcenter

import (
	"context"
	"fmt"
	"time"

	"github.com/hashicorp/go-azure-helpers/lang/pointer"
	"github.com/hashicorp/go-azure-helpers/lang/response"
	"github.com/hashicorp/go-azure-helpers/resourcemanager/commonschema"
	"github.com/hashicorp/go-azure-helpers/resourcemanager/identity"
	"github.com/hashicorp/go-azure-helpers/resourcemanager/location"
	"github.com/hashicorp/go-azure-helpers/resourcemanager/tags"
	"github.com/hashicorp/go-azure-sdk/resource-manager/devcenter/2025-02-01/devcenters"
	"github.com/hashicorp/terraform-provider-azurerm/internal/sdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
)

var (
	_ sdk.Resource           = DevCenterResource{}
	_ sdk.ResourceWithUpdate = DevCenterResource{}
)

type DevCenterResource struct{}

func (r DevCenterResource) ModelObject() interface{} {
	return &DevCenterResourceSchema{}
}

type DevCenterResourceSchema struct {
	DevCenterUri      string                                     `tfschema:"dev_center_uri"`
	Identity          []identity.ModelSystemAssignedUserAssigned `tfschema:"identity"`
	Location          string                                     `tfschema:"location"`
	Name              string                                     `tfschema:"name"`
	ResourceGroupName string                                     `tfschema:"resource_group_name"`
	Tags              map[string]interface{}                     `tfschema:"tags"`
}

func (r DevCenterResource) IDValidationFunc() pluginsdk.SchemaValidateFunc {
	return devcenters.ValidateDevCenterID
}

func (r DevCenterResource) ResourceType() string {
	return "azurerm_dev_center"
}

func (r DevCenterResource) Arguments() map[string]*pluginsdk.Schema {
	return map[string]*pluginsdk.Schema{
		"location": commonschema.Location(),
		"name": {
			ForceNew: true,
			Required: true,
			Type:     pluginsdk.TypeString,
		},
		"resource_group_name": commonschema.ResourceGroupName(),
		"identity":            commonschema.SystemAssignedUserAssignedIdentityOptional(),
		"tags":                commonschema.Tags(),
	}
}

func (r DevCenterResource) Attributes() map[string]*pluginsdk.Schema {
	return map[string]*pluginsdk.Schema{
		"dev_center_uri": {
			Computed: true,
			Type:     pluginsdk.TypeString,
		},
	}
}

func (r DevCenterResource) Create() sdk.ResourceFunc {
	return sdk.ResourceFunc{
		Timeout: 30 * time.Minute,
		Func: func(ctx context.Context, metadata sdk.ResourceMetaData) error {
			client := metadata.Client.DevCenter.V20250201.DevCenters

			var config DevCenterResourceSchema
			if err := metadata.Decode(&config); err != nil {
				return fmt.Errorf("decoding: %+v", err)
			}

			subscriptionId := metadata.Client.Account.SubscriptionId

			id := devcenters.NewDevCenterID(subscriptionId, config.ResourceGroupName, config.Name)

			existing, err := client.Get(ctx, id)
			if err != nil {
				if !response.WasNotFound(existing.HttpResponse) {
					return fmt.Errorf("checking for the presence of an existing %s: %+v", id, err)
				}
			}
			if !response.WasNotFound(existing.HttpResponse) {
				return metadata.ResourceRequiresImport(r.ResourceType(), id)
			}

			var payload devcenters.DevCenter
			if err := r.mapDevCenterResourceSchemaToDevCenter(config, &payload); err != nil {
				return fmt.Errorf("mapping schema model to sdk model: %+v", err)
			}

			if err := client.CreateOrUpdateThenPoll(ctx, id, payload); err != nil {
				return fmt.Errorf("creating %s: %+v", id, err)
			}

			metadata.SetID(id)
			return nil
		},
	}
}

func (r DevCenterResource) Read() sdk.ResourceFunc {
	return sdk.ResourceFunc{
		Timeout: 5 * time.Minute,
		Func: func(ctx context.Context, metadata sdk.ResourceMetaData) error {
			client := metadata.Client.DevCenter.V20250201.DevCenters
			schema := DevCenterResourceSchema{}

			id, err := devcenters.ParseDevCenterID(metadata.ResourceData.Id())
			if err != nil {
				return err
			}

			resp, err := client.Get(ctx, *id)
			if err != nil {
				if response.WasNotFound(resp.HttpResponse) {
					return metadata.MarkAsGone(*id)
				}
				return fmt.Errorf("retrieving %s: %+v", *id, err)
			}

			if model := resp.Model; model != nil {
				schema.Name = id.DevCenterName
				schema.ResourceGroupName = id.ResourceGroupName
				if err := r.mapDevCenterToDevCenterResourceSchema(*model, &schema); err != nil {
					return fmt.Errorf("flattening model: %+v", err)
				}
			}

			return metadata.Encode(&schema)
		},
	}
}

func (r DevCenterResource) Delete() sdk.ResourceFunc {
	return sdk.ResourceFunc{
		Timeout: 30 * time.Minute,
		Func: func(ctx context.Context, metadata sdk.ResourceMetaData) error {
			client := metadata.Client.DevCenter.V20250201.DevCenters

			id, err := devcenters.ParseDevCenterID(metadata.ResourceData.Id())
			if err != nil {
				return err
			}

			if err := client.DeleteThenPoll(ctx, *id); err != nil {
				return fmt.Errorf("deleting %s: %+v", *id, err)
			}

			return nil
		},
	}
}

func (r DevCenterResource) Update() sdk.ResourceFunc {
	return sdk.ResourceFunc{
		Timeout: 30 * time.Minute,
		Func: func(ctx context.Context, metadata sdk.ResourceMetaData) error {
			client := metadata.Client.DevCenter.V20250201.DevCenters

			id, err := devcenters.ParseDevCenterID(metadata.ResourceData.Id())
			if err != nil {
				return err
			}

			var config DevCenterResourceSchema
			if err := metadata.Decode(&config); err != nil {
				return fmt.Errorf("decoding: %+v", err)
			}

			existing, err := client.Get(ctx, *id)
			if err != nil {
				return fmt.Errorf("retrieving existing %s: %+v", *id, err)
			}
			if existing.Model == nil {
				return fmt.Errorf("retrieving existing %s: properties was nil", *id)
			}
			payload := *existing.Model

			if err := r.mapDevCenterResourceSchemaToDevCenter(config, &payload); err != nil {
				return fmt.Errorf("mapping schema model to sdk model: %+v", err)
			}

			if err := client.CreateOrUpdateThenPoll(ctx, *id, payload); err != nil {
				return fmt.Errorf("updating %s: %+v", *id, err)
			}

			return nil
		},
	}
}

func (r DevCenterResource) mapDevCenterResourceSchemaToDevCenter(input DevCenterResourceSchema, output *devcenters.DevCenter) error {
	identity, err := identity.ExpandSystemAndUserAssignedMapFromModel(input.Identity)
	if err != nil {
		return fmt.Errorf("expanding SystemAndUserAssigned Identity: %+v", err)
	}
	output.Identity = identity

	output.Location = location.Normalize(input.Location)
	output.Tags = tags.Expand(input.Tags)

	if output.Properties == nil {
		output.Properties = &devcenters.DevCenterProperties{}
	}

	return nil
}

func (r DevCenterResource) mapDevCenterToDevCenterResourceSchema(input devcenters.DevCenter, output *DevCenterResourceSchema) error {
	identity, err := identity.FlattenSystemAndUserAssignedMapToModel(input.Identity)
	if err != nil {
		return fmt.Errorf("flattening SystemAndUserAssigned Identity: %+v", err)
	}
	output.Identity = *identity

	output.Location = location.Normalize(input.Location)
	output.Tags = tags.Flatten(input.Tags)

	if input.Properties == nil {
		input.Properties = &devcenters.DevCenterProperties{}
	}

	output.DevCenterUri = pointer.From(input.Properties.DevCenterUri)

	return nil
}
